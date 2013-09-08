// +build !windows

package cmd

import (
	"fmt"
)

const (
	error_color = "31"
	warn_color  = "33"
	esc_start   = "\x1b["
	esc_stop    = "m"
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
	println("%s\n", error_color, args)
}

func errorf(fmt string, args ...interface{}) {
	println(fmt, error_color, args)
}

func warnln(args ...interface{}) {
	println("%s\n", warn_color, args)
}

func warnf(fmt string, args ...interface{}) {
	println(fmt, WARN_COLOR, args)
}

func println(format string, color string, args interface{}) {
	fmt.Printf("%s%s", esc_start, color)
	if format == "" {
		fmt.Println((args.([]interface{}))...)
	} else {
		fmt.Printf(format, (args.([]interface{}))...)
	}
	fmt.Printf(esc_stop)
}
