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
}

func (c *Conn) Write(s string) error {
	_, err := c.TcpConn.Write([]byte(s + CRLF))
	return err
}

func (c *Conn) Read() (string, error) {
	c.once.Do(func() {
		c.s = bufio.NewScanner(c.TcpConn)
		c.s.Split(ScanCRLF)
	})
	if c.s.Scan() {
		return c.s.Text(), nil
	}
	err := c.s.Err()
	if err == nil {
		return "", io.EOF
	}
	return "", err
}

func DropCR(data []byte) []byte {
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
		return i + 2, DropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), DropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
