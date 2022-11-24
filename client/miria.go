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
	"syscall"

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

func (m *Miria) AuthenticateInteractive() (err error) {
	err = m.Client.CheckAuthentication()
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

func (m *Miria) Find(pattern string) {

}
