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
	"log"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		authPath, err := miria.Client.AuthenticationCache()
		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
		err = os.RemoveAll(authPath)
		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
		err = miria.AuthenticateInteractive()
		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
		log.Println("Authentication token successfully reset")
	},
}

var authCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check authentication",
	Long: `Check authentication token and try to refresh it if necessary, 
exit with status 1 in case of failure`,
	Run: func(cmd *cobra.Command, args []string) {
		err := miria.Client.CheckAuthentication()
		if err != nil {
			log.Fatalf("error: %s\nerror: authentication check failed", err.Error())
		}
		log.Println("Authentication token valid")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authResetCmd)
	authCmd.AddCommand(authCheckCmd)
}
