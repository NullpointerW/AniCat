package server

import (
	"fmt"
	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	N "github.com/NullpointerW/anicat/net"
	"github.com/NullpointerW/anicat/net/cmd"
	"github.com/NullpointerW/anicat/net/cmd/view"
	"net"
	"os"
	"runtime"
	"strconv"
)

func Listen() {
	p := CFG.Env.Port
	if p == 0 {
		p = 12314
	}
	adr := ":" + strconv.Itoa(p)
	ls, err := net.Listen("tcp", adr)
	errCallbackFunc := func() {
		if runtime.GOOS == "windows" {
			log.Error(log.Struct{"err", err}, "PANIC! process crashed")
		}
	}
	if err != nil {
		errs.PanicErr(err, errCallbackFunc)
	}
	for {
		c, err := ls.Accept()
		if err != nil {
			log.Error(log.Struct{"err", err}, "accept connection failed")
			continue
		}

		go process(&N.Conn{
			TcpConn: c,
		})
	}
}

func process(c *N.Conn) {
	c.Write("PONG")
	var (
		fsMsg              = true
		render view.Render = view.AsciiRender{}
	)
	for {
		if msg, err := c.Read(); err == nil {
			if fsMsg && msg == "NEW_CLI" {
				c.Write("RECV_CLI_VER")
				render = view.JsonRender{}
				continue
			}
			// old cli
			fsMsg = false
			log.Debug(log.Struct{"len", len(msg), "cmd", msg}, "recv command")
			if len(msg) == 0 {
				c.Write("PONG")
				continue
			}
			cmds := cmd.ParseArgs(msg)
			if len(cmds) == 0 {
				c.Write("PONG")
				continue
			}
			rep := cmd.Parse(cmds)
			if rep.Err != nil {
				s := ""
				if rep.Err == errs.ErrAddCommandMissingHelping {
					s = rep.N
				} else {
					s = fmt.Sprintln(cmd.Red, rep.Err.Error(), cmd.Reset)
				}
				c.Write(s)
				continue
			}
			if rep.Opt == cmd.Help {
				c.Write(rep.N)
				continue
			}
			route(&rep, render)
			if rep.Err != nil {
				var s string
				if rep.Err == errs.WarnRssRuleNotMatched || rep.Err == errs.WarnReservedCommand_lsg {
					s = fmt.Sprintln(cmd.Yellow, rep.Err.Error(), cmd.Reset)
				} else {
					s = fmt.Sprintln(cmd.Red, rep.Err.Error(), cmd.Reset)
				}

				c.Write(s)
				continue
			}
			if rep.Opt == cmd.Ls || rep.Opt == cmd.LsItems ||
				rep.Opt == cmd.LsItems_searchlist || rep.Opt == cmd.LsGroup ||
				rep.Opt == cmd.Status || rep.Opt == cmd.Stop || rep.Opt == cmd.Add || rep.Opt == cmd.AddFeed {
				c.Write(rep.N)
				if rep.Opt == cmd.Stop {
					os.Exit(0)
				}
			} else {
				c.Write("OK")
			}
		} else {
			c.TcpConn.Close()
			log.Error(log.Struct{"err", err}, "connection closed")
			break
		}
	}
}
