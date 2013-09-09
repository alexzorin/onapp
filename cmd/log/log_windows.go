// build windows

package log

import (
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

func Infof(fmt string, args ...interface{}) {
	println(fmt, info_color, args)
}

func Infoln(args ...interface{}) {
	println("", info_color, args)
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
	println("", error_color, args)
}

func Errorf(fmt string, args ...interface{}) {
	println(fmt, error_color, args)
}

func Warnln(args ...interface{}) {
	println("", warn_color, args)
}

func Warnf(fmt string, args ...interface{}) {
	println(fmt, warn_color, args)
}

func println(format string, color doscolor.Color, args interface{}) {
	if wrapper == nil {
		wrapper = doscolor.NewWrapper(os.Stdout)
	}
	wrapper.Save()
	var c doscolor.Color
	c |= color
	wrapper.Set(c)
	if format == "" {
		fmt.Println((args.([]interface{}))...)
	} else {
		fmt.Printf(format, (args.([]interface{}))...)
	}
	wrapper.Restore()
}
