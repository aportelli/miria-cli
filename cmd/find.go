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
		pattern := args[0]
		miria.Find(pattern, findOpt.Path)
	},
}

var findOpt = struct{ Path string }{""}

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.Flags().StringVarP(&findOpt.Path, "name", "n", "", "search pattern")
}
