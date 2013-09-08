// +build !windows

package cmd

import (
	"fmt"
)

const (
	ERROR_COLOR = "31"
	WARN_COLOR  = "33"
	ESC_START   = "\x1b["
	ESC_STOP    = "m"
)

func infof(format string, args ...interface{}) {
	println(format, "", args)
}

func infoln(args ...interface{}) {
	println("", "", args)
}

func infoToggle(on bool) {
}

func errorln(args ...interface{}) {
	println("%s\n", ERROR_COLOR, args)
}

func errorf(fmt string, args ...interface{}) {
	println(fmt, ERROR_COLOR, args)
}

func warnln(args ...interface{}) {
	println("%s\n", WARN_COLOR, args)
}

func warnf(fmt string, args ...interface{}) {
	println(fmt, WARN_COLOR, args)
}

func println(format string, color string, args interface{}) {
	fmt.Printf("%s%s", ESC_START, color)
	if format == "" {
		fmt.Println((args.([]interface{}))...)
	} else {
		fmt.Printf(format, (args.([]interface{}))...)
	}
	fmt.Printf(ESC_STOP)
}
