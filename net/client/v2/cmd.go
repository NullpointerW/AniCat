package main

import "strings"

type cmdtyp int

const (
	Ls = cmdtyp(iota)
	Lsi
	Text
	Stat
)

func getCmdTyp(s string) cmdtyp {
	args := strings.Fields(s)
	if len(s) < 2 {
		return Text
	}
	var cmdh string
	cmdh = args[0]
	if strings.ToLower(args[0]) == "anicat" {
		cmdh = args[1]
	}
	cmdh = strings.ToLower(cmdh)
	switch cmdh {
	case "ls":
		return Ls

	case "lsi":
		return Lsi

	default:
		return Text
	}
}

func getArg(s string) string {
	args := strings.Fields(s)
	si, ei := 1, len(args)
	cts := strings.ToLower(args[len(args)-1]) == "-s"
	if cts {
		ei = len(args) - 1
	}

	if strings.ToLower(args[0]) == "anicat" {
		si = 2
	}
	return strings.Join(args[si:ei], " ")
}
