// build windows

package log

import (
	"fmt"
	"github.com/anschelsc/doscolor"
	"os"
)

const (
	info_color    = doscolor.White | doscolor.Bright
	error_color   = doscolor.Red | doscolor.Bright
	warn_color    = doscolor.Yellow
	success_color = doscolor.Green | doscolor.Bright
)

var wrapper *doscolor.Wrapper
var padded bool

func Infof(fmt string, args ...interface{}) {
	println(fmt, info_color, false, args)
}

func Infoln(args ...interface{}) {
	println("", info_color, false, args)
}

func Successf(fmt string, args ...interface{}) {
	println(fmt, success_color, false, args)
}

func Successln(args ...interface{}) {
	println("", success_color, false, args)
}

func InfoToggle(on bool) {
	if on {
		wrapper.Save()
		wrapper.Set(info_color)
	} else {
		wrapper.Restore()
	}
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

func println(format string, color doscolor.Color, pad bool, args interface{}) {
	if wrapper == nil {
		wrapper = doscolor.NewWrapper(os.Stdout)
	}
	wrapper.Save()
	var c doscolor.Color
	c |= color
	wrapper.Set(c)
	if pad && !padded {
		fmt.Println()
	}
	if format == "" {
		fmt.Println((args.([]interface{}))...)
	} else {
		fmt.Printf(format, (args.([]interface{}))...)
	}
	wrapper.Restore()
	if pad && !padded {
		fmt.Println()
		padded = true
	} else {
		padded = false
	}
}
