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
package log

import (
	"log"

	"github.com/fatih/color"
)

type Logger struct {
	Level     int
	Color     *color.Color
	StdLogger *log.Logger
}

var Level int = 0

func AtMostLevel(maxLevel int) int {
	logCopy := Level
	if logCopy >= maxLevel {
		Level = maxLevel
	}
	return logCopy
}

func (l Logger) printFn(fn func()) {
	if Level >= l.Level {
		l.Color.Set()
		defer color.Unset()
		fn()
	}
}

func (l Logger) Println(a ...any) {
	l.printFn(func() { l.StdLogger.Println(a...) })
}

func (l Logger) Fatalln(a ...any) {
	l.printFn(func() { l.StdLogger.Fatalln(a...) })
}

func (l Logger) Printf(format string, a ...any) {
	l.printFn(func() { l.StdLogger.Printf(format, a...) })
}

func (l Logger) Fatalf(format string, a ...any) {
	l.printFn(func() { l.StdLogger.Fatalf(format, a...) })
}
