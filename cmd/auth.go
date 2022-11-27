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
	"os"

	"github.com/aportelli/miria-cli/log"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset authentication",
	Long:  `Reset authentication token, you will be asked your user name and password.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		authPath, err := miria.Client.AuthenticationCache()
		log.ErrorCheck(err, "")
		err = os.RemoveAll(authPath)
		log.ErrorCheck(err, "")
		err = miria.AuthenticateInteractive()
		log.ErrorCheck(err, "")
		log.Msg.Println("Authentication token successfully reset")
	},
}

var authCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check authentication",
	Long: `Check authentication token and try to refresh it if necessary, 
exit with status 1 in case of failure`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := miria.Client.CheckAuthentication()
		log.ErrorCheck(err, "authentication check failed")
		log.Msg.Println("Authentication token valid")
	},
}

var authFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Get path of authentication cache",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := miria.Client.AuthenticationCache()
		log.ErrorCheck(err, "cannot get authentication cache path")
		log.Msg.Println(path)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authCheckCmd)
	authCmd.AddCommand(authFileCmd)
	authCmd.AddCommand(authResetCmd)
}
