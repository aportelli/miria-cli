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

	log "github.com/aportelli/golog"
)

// Private method to handle abstract requests /////////////////////////////////
func (m *MiriaClient) executeRequest(request *http.Request, authenticate bool) (map[string]any, error) {
	// error if host is empty
	if m.host == "" {
		return nil, fmt.Errorf("host empty, please configure a host with `miria config set host <host>`")
	}

	// complete request
	request.Header.Set("Content-Type", "application/json")
	if authenticate {
		levelCopy := log.AtMostLevel(0)
		err := m.CheckAuthentication()
		log.Level = levelCopy
		if err != nil {
			return nil, err
		}
		request.Header.Add("Authorization", "Bearer "+m.auth.Access)
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
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("the Miria server returned HTTP response %d", response.StatusCode)
	}
	dec := json.NewDecoder(response.Body)
	dec.Decode(&raw)
	if err != nil {
		return nil, err
	}
	if _, ok := raw["error"]; ok {
		message, _ := json.MarshalIndent(raw, "", "  ")
		return nil, fmt.Errorf("the Miria server returned an error\n%s", string(message))
	}
	return raw, nil
}

// POST ///////////////////////////////////////////////////////////////////////
func (m *MiriaClient) Post(path string, body any, authenticate bool) (map[string]any, error) {
	jbuf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", m.apiUrl+path, bytes.NewBuffer(jbuf))
	if err != nil {
		return nil, err
	}
	log.Inf.Printf("POST %s", m.apiUrl+path)
	if log.Level >= 2 {
		jbufPretty, _ := json.MarshalIndent(body, "", "  ")
		log.Dbg.Printf("* Request data\n%s", string(jbufPretty))
	}

	return m.executeRequest(request, authenticate)
}

// GET ////////////////////////////////////////////////////////////////////////
func (m *MiriaClient) Get(path string, authenticate bool) (map[string]any, error) {
	request, err := http.NewRequest("GET", m.apiUrl+path, nil)
	if err != nil {
		return nil, err
	}
	log.Inf.Printf(" GET %s", m.apiUrl+path)

	return m.executeRequest(request, authenticate)
}
