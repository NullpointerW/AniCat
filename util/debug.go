package util

import "log"

var debug bool

func DebugEnv() {
	debug = true
}

func Debugln(v ...any) {
	if debug {
		log.Println(v...)
	}
}

func Debugf(format string, v ...any) {
	if debug {
		log.Printf(format, v...)
	}
}
