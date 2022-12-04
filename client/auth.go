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
	"encoding/json"
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/aportelli/miria-cli/log"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/term"
)

// Private method to cache token //////////////////////////////////////////////
func (m *MiriaClient) AuthenticationCache() (string, error) {
	authPath, err := AppCacheDir()
	if err != nil {
		return "", err
	}
	authPath += "/auth.json"
	return authPath, nil
}

func (m *MiriaClient) cacheAuthentication() error {
	fileContent, err := json.MarshalIndent(m.auth, "", "  ")
	if err != nil {
		return err
	}
	authPath, err := m.AuthenticationCache()
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
func (m *MiriaClient) Authenticate(username string, password string) error {
	var request AuthRequest

	request.Db = "ADA"
	request.SuperUser = false
	request.Name = username
	request.Password = password
	log.Dbg.Println("warning: debug output deactivated during authentication")
	levelCopy := log.AtMostLevel(1)
	auth, err := m.Post("/auth/token/", request, false)
	log.Level = levelCopy
	if err != nil {
		return err
	}
	err = mapstructure.Decode(auth, &m.auth)
	if err != nil {
		return err
	}
	err = m.cacheAuthentication()
	if err != nil {
		return err
	}

	return nil
}

// Check if token exists and is still valid ///////////////////////////////////
func (m *MiriaClient) CheckAuthentication() error {
	// check for cache auth.json file
	authPath, err := m.AuthenticationCache()
	if err != nil {
		return err
	}
	_, err = os.Open(authPath)
	if err == nil {
		authj, err := os.ReadFile(authPath)
		if err != nil {
			return err
		}
		json.Unmarshal(authj, &m.auth)
	} else {
		return err
	}

	// check for token validity and refresh if invalid
	body := map[string]string{"token": m.auth.Access}
	_, err = m.Post("/auth/token/verify/", body, false)
	if err != nil {
		body = map[string]string{"refresh": m.auth.Refresh}
		response, err := m.Post("/auth/token/refresh/", body, false)
		if err != nil {
			return err
		}
		err = mapstructure.Decode(response, &m.auth)
		if err != nil {
			return err
		}
		err = m.cacheAuthentication()
		if err != nil {
			return err
		}
	}
	return nil
}

// Interactive authentication /////////////////////////////////////////////////
func (m *MiriaClient) AuthenticateInteractive(force bool) error {
	err := m.CheckAuthentication()
	if force || err != nil {
		var username string

		fmt.Print("Enter username: ")
		fmt.Scanln(&username)
		fmt.Print("Enter password: ")
		bytepw, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println("")
		err = m.Authenticate(username, string(bytepw))
		if err != nil {
			return err
		}
	}
	return nil
}
