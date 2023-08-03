package net

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"sync"
)

// type RssGroup struct {
// 	RssName string     `json:"rssName"`
// 	Items   []TorrItem `json:"items"`
// }

type TorrItem struct {
	Name       string `json:"name"`
	Size       string `json:"size"`
	UpdateTime string `json:"uptime"`
}

type Subj struct {
	Sid     int    `json:"sid"`
	Typ     string `json:"type"`
	Name    string `json:"name"`
	Episode string `json:"epi"`
	Status  string `json:"status"`
	Compl   string `json:"compl"`
}

type Stat struct {
	File     string `json:"file"`
	Size     string `json:"size"`
	Progress string `json:"progress"`
}

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
