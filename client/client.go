/*
 * client.go, part of miria-cli (https://github.com/aportelli/miria-cli)
 * Copyright (C) 2022 Antonin Portelli
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the GNU General Public License as published by the Free Software
 * Foundation, either version 3 of the License, or (at your option) any later
 * version.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of  MERCHANTABILITY or FITNESS FOR
 * A PARTICULAR PURPOSE. See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	glob "github.com/aportelli/miria-cli/global"
	req "github.com/aportelli/miria-cli/requests"
	resp "github.com/aportelli/miria-cli/responses"
)

type client struct {
	apiUrl string
	auth   resp.AuthToken
}

/**************************************************************************************************
 *                                     Client creation                                            *
 **************************************************************************************************/
func NewClient(host string) *client {
	c := new(client)
	c.apiUrl = "http://" + host + "/restapi"

	return c
}

/**************************************************************************************************
 *                                      HTTP requests                                             *
 **************************************************************************************************/
// Private method to handle abstract requests //////////////////////////////////////////////////////
func (c *client) executeRequest(request *http.Request, authenticate bool) (map[string]any, error) {
	// complete request
	request.Header.Set("Content-Type", "application/json")
	if authenticate {
		err := c.CheckAuthentication()
		if err != nil {
			return nil, err
		}
		request.Header.Add("Authorization", "Bearer "+c.auth.Access)
	}

	// execute request
	var raw map[string]any

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	jbuf, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jbuf, &raw)
	if err != nil {
		return nil, err
	}
	if _, ok := raw["error"]; ok {
		message, _ := json.MarshalIndent(raw, "", "  ")
		return nil, fmt.Errorf("error: Miria server returned an error\n%s", string(message))
	}
	return raw, nil
}

// POST ////////////////////////////////////////////////////////////////////////////////////////////
func (c *client) Post(path string, body any, authenticate bool) (map[string]any, error) {
	jbuf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", c.apiUrl+path, bytes.NewBuffer(jbuf))
	if err != nil {
		return nil, err
	}

	return c.executeRequest(request, authenticate)
}

// GET /////////////////////////////////////////////////////////////////////////////////////////////
func (c *client) Get(path string, authenticate bool) (map[string]any, error) {
	request, err := http.NewRequest("GET", c.apiUrl+path, nil)
	if err != nil {
		return nil, err
	}

	return c.executeRequest(request, authenticate)
}

/**************************************************************************************************
 *                                      Authentication                                            *
 **************************************************************************************************/
// Private method to cache token ///////////////////////////////////////////////////////////////////
func (c *client) cacheAuthentication() error {
	fileContent, err := json.MarshalIndent(c.auth, "", "  ")
	if err != nil {
		return err
	}
	authPath, err := glob.AppCacheDir()
	if err != nil {
		return err
	}
	err = os.MkdirAll(authPath, 0700)
	if err != nil {
		return err
	}
	authPath += "/auth.json"
	err = os.WriteFile(authPath, fileContent, 0600)
	if err != nil {
		return err
	}
	return nil
}

// Obtain token from username/password /////////////////////////////////////////////////////////////
func (c *client) Authenticate(username string, password string) error {
	var request req.AuthRequest

	request.Db = "ADA"
	request.SuperUser = false
	request.Name = username
	request.Password = password
	auth, err := c.Post("/auth/token/", request, false)
	if err != nil {
		return err
	}
	err = glob.JsonMapToStruct(&c.auth, auth)
	if err != nil {
		return err
	}
	err = c.cacheAuthentication()
	if err != nil {
		return err
	}

	return nil
}

// Check if token exists and is still valid ////////////////////////////////////////////////////////
func (c *client) CheckAuthentication() error {
	// check for cache auth.json file
	authPath, err := glob.AppCacheDir()
	if err != nil {
		return err
	}
	authPath += "/auth.json"
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
			return fmt.Errorf("error: could not refresh token")
		}
		err = glob.JsonMapToStruct(&c.auth, response)
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
