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
package cmd

import (
	"encoding/json"

	"github.com/aportelli/miria-cli/log"
	"github.com/spf13/cobra"
)

var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Manual REST API request",
	Long:  `Manual REST API request to Miria server, requests are authenticated by default.`,
}

var restGetCmd = &cobra.Command{
	Use:   "GET <path>",
	Short: "GET request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		if !restOpt.NoAuth {
			AuthenticateIfNecessary()
		}
		response, err := miria.Get(path, !restOpt.NoAuth)
		log.ErrorCheck(err, "")
		jbuf, err := json.MarshalIndent(response, "", "  ")
		log.ErrorCheck(err, "")
		log.Msg.Println(string(jbuf))
	},
}

var restPostCmd = &cobra.Command{
	Use:   "POST <path> <body (JSON)>",
	Short: "POST request",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var body map[string]any
		path := args[0]
		bodyJson := args[1]
		if !restOpt.NoAuth {
			AuthenticateIfNecessary()
		}
		err := json.Unmarshal([](byte)(bodyJson), &body)
		log.ErrorCheck(err, "")
		response, err := miria.Post(path, body, !restOpt.NoAuth)
		log.ErrorCheck(err, "")
		jbuf, err := json.MarshalIndent(response, "", "  ")
		log.ErrorCheck(err, "")
		log.Msg.Println(string(jbuf))
	},
}

var restOpt = struct{ NoAuth bool }{false}

func init() {
	rootCmd.AddCommand(restCmd)
	restCmd.AddCommand(restGetCmd)
	restCmd.AddCommand(restPostCmd)
	restCmd.PersistentFlags().BoolVarP(&restOpt.NoAuth, "noauth", "", false, "do not authenticate")
}
