package util

import (
	"fmt"
	"log"
)

var debug bool

func DebugEnv() {
	debug = true
}

func Debugln(v ...any) {
	if debug {
		log.Output(2, fmt.Sprintln(v...))
	}
}

func Debugf(format string, v ...any) {
	if debug {
		log.Output(2, fmt.Sprintf(format, v...))
	}
}
