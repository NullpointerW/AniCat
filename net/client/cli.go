package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	N "github.com/NullpointerW/mikanani/net"
	"github.com/NullpointerW/mikanani/net/cmd"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	var port string
	for {
		fmt.Println(cmd.GreenBg, "PORT:", cmd.Reset)
		l, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}
		p := string(N.DropCR([]byte(l[:len(l)-1])))
		_, err = strconv.Atoi(p)
		if err != nil {
			fmt.Println(cmd.RedBg, err, cmd.Reset)
			_, _ = r.ReadString('\n')
			fmt.Print(cmd.Cls)
			continue
		}
		port = ":" + p
		log.Print(cmd.Cls)
		break
	}

	c, err := net.Dial("tcp", port)
	if err != nil {
		log.Fatalln(cmd.RedBg, err, cmd.Reset)
	}
	s := bufio.NewScanner(c)
	s.Split(N.ScanCRLF)
	for s.Scan() {
		log.Println(s.Text())
		l, err := r.ReadString('\n')
		l = string(N.DropCR([]byte(l[:len(l)-1])))
		if err != nil {
			panic(err)
		}
		c.Write([]byte(l + N.CRLF))
	}
	log.Fatalln(cmd.RedBg, s.Err(), cmd.Reset)
}
