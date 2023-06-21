package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/errs"
	N "github.com/NullpointerW/anicat/net"
	"github.com/NullpointerW/anicat/net/cmd"
	"github.com/NullpointerW/anicat/util"
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
			util.Debugln("ls::after parse cmd", rep)
			route(&rep)
			util.Debugln("ls::after route(&rep)")
			if rep.Err != nil {
				util.Debugln("ls::in err{}")
				var s string
				if rep.Err == errs.WarnRssRuleNotMatched || rep.Err == errs.WarnReservedCommand_lsg {
					s = fmt.Sprintln(cmd.Yellow, rep.Err.Error(), cmd.Reset)
				} else {
					s = fmt.Sprintln(cmd.Red, rep.Err.Error(), cmd.Reset)
				}

				c.Write(s)
				continue
			}
			if rep.Opt == cmd.Ls || rep.Opt == cmd.LsItems || rep.Opt == cmd.LsGroup || rep.Opt == cmd.Status || rep.Opt == cmd.Stop {
				c.Write(rep.N)
				util.Debugln("ls::rep.N")
				util.Debugln(rep.N)
				if rep.Opt == cmd.Stop {
					os.Exit(0)
				}
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
