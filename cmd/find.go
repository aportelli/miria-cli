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
	"strings"

	log "github.com/aportelli/golog"
	"github.com/aportelli/miria-cli/client"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find <path>",
	Short: "Find files in archive",
	Long: `Find files in the tape archive, mimicking the ` + "`find`" + ` Unix command.

Example:
  miria find archive@project:/dir --name '*.txt'`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		AuthenticateIfNecessary()
		findOpt.Opt.Path = args[0]
		cout := make(chan []client.SearchResult)
		cerr := make(chan error)
		rootDepth := depth(findOpt.Opt.Path)
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
						if findOpt.MaxDepth < 0 || depth(r.ObjectPath)-rootDepth <= findOpt.MaxDepth {
							if findOpt.List {
								if findOpt.Humanize {
									log.Msg.Printf("%6s %6s %20s %s", r.ObjectType,
										log.SizeString((log.ByteSize)(r.ObjectSize)), r.InstanceBackupDate, r.ObjectPath)
								} else {
									log.Msg.Printf("%6s %12d %20s %s", r.ObjectType, r.ObjectSize,
										r.InstanceBackupDate, r.ObjectPath)
								}
							} else {
								log.Msg.Println(r.ObjectPath)
							}
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
	MaxDepth int
}{client.FindOptions{Path: "", Type: "", Pattern: ""}, false, false, -1}

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
	findCmd.Flags().IntVarP(&findOpt.MaxDepth, "max-depth", "d", -1, "maximum depth (-1 is unlimited)")
}

func depth(path string) int {
	var depth int = 0
	colSplit := strings.Split(path, ":")
	slashSplit := strings.Split(colSplit[1], "/")
	for _, name := range slashSplit {
		if name != "" {
			depth++
		}
	}
	if depth == 0 {
		return depth
	} else {
		return depth - 1
	}
}
