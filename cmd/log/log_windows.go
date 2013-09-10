// build windows

package log

import (
	"bytes"
	"fmt"
	"github.com/anschelsc/doscolor"
	"os"
)

const (
	info_color  = doscolor.White | doscolor.Bright
	error_color = doscolor.Red | doscolor.Bright
	warn_color  = doscolor.Yellow
)

var wrapper *doscolor.Wrapper
var padded bool

func Infof(fmt string, args ...interface{}) {
	println(fmt, info_color, false, args)
}

func Infoln(args ...interface{}) {
	println("", info_color, false, args)
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
	var buf bytes.Buffer
	if wrapper == nil {
		wrapper = doscolor.NewWrapper(os.Stdout)
	}
	wrapper.Save()
	var c doscolor.Color
	c |= color
	wrapper.Set(c)
	if pad && !padded {
		buf.WriteByte('\n')
	}
	if format == "" {
		fmt.Fprintln(&buf, (args.([]interface{}))...)
	} else {
		fmt.Fprintf(&buf, format, (args.([]interface{}))...)
	}
	wrapper.Restore()
	if pad && !padded {
		buf.WriteByte('\n')
		padded = true
	} else {
		padded = false
	}
	fmt.Printf(buf.String())
}
