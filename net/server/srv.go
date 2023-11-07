package server

import (
	"encoding/json"
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
	var (
		render  view.Render = view.AsciiRender{}
		command cmd.Cmd
	)
	for {
		if msg, err := c.Read(); err == nil {
			log.Debug(log.Struct{"len", len(msg), "cmd", msg}, "recv command")
			err := json.Unmarshal([]byte(msg), &command)
			if err != nil {
				log.Errorf(log.Struct{"error", err.Error()}, "net: json Unmarshal failed")
				c.Write(err.Error())
				return
			}
			route(command, render)
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
