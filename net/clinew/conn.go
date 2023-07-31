package main

import (
	"bufio"
	"net"

	"strconv"

	N "github.com/NullpointerW/anicat/net"
)

func connect() (*bufio.Scanner, net.Conn, error) {
	dialadress := host + ":" + strconv.Itoa(port)
	c, err := net.Dial("tcp", dialadress)
	if err != nil {
		return nil, nil, err
	}
	s := bufio.NewScanner(c)
	s.Split(N.ScanCRLF)
	alloc := 64 * 1024 // 64k
	buf := make([]byte, 0, alloc)
	s.Buffer(buf, 10<<alloc) // 1024k
	if recv := s.Scan(); recv {
		
	} else {
		return nil, nil, s.Err()
	}
	c.Write([]byte("NEW_CLI" + N.CRLF))
	return s, c, nil
}
