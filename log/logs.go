package log

import (
	"log"
	"os"
	"runtime/debug"

	"github.com/fatih/color"
)

var Dbg = Logger{
	Level:     2,
	Color:     color.New(color.FgYellow),
	StdLogger: log.New(os.Stdout, "DEBUG: ", 0),
}

var Err = Logger{
	Level:     0,
	Color:     color.New(color.FgRed, color.Bold),
	StdLogger: log.New(os.Stderr, "ERROR: ", 0),
}

var Inf = Logger{
	Level:     1,
	Color:     color.New(color.FgHiBlack),
	StdLogger: log.New(os.Stdout, " INFO: ", 0),
}

var Msg = Logger{
	Level:     0,
	Color:     color.New(),
	StdLogger: log.New(os.Stdout, "", 0),
}

func ErrorCheck(err error, message string) {
	if err != nil {
		Err.Println(err.Error())
		if message != "" {
			Err.Println(message)
		}
		Dbg.Println(string(debug.Stack()))
		os.Exit(1)
	}
}
