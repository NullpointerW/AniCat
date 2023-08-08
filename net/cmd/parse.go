package cmd

import (
	"fmt"
	"strings"

	"github.com/NullpointerW/anicat/errs"
	"github.com/jessevdk/go-flags"
)

var cmdprfx = "anicat"

type Command struct {
	N    string
	Opt  Option
	Flag struct {
		Using          bool
		SubtitleGroup  string `short:"g" long:"group" required:"false"`
		MustContain    string `long:"mc"   required:"false"`
		MustNotContain string `long:"mn"   required:"false"`
		UseRegex       bool   `long:"rg"   required:"false"`
		Index          int    `short:"i" long:"index" required:"false"`
	}
	Err error
}
type Option int

const (
	Add Option = iota
	Del
	Ls
	LsItems
	LsItems_searchlist //lsi -s
	LsGroup            // reserved command
	Help

	Status // TODO
	Stop
)

func optionMode(o string) (Option, bool) {
	o = strings.ToLower(o)
	switch o {
	case "add":
		return Add, true
	case "rm":
		return Del, true
	case "ls":
		return Ls, true
	case "lsi":
		return LsItems, true
	case "lsg":
		return LsGroup, true
	case "stat":
		return Status, true
	case "h", "help", "":
		return Help, true
	case "stop":
		return Stop, true
	default:
		return -1, false
	}
}

func Parse(cmds []string) (reply Command) {
	sfxok := hasPrfx(cmds)
	if !sfxok {
		reply.Err = fmt.Errorf("%w:%s", errs.ErrUnknownCommand, cmds[0])
		return
	}
	var o string
	if len(cmds) >= 2 {
		o = cmds[1]
	}
	opt, parsed := optionMode(o)
	if !parsed {
		reply.Err = fmt.Errorf("%w:%s", errs.ErrUnknownCommand, o)
		return
	}
	reply.Opt = opt
	if opt != Help && opt != Ls && opt != Stop {
		if len(cmds) < 3 {
			if opt == Add {
				reply.N = addCMDUsageHelp
				reply.Err = errs.ErrAddCommandMissingHelping
				return
			}
			reply.Err = fmt.Errorf("%w:%s", errs.ErrMissingCommandArgument, `Use "(anicat) help " for more information about a command.`)
			return
		}
		if opt == Del || (opt == Add && len(cmds) == 3) || opt == Status {
			reply.N = cmds[2]
			return
		}
		if len(cmds) > 3 && opt == Add {
			ext := cmds[3:]
			e := 3
			for _, n := range ext {
				if strings.HasPrefix(n, "-") || strings.HasPrefix(n, "--") {
					break
				}
				e++
			}
			n := cmds[2:e]
			reply.N = strings.Join(n, " ")
			s := e - 3
			if s < len(ext) {
				ext = ext[s:]
				ext = ParseFlagArgs(ext)
				_, reply.Err = flags.ParseArgs(&reply.Flag, ext)
				reply.Flag.Using = true
			}
			return
		} else { // lsi , lsg  //lsi -s
			args := cmds[2:]
			ext := strings.ToLower(args[len(args)-1])
			// command is lsi ... -s,show search list 
			if ext == "-s" {
				reply.Opt = LsItems_searchlist
				args = args[:len(args)-1]
			}
			reply.N = strings.Join(args, " ")
			return
		}
	}
	if opt == Help {
		reply.N = usageHelp
	}
	return
}

func ParseArgs(s string) []string {
	// s = strings.ToLower(s)
	args := strings.Fields(s)
	if len(args) == 0 {
		return args
	}
	if hasPrfx(args) {
		return args
	} else {
		prfx := []string{cmdprfx}
		prfx = append(prfx, args...)
		return prfx
	}
}

func ParseFlagArgs(flag []string) (f []string) {
	argument := ""
	for i, a := range flag {
		arg := !(strings.HasPrefix(a, "-") || strings.HasPrefix(a, "--"))
		if arg {
			argument += " " + a
			if len(flag)-1 == i {
				argument = strings.TrimPrefix(argument, " ")
				f = append(f, argument)
			}
		} else {
			if argument != "" {
				argument = strings.TrimPrefix(argument, " ")
				f = append(f, argument)
				argument = ""
			}
			f = append(f, a)
		}
	}
	return
}

func hasPrfx(c []string) bool {
	return strings.ToLower(c[0]) == cmdprfx
}
