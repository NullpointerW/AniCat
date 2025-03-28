package net

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"sync"
)

var (
	CR   = byte('\r')
	LF   = byte('\n')
	CRLF = string([]byte{CR, LF})
)

type Conn struct {
	TcpConn net.Conn
	s       *bufio.Scanner
	once    sync.Once
	Max     int
	Hajacked bool
}

func (c *Conn) Write(s string) error {
	_, err := c.TcpConn.Write([]byte(s + CRLF))
	return err
}

func (c *Conn) Read() ([]byte, error) {
	c.once.Do(func() {
		c.s = bufio.NewScanner(c.TcpConn)
		if c.Max > 0 {
			c.s.Buffer(make([]byte, 0, c.Max), c.Max)
		}
		c.s.Split(ScanCRLF)
	})
	if c.s.Scan() {
		return c.s.Bytes(), nil
	}
	err := c.s.Err()
	if err == nil {
		return nil, io.EOF
	}
	return nil, err
}

func DropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == CR {
		return data[0 : len(data)-1]
	}
	return data
}

func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{CR, LF}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, DropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), DropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
