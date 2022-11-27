/*
Copyright © 2022 Antonin Portelli <antonin.portelli@me.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/aportelli/miria-cli/log"
	"github.com/mitchellh/mapstructure"
)

type restClient struct {
	apiUrl string
	host   string
	auth   AuthToken
}

/******************************************************************************
 *                              Client creation                               *
 ******************************************************************************/
func NewRestClient(host string) *restClient {
	c := new(restClient)
	c.host = host
	c.apiUrl = "http://" + c.host + "/restapi"

	return c
}

/******************************************************************************
 *                              HTTP requests                                 *
 ******************************************************************************/
// Private method to handle abstract requests /////////////////////////////////
func (c *restClient) executeRequest(request *http.Request, authenticate bool) (map[string]any, error) {
	// error if host is empty
	if c.host == "" {
		return nil, fmt.Errorf("host empty, please configure a host with `miria config set host <host>`")
	}

	// complete request
	request.Header.Set("Content-Type", "application/json")
	if authenticate {
		levelCopy := log.AtMostLevel(0)
		err := c.CheckAuthentication()
		log.Level = levelCopy
		if err != nil {
			return nil, err
		}
		request.Header.Add("Authorization", "Bearer "+c.auth.Access)
	}
	if log.Level >= 2 {
		log.Dbg.Println("* Request headers")
		for k, v := range request.Header {
			buf := k + ": "
			if k == "Authorization" {
				buf += "Bearer <redacted>"
			} else {
				for i, val := range v {
					buf += val
					if i < len(v)-1 {
						buf += ", "
					}
				}
			}
			log.Dbg.Println(buf)
		}
	}

	// execute request
	var raw map[string]any

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&raw)
	if err != nil {
		return nil, err
	}
	if _, ok := raw["error"]; ok {
		message, _ := json.MarshalIndent(raw, "", "  ")
		return nil, fmt.Errorf("Miria server returned an error\n%s", string(message))
	}
	return raw, nil
}

// POST ///////////////////////////////////////////////////////////////////////
func (c *restClient) Post(path string, body any, authenticate bool) (map[string]any, error) {
	jbuf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", c.apiUrl+path, bytes.NewBuffer(jbuf))
	if err != nil {
		return nil, err
	}
	log.Inf.Printf("POST %s", c.apiUrl+path)
	if log.Level >= 2 {
		jbufPretty, _ := json.MarshalIndent(body, "", "  ")
		log.Dbg.Printf("* Request data\n%s", string(jbufPretty))
	}

	return c.executeRequest(request, authenticate)
}

// GET ////////////////////////////////////////////////////////////////////////
func (c *restClient) Get(path string, authenticate bool) (map[string]any, error) {
	request, err := http.NewRequest("GET", c.apiUrl+path, nil)
	if err != nil {
		return nil, err
	}
	log.Inf.Printf(" GET %s", c.apiUrl+path)

	return c.executeRequest(request, authenticate)
}

/******************************************************************************
 *                              Authentication                                *
 ******************************************************************************/
// Private method to cache token //////////////////////////////////////////////
func (c *restClient) AuthenticationCache() (string, error) {
	authPath, err := AppCacheDir()
	if err != nil {
		return "", err
	}
	authPath += "/auth.json"
	return authPath, nil
}

func (c *restClient) cacheAuthentication() error {
	fileContent, err := json.MarshalIndent(c.auth, "", "  ")
	if err != nil {
		return err
	}
	authPath, err := c.AuthenticationCache()
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Dir(authPath), 0700)
	if err != nil {
		return err
	}
	err = os.WriteFile(authPath, fileContent, 0600)
	if err != nil {
		return err
	}
	return nil
}

// Obtain token from username/password ////////////////////////////////////////
func (c *restClient) Authenticate(username string, password string) error {
	var request AuthRequest

	request.Db = "ADA"
	request.SuperUser = false
	request.Name = username
	request.Password = password
	log.Dbg.Println("warning: debug output deactivated during authentication")
	levelCopy := log.AtMostLevel(1)
	auth, err := c.Post("/auth/token/", request, false)
	log.Level = levelCopy
	if err != nil {
		return err
	}
	err = mapstructure.Decode(auth, &c.auth)
	if err != nil {
		return err
	}
	err = c.cacheAuthentication()
	if err != nil {
		return err
	}

	return nil
}

// Check if token exists and is still valid ///////////////////////////////////
func (c *restClient) CheckAuthentication() error {
	// check for cache auth.json file
	authPath, err := c.AuthenticationCache()
	if err != nil {
		return err
	}
	_, err = os.Open(authPath)
	if err == nil {
		authj, err := os.ReadFile(authPath)
		if err != nil {
			return err
		}
		json.Unmarshal(authj, &c.auth)
	} else {
		return err
	}

	// check for token validity and refresh if invalid
	body := map[string]string{"token": c.auth.Access}
	response, err := c.Post("/auth/token/verify/", body, false)
	if err != nil {
		return err
	}
	if val, ok := response["code"]; (val == "token_not_valid") && ok {
		body = map[string]string{"refresh": c.auth.Refresh}
		response, err = c.Post("/auth/token/refresh/", body, false)
		if err != nil {
			return err
		}
		if _, ok := response["code"]; ok {
			return fmt.Errorf("could not refresh token")
		}
		err = mapstructure.Decode(response, &c.auth)
		if err != nil {
			return err
		}
		err = c.cacheAuthentication()
		if err != nil {
			return err
		}
	}
	return nil
}
