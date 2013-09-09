// +build !windows

package log

import (
	"fmt"
)

const (
	error_color = "31"
	warn_color  = "33"
	esc_start   = "\x1b["
	esc_stop    = "m"
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
	println("%s\n", error_color, args)
}

func Errorf(fmt string, args ...interface{}) {
	println(fmt, error_color, args)
}

func Warnln(args ...interface{}) {
	println("%s\n", warn_color, args)
}

func Warnf(fmt string, args ...interface{}) {
	println(fmt, warn_color, args)
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
