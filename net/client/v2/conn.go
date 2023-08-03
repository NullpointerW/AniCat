package main

import (
	"bufio"
	"net"

	"strconv"

	N "github.com/NullpointerW/anicat/net"
)

var (
	send, recv = make(chan string), make(chan string)
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
	s.Buffer(buf, alloc<<2) // 256k
	if recv := s.Scan(); recv {

	} else {
		return nil, nil, s.Err()
	}
	c.Write([]byte("NEW_CLI" + N.CRLF))
	if recv := s.Scan(); recv {

	} else {
		return nil, nil, s.Err()
	}
	go sndrecv(s, c)
	return s, c, nil
}

func sndrecv(s *bufio.Scanner, c net.Conn) {
	for {
		sndmsg := <-recv
		c.Write([]byte(sndmsg + N.CRLF))
		if recv := s.Scan(); recv {
			send <- s.Text()
		} else {
			panic(s.Err())
		}
	}
}
