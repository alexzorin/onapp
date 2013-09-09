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

var padded bool

func Infof(format string, args ...interface{}) {
	println(format, "", false, args)
}

func Infoln(args ...interface{}) {
	println("", "", false, args)
}

func InfoToggle(on bool) {
}

func Errorln(args ...interface{}) {
	out := make([]interface{}, len(args)+1)
	out[0] = "ERROR:"
	copy(out[1:], args)
	println("", error_color, true, out)
}

func Errorf(format string, args ...interface{}) {
	println(fmt.Sprintf("ERROR: %s", format), error_color, true, args)
}

func Warnln(args ...interface{}) {
	out := make([]interface{}, len(args)+1)
	out[0] = "WARNING:"
	copy(out[1:], args)
	println("", warn_color, true, out)
}

func Warnf(format string, args ...interface{}) {
	println(fmt.Sprintf("WARNING: %s", format), warn_color, true, args)
}

func println(format string, color string, pad bool, args interface{}) {
	if pad && !padded {
		fmt.Println()
	}
	fmt.Printf("%s%s", esc_start, color)
	if format == "" {
		fmt.Println((args.([]interface{}))...)
	} else {
		fmt.Printf(format, (args.([]interface{}))...)
	}
	fmt.Printf(esc_stop)
	if pad && !padded {
		fmt.Println()
		padded = true
	} else {
		padded = false
	}
}
