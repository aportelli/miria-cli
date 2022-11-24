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
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure miria-cli",
	Long:  `Modify the options for miria-cli. The list of all options can be ob`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List possible options",
	Run: func(cmd *cobra.Command, args []string) {
		for _, opt := range options {
			fmt.Println(opt)
		}
	},
}

var configFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Get path of current config file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.ConfigFileUsed())
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set option",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		for _, opt := range options {
			if opt == args[0] {
				viper.Set(opt, args[1])
				err := viper.WriteConfig()
				if err != nil {
					log.Fatalf("error: %s\nerror: cannot write config file", err.Error())
				}
				return
			}
		}
		log.Fatal("error: option '" + args[0] + "' does not exist, use `miria config list` to see all possible options")
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get option",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("error: %s\nerror: cannot read config file", err.Error())
		}
		for _, opt := range options {
			if opt == args[0] {
				val := viper.Get(opt)
				fmt.Println(val)
				return
			}
		}
		log.Fatal("error: option '" + args[0] + "' does not exist, use `miria config list` to see all possible options")
	},
}

var options = []string{"host"}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configFileCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
}
