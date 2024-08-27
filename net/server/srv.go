package server

import (
	"encoding/json"
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

type hijackConn struct {
	*N.Conn
	hijacked bool
}

func (hc *hijackConn) Write(s string) error {
	if s != "" {
		return hc.Conn.Write(s)
	}
	return nil
}

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

		go process(&hijackConn{&N.Conn{
			TcpConn: c,
		}, false})
	}
}

func process(c *hijackConn) {
	defer func() {
		if !c.hijacked {
			_ = c.TcpConn.Close()
		}
	}()
	var (
		render view.Render = view.AsciiRender{
			Conn: c.Conn,
		}
		command cmd.Cmd
	)
	if msg, err := c.Read(); err == nil {
		log.Debug(log.Struct{"len", len(msg), "cmd", msg}, "recv command")
		err := json.Unmarshal(msg, &command)
		if err != nil {
			log.Errorf(log.Struct{"error", err.Error()}, "net: json Unmarshal failed")
			_ = c.Write(err.Error())
			return
		}
		if command.Cmd == cmd.Stop {
			defer os.Exit(0)
		}
		resp, err := route(command, render)
		if err != nil {
			_ = c.Write(err.Error())
			return
		}
		_ = c.Write(resp)
	}
}
