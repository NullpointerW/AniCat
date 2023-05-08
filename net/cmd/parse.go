package cmd

import (
	"github.com/NullpointerW/mikanani/errs"
	"github.com/jessevdk/go-flags"
)

var (
	GreenBg = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	RedBg   = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	Cls     = "\033[2J\033[H"
	Reset   = string([]byte{27, 91, 48, 109})
)

const (
	usageHelp = `

	Usage:
	             mikan  <command> [anine-name]
	The commands are:

	             add    add a anine-subject
	             rm     delete a anine-subject
	             ls     show all anine-subjects   
	             `
)

type Command struct {
	N    string
	Opt  Option
	Flag struct {
		Using          bool
		SubtitleGroup  string `short:"g" long:"group" required:"false"`
		MustContain    string `long:"mc"   required:"false"`
		MustNotContain string `long:"mn"   required:"false"`
		useRegex       bool   `long:"rg"  required:"false"`
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
		reply.N = cmds[2]
		if len(cmds) > 3 && opt == Add {
			ext := cmds[3:]
			_, reply.Err = flags.ParseArgs(&reply.Flag, ext)
			reply.Flag.Using = true
			return
		}
	}
	if opt == Help {
		reply.N = usageHelp
	}
	return
}
