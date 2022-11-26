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
	"fmt"
	"net/url"
	"strings"
	"syscall"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/term"
)

type Miria struct {
	Client *restClient
}

/******************************************************************************
 *                      High-level client creation                            *
 ******************************************************************************/
func NewMiria(host string) *Miria {
	m := new(Miria)
	m.Client = NewRestClient(host)
	return m
}

func (m *Miria) AuthenticateInteractive() error {
	err := m.Client.CheckAuthentication()
	if err != nil {
		var username string

		fmt.Print("Enter username: ")
		fmt.Scanln(&username)
		fmt.Print("Enter password: ")
		bytepw, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println("")
		err = m.Client.Authenticate(username, string(bytepw))
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Miria) Find(path string, pattern string, cout chan []SearchResult, cerr chan error) {
	// form search request from pattern
	var req FindInstanceRequest

	req.RootObjectPath = path
	req.ResultType = "INST"
	req.PageSize = 3000
	req.Criteria.Condition = "AND"
	split := strings.Split(pattern, "*")
	if pattern == "*" || pattern == "" {
		// match everything
		rule := FindRule{
			Type:     "FILE_NAME",
			Value:    "",
			Value2:   nil,
			Operator: "contains",
		}
		req.Criteria.Rules = append(req.Criteria.Rules, rule)
	} else if len(split) == 1 {
		// no wildcard
		rule := FindRule{
			Type:     "FILE_NAME",
			Value:    split[0],
			Value2:   nil,
			Operator: "equal",
		}
		req.Criteria.Rules = append(req.Criteria.Rules, rule)
	} else {
		// first character is not a wildcard
		if split[0] != "" {
			rule := FindRule{
				Type:     "FILE_NAME",
				Value:    split[0],
				Value2:   nil,
				Operator: "starts with",
			}
			req.Criteria.Rules = append(req.Criteria.Rules, rule)
		}
		// strings between wildcards
		for _, s := range split[1 : len(split)-1] {
			rule := FindRule{
				Type:     "FILE_NAME",
				Value:    s,
				Value2:   nil,
				Operator: "contains",
			}
			req.Criteria.Rules = append(req.Criteria.Rules, rule)
		}
		// last character is not a wildcard
		if split[len(split)-1] != "" {
			rule := FindRule{
				Type:     "FILE_NAME",
				Value:    split[len(split)-1],
				Value2:   nil,
				Operator: "ends with",
			}
			req.Criteria.Rules = append(req.Criteria.Rules, rule)
		}
	}

	// execute request
	var searchResp SearchResponse
	resp, err := m.Client.Post("/files/advanced-search/", req, true)
	if err != nil {
		cerr <- err
		return
	}
	mapstructure.Decode(resp, &searchResp)
	cout <- searchResp.Results
	nextPage := resp["nextPage"]
	for nextPage != nil {
		nextPageEnc := url.QueryEscape(nextPage.(string))
		resp, err := m.Client.Post("/files/advanced-search/?page="+nextPageEnc, req, true)
		if err != nil {
			cerr <- err
			return
		}
		mapstructure.Decode(resp, &searchResp)
		cout <- searchResp.Results
		nextPage = resp["nextPage"]
	}
	cout <- nil
}
