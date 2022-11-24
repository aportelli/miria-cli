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

	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find <path>",
	Short: "Find files in archive",
	Long: `Find files in the tape archive.
Example:
  miria find archive@project:/dir -name '*.txt'`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("find called", findPath)
	},
}

var findPath string

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.Flags().StringVarP(&findPath, "name", "n", "", "search pattern")
	findCmd.MarkFlagRequired("name")
}
