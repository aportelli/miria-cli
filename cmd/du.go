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

	"github.com/aportelli/miria-cli/client"
	"github.com/aportelli/miria-cli/log"
	"github.com/spf13/cobra"
)

// duCmd represents the du command
var duCmd = &cobra.Command{
	Use:   "du <path>",
	Short: "Show total size contained in a given path",
	Long: `Show total size contained in a given path, mimicking the ` + "`du -s`" + ` Unix command.
Miria does not have a direct interface to query directory sizes, this command 
will perform a full scan similar to the ` + "`find`" + ` command, and might take time 
for large directories.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		AuthenticateIfNecessary()
		var total uint64 = 0
		findOpt.Opt.Path = args[0]
		cout := make(chan []client.SearchResult)
		cerr := make(chan error)
		go miria.Find(findOpt.Opt, cout, cerr)
	out:
		for {
			select {
			case err := <-cerr:
				log.ErrorCheck(err, "")
				return
			case buf := <-cout:
				if buf == nil {
					break out
				} else {
					for _, r := range buf {
						total += r.ObjectSize
					}
				}
			}
		}
		var size string
		if findOpt.Humanize {
			size = SizeString(ByteSize(total))
		} else {
			size = fmt.Sprint(total)
		}
		log.Msg.Printf("%s %s", size, findOpt.Opt.Path)
	},
}

var duOpt = struct {
	Opt      client.FindOptions
	Humanize bool
}{client.FindOptions{Path: "", Type: "f", Pattern: "*"}, false}

func init() {
	rootCmd.AddCommand(duCmd)
	duCmd.Flags().BoolVarP(&findOpt.Humanize, "human-readable", "H", false,
		"human-readable sizes")
	duCmd.Flags().Lookup("human-readable").NoOptDefVal = "true"
}
