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

	"github.com/aportelli/miria-cli/client"
	"github.com/aportelli/miria-cli/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var miria *client.Miria

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "miria",
	Short: "CLI for Atempo Miria",
	Long:  `Command-line interface tool to interact with an Atempo Miria tape storage solution.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&log.Level, "verbosity", "v", 0,
		"verbosity level (0: default, 1: info, 2: debug)")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	userConfigDir, err := os.UserConfigDir()
	log.ErrorCheck(err, "cannot find user config directory")
	viper.AddConfigPath(userConfigDir + "/miria")
	viper.SafeWriteConfig()
	err = viper.ReadInConfig()
	log.ErrorCheck(err, "cannot read user config")
	host := viper.Get("host")
	if host != nil {
		miria = client.NewMiria(host.(string))
	} else {
		miria = client.NewMiria("")
	}
}
