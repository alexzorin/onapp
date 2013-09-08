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

func Infof(format string, args ...interface{}) {
	println(format, "", args)
}

func Infoln(args ...interface{}) {
	println("", "", args)
}

func InfoToggle(on bool) {
}

func Errorln(args ...interface{}) {
	println("%s\n", ERROR_COLOR, args)
}

func Errorf(fmt string, args ...interface{}) {
	println(fmt, ERROR_COLOR, args)
}

func Warnln(args ...interface{}) {
	println("%s\n", WARN_COLOR, args)
}

func Warnf(fmt string, args ...interface{}) {
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
