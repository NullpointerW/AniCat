package cmd

import (
	"strings"

	"github.com/NullpointerW/mikanani/errs"
	"github.com/jessevdk/go-flags"
)

type Command struct {
	N    string
	Opt  Option
	Flag struct {
		Using          bool
		SubtitleGroup  string `short:"g" long:"group" required:"false"`
		MustContain    string `long:"mc"   required:"false"`
		MustNotContain string `long:"mn"   required:"false"`
		UseRegex       bool   `long:"rg"   required:"false"`
	}
	Err error
}
type Option int

const (
	Add Option = iota
	Del
	Ls
	Help
)

func optionMode(o string) (Option, bool) {
	switch o {
	case "add":
		return Add, true
	case "rm":
		return Del, true
	case "ls":
		return Ls, true
	case "h", "help", "":
		return Help, true
	default:
		return -1, false
	}
}

func Parse(cmds []string) (reply Command) {
	sfxok := cmds[0] == "mikan"
	if !sfxok {
		reply.Err = errs.Custom("%w:%s", errs.ErrUnknownCommand, cmds[0])
		return
	}
	var o string
	if len(cmds) >= 2 {
		o = cmds[1]
	}
	opt, parsed := optionMode(o)
	if !parsed {
		reply.Err = errs.Custom("%w:%s", errs.ErrUnknownCommand, o)
		return
	}
	reply.Opt = opt
	if opt != Help && opt != Ls {
		if len(cmds) < 3 {
			reply.Err = errs.Custom("%w:%s", errs.ErrMissingCommandArgument, `Use "mikan help " for more information about a command.`)
			return
		}
		if opt == Del || (opt == Add && len(cmds) == 3) {
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
		}
	}
	if opt == Help {
		reply.N = usageHelp
	}
	return
}

func ParseArgs(s string) []string {
	s = strings.ToLower(s)
	return strings.Fields(s)
}

func ParseFlagArgs(flag []string) (f []string) {
	argument := ""
	for i, a := range flag {
		arg := !(strings.HasPrefix(a, "-") || strings.HasPrefix(a, "--"))
		if arg {
			argument += " " + a
			if len(flag)-1 == i {
				f = append(f, argument)
			}
		} else {
			if argument != "" {
				f = append(f, argument)
				argument = ""
			}
			f = append(f, a)
		}
	}
	return
}
