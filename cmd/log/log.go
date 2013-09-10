// +build !windows

package log

import (
	"bytes"
	"fmt"
)

const (
	error_color = "31"
	warn_color  = "33"
	esc_start   = "\x1b["
	esc_stop    = "m"
	esc_default = "0"
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
	var buf bytes.Buffer
	if pad && !padded {
		buf.WriteByte('\n')
	}
	fmt.Fprintf(&buf, "%s%s%s", esc_start, color, esc_stop)
	if format == "" {
		fmt.Fprintln(&buf, (args.([]interface{}))...)
	} else {
		fmt.Fprintf(&buf, format, (args.([]interface{}))...)
	}
	fmt.Fprintf(&buf, "%s%s%s", esc_start, esc_default, esc_stop)
	if pad && !padded {
		buf.WriteByte('\n')
		padded = true
	} else {
		padded = false
	}
	fmt.Printf(buf.String())
}
