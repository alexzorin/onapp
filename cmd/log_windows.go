// build windows

package cmd

import (
	"fmt"
	"github.com/anschelsc/doscolor"
	"os"
)

const (
	INFO_COLOR  = doscolor.White | doscolor.Bright
	ERROR_COLOR = doscolor.Red | doscolor.Bright
	WARN_COLOR  = doscolor.Yellow
)

var wrapper *doscolor.Wrapper

func Infof(fmt string, args ...interface{}) {
	println(fmt, INFO_COLOR, args)
}

func Infoln(args ...interface{}) {
	println("", INFO_COLOR, args)
}

func InfoToggle(on bool) {
	if on {
		wrapper.Save()
		wrapper.Set(INFO_COLOR)
	} else {
		wrapper.Restore()
	}
}

func Errorln(args ...interface{}) {
	println("", ERROR_COLOR, args)
}

func Errorf(fmt string, args ...interface{}) {
	println(fmt, ERROR_COLOR, args)
}

func Warnln(args ...interface{}) {
	println("", WARN_COLOR, args)
}

func Warnf(fmt string, args ...interface{}) {
	println(fmt, WARN_COLOR, args)
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
