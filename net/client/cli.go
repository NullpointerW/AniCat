package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	N "github.com/NullpointerW/mikanani/net"
)

var (
	greenBg = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	redBg   = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	cls     = "\033[2J\033[H"
	reset   = string([]byte{27, 91, 48, 109})
)

func main() {
	r := bufio.NewReader(os.Stdin)
	var p string
	for {
		fmt.Println(greenBg, "PORT:", reset)
		l, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}
		pstr := string(N.DropCR([]byte(l[:len(l)-1])))
		_, err = strconv.Atoi(pstr)
		if err != nil {
			fmt.Println(redBg, err, reset)
			_, _ = r.ReadString('\n')
			fmt.Print(cls)
			continue
		}
		p = ":" + p
		log.Print(cls)
		break
	}

	c, err := net.Dial("tcp", p)
	if err != nil {
		log.Fatalln(redBg, err, reset)
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
	log.Fatalln(redBg, s.Err(), reset)
}
