package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	N "github.com/NullpointerW/mikanani/net"
	"github.com/NullpointerW/mikanani/net/cmd"
)

var port int

func init() {
	flag.IntVar(&port, "p", 8080, "server dial port")
	flag.Parse()
}

func main() {
	r := bufio.NewReader(os.Stdin)
	defer func() {
		r.ReadString('\n')
	}()
	dialport := ":" + strconv.Itoa(port)
	c, err := net.Dial("tcp", dialport)
	if err != nil {
		log.Println(cmd.Red, err, cmd.Reset)
		exit(r)
	}
	s := bufio.NewScanner(c)
	s.Split(N.ScanCRLF)
	buf := make([]byte, 0, 64*1024)
	s.Buffer(buf, 1024*1024)
	for s.Scan() {
		log.Println(s.Text())
		fmt.Print(cmd.Cyan, cmd.Cursor, cmd.Reset)
		l, err := r.ReadString('\n')
		l = string(N.DropCR([]byte(l[:len(l)-1])))
		if err != nil {
			panic(err)
		}
		c.Write([]byte(l + N.CRLF))
	}
	log.Println(cmd.Red, s.Err(), cmd.Reset)
}

func exit(r *bufio.Reader) {
	r.ReadString('\n')
	os.Exit(1)
}
