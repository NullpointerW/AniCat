package net

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"strconv"

	CFG "github.com/NullpointerW/mikanani/conf"
)

var (
	CR   = byte('\r')
	LF   = byte('\n')
	CRLF = string([]byte{CR, LF})
)

type Conn struct {
	tcpConn net.Conn
	s    *bufio.Scanner
}

func (c Conn) Write(s string) {
	c.tcpConn.Write([]byte(s + CRLF))
}

func (c Conn) Read() {
	if c.s.Scan()
}

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

		go process(c)
	}
}

func process(c net.Conn) {
	c.Write([]byte("PONG" + CRLF))
	s := bufio.NewScanner(c)
	s.Split(ScanCRLF)
	for s.Scan() {
		cmd := s.Text()
		log.Println("cli::" + cmd)
		rp := parseCMD(cmd)
		c.Write([]byte(rp + CRLF))
	}
	if err := s.Err(); err != nil {
		log.Printf("conn closed: %s", err)
	}
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
