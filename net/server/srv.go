package main

import (
	"log"
	"net"
	"strconv"

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
			rep, err := cmd.Parse(msg)
			if err != nil {
				log.Printf("conn closed case inner error: %s", err)
				c.TcpConn.Close()
			}
			c.Write(rep.N)
		} else {
			log.Printf("conn closed: %s", err)
			c.TcpConn.Close()
		}
	}
}
