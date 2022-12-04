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

	"github.com/mitchellh/mapstructure"
)

type FindOptions struct {
	Path    string
	Pattern string
	Type    string
}

// Find, meant to be used as a goroutine //////////////////////////////////////
func (m *MiriaClient) Find(opt FindOptions, cout chan []SearchResult, cerr chan error) {
	// form search request from pattern
	var req FindInstanceRequest

	req.RootObjectPath = opt.Path
	req.ResultType = "INST"
	req.PageSize = 3000
	req.Criteria.Condition = "AND"
	split := strings.Split(opt.Pattern, "*")

	addRule := func(t string, val any, op string) {
		req.Criteria.Rules = append(req.Criteria.Rules, FindRule{
			Type:     t,
			Value:    val,
			Value2:   nil,
			Operator: op,
		})
	}

	// parsing serach pattern
	if opt.Pattern == "*" || opt.Pattern == "" {
		// match everything
		addRule("FILE_NAME", "", "contains")
	} else if len(split) == 1 {
		// no wildcard
		addRule("FILE_NAME", split[0], "equal")
	} else {
		// first character is not a wildcard
		if split[0] != "" {
			addRule("FILE_NAME", split[0], "starts with")
		}
		// strings between wildcards
		for _, s := range split[1 : len(split)-1] {
			addRule("FILE_NAME", s, "contains")
		}
		// last character is not a wildcard
		if split[len(split)-1] != "" {
			addRule("FILE_NAME", split[len(split)-1], "ends with")
		}
	}

	// adding rule for type filtering
	if opt.Type != "" {
		if opt.Type == "f" {
			addRule("FILE_TYPE", 1, "equals to")
		} else if opt.Type == "d" {
			addRule("FILE_TYPE", 2, "equals to")
		} else {
			cerr <- fmt.Errorf("unknown file type '%s'", opt.Type)
			return
		}
	}

	// execute request
	var searchResp SearchResponse
	resp, err := m.Post("/files/advanced-search/", req, true)
	if err != nil {
		cerr <- err
		return
	}
	mapstructure.Decode(resp, &searchResp)
	cout <- searchResp.Results
	nextPage := resp["nextPage"]
	for nextPage != nil {
		nextPageEnc := url.QueryEscape(nextPage.(string))
		resp, err := m.Post("/files/advanced-search/?page="+nextPageEnc, req, true)
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
