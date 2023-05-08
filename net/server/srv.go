package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	// CFG "github.com/NullpointerW/mikanani/conf"
	N "github.com/NullpointerW/mikanani/net"
	"github.com/NullpointerW/mikanani/net/cmd"
)

func Listen() {
	p := 8007
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

func main() {
	Listen()
}

func process(c *N.Conn) {
	c.Write("PONG")
	for {
		if msg, err := c.Read(); err == nil {
			log.Printf("msg_len:%d _:%s \n", len(msg), msg)
			if len(msg) == 0 {
				c.Write("PONG")
				continue
			}
			msg = strings.ToLower(msg)
			cmds := strings.Fields(msg)
			if len(cmds) == 0 {
				c.Write("PONG")
				continue
			}
			rep := cmd.Parse(cmds)
			if rep.Err != nil {
				s := fmt.Sprintln(cmd.RedBg, rep.Err.Error(), cmd.Reset)
				c.Write(s)
				continue
			}
			if rep.Opt == cmd.Help {
				c.Write(rep.N)
				continue
			}
			c.Write("ok")
		} else {
			log.Printf("conn closed: %s", err)
			c.TcpConn.Close()
		}
	}
}
