/*
 * miria-cli.go, part of miria-cli (https://github.com/aportelli/miria-cli)
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

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/aportelli/miria-cli/client"
	glob "github.com/aportelli/miria-cli/global"
	"golang.org/x/term"
)

func main() {
	if len(os.Args[1:]) != 1 {
		fmt.Fprintln(os.Stderr, "usage:", path.Base(os.Args[0]), "<Miria host>")
		os.Exit(1)
	}
	host := os.Args[1]

	c := client.NewClient(host)

	var username string
	err := c.CheckAuthentication()
	if err != nil {
		fmt.Print("Enter username: ")
		fmt.Scanln(&username)
		fmt.Print("Enter password: ")
		bytepw, _ := term.ReadPassword(int(syscall.Stdin))
		fmt.Println("")
		err = c.Authenticate(username, string(bytepw))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	fmt.Println("Authenticated!")
	response, err := c.Get("/datamanagement/repositories/", true)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	glob.PrettyPrintResponse(response)

	search := []byte(`{
			"rootObjectPath": "archive@dp207",
	  	"resultType": "INST",
			"pageSize": 1000,
			"criteria": {
				"condition": "AND",
				"rules": [
					{
						"type": "FILE_NAME",
						"value": "",
						"value2": null,
						"operator": "contains"
					}
				]
			}
		}
	`)
	var searchMap map[string]any
	json.Unmarshal(search, &searchMap)
	response, err = c.Post("/files/advanced-search/", searchMap, true)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	glob.PrettyPrintResponse(response)
}
