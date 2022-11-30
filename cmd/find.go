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
	"github.com/aportelli/miria-cli/client"
	"github.com/aportelli/miria-cli/log"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find <path>",
	Short: "Find files in archive, ",
	Long: `Find files in the tape archive, mimicking the ` + "`find`" + ` Unix command.

Example:
  miria find archive@project:/dir --name '*.txt'`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		AuthenticateIfNecessary()
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
					if findOpt.List {
						if findOpt.Humanize {
							for _, r := range buf {
								log.Msg.Printf("%6s %6s %s %s", r.ObjectType,
									SizeString((ByteSize)(r.ObjectSize)), r.InstanceBackupDate, r.ObjectPath)
							}
						} else {
							for _, r := range buf {
								log.Msg.Printf("%6s %12d %s %s", r.ObjectType, r.ObjectSize, r.InstanceBackupDate,
									r.ObjectPath)
							}
						}

					} else {
						for _, r := range buf {
							log.Msg.Println(r.ObjectPath)
						}
					}
				}
			}
		}
	},
}

var findOpt = struct {
	Opt      client.FindOptions
	List     bool
	Humanize bool
}{client.FindOptions{Path: "", Type: "", Pattern: ""}, false, false}

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.Flags().StringVarP(&findOpt.Opt.Pattern, "name", "n", "", "search pattern")
	findCmd.Flags().StringVarP(&findOpt.Opt.Type, "type", "t", "",
		"filter file type (d or f)")
	findCmd.Flags().BoolVarP(&findOpt.Humanize, "human-readable", "H", false,
		"human-readable sizes")
	findCmd.Flags().Lookup("human-readable").NoOptDefVal = "true"
	findCmd.Flags().BoolVarP(&findOpt.List, "list", "l", false,
		"columns with file type and size")
	findCmd.Flags().Lookup("list").NoOptDefVal = "true"
}
