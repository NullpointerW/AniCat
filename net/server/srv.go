package server

import (
	"fmt"
	"log"
	"net"
	"strconv"

	CFG "github.com/NullpointerW/mikanani/conf"
	"github.com/NullpointerW/mikanani/errs"
	N "github.com/NullpointerW/mikanani/net"
	"github.com/NullpointerW/mikanani/net/cmd"
	"github.com/NullpointerW/mikanani/util"
)

func Listen() {
	p := CFG.Env.Port
	if p == 0 {
		p = 8080
	}
	adr := ":" + strconv.Itoa(p)
	ls, err := net.Listen("tcp", adr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := ls.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go process(&N.Conn{
			TcpConn: c,
		})
	}
}

func process(c *N.Conn) {
	c.Write("PONG")
	for {
		if msg, err := c.Read(); err == nil {
			util.Debugf("msg_len:%d cmd:%s \n", len(msg), msg)
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
			route(&rep)
			if rep.Err != nil {
				var s string
				if rep.Err == errs.WarnRssRuleNotMatched || rep.Err == errs.WarnReservedCommand_lsg {
					s = fmt.Sprintln(cmd.YellowBg, rep.Err.Error(), cmd.Reset)
				} else {
					s = fmt.Sprintln(cmd.Red, rep.Err.Error(), cmd.Reset)
				}

				c.Write(s)
				continue
			}
			if rep.Opt == cmd.Ls || rep.Opt == cmd.LsItems || rep.Opt == cmd.LsGroup || rep.Opt == cmd.Status {
				c.Write(rep.N)
			} else {
				c.Write("OK")
			}
		} else {
			log.Printf("conn closed: %s", err)
			c.TcpConn.Close()
			break
		}
	}
}
