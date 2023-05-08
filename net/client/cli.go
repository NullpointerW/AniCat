package main

import (
	"bufio"
	N "github.com/NullpointerW/mikanani/net"
	"log"
	"net"
	"os"
)

func main() {
	c, err := net.Dial("tcp", ":8007")
	if err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(c)
	r := bufio.NewReader(os.Stdin)
	s.Split(N.ScanCRLF)
	for s.Scan() {
		log.Println(s.Text())
		l, err := r.ReadString('\n')
		l = l[:len(l)-1]
		l = string(N.DropCR([]byte(l)))
		if err != nil {
			panic(err)
		}
		c.Write([]byte(l + N.CRLF))
	}
	log.Fatal(s.Err())
}
