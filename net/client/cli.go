package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	N "github.com/NullpointerW/anicat/net"
	"github.com/NullpointerW/anicat/net/cmd"
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
		if s.Text() == "exited." {
			return
		}
		var (
			err error
			l   string
		)
		for {
			fmt.Print(cmd.Cyan, cmd.Cursor, cmd.Reset)
			l, err = r.ReadString('\n')
			if err != nil {
				panic(err)
			}
			l = string(N.DropCR([]byte(l[:len(l)-1])))
			if l != "cls" && l != "clear" {
				break
			}
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
		c.Write([]byte(l + N.CRLF))
	}
	log.Println(cmd.Red, s.Err(), cmd.Reset)
}

func exit(r *bufio.Reader) {
	r.ReadString('\n')
	os.Exit(1)
}

func waitProgress(c chan struct{}) {
	fmt.Print("\033[K\r")
	fmt.Print("\033[?25h")
Wait:
	<-c
	for {
		fmt.Print("\033[?25l")
		fmt.Printf("\\\r")
		select {
		case <-c:
			goto Wait
		default:
		}
		time.Sleep(100 * time.Millisecond)
		select {
		case <-c:
			goto Wait
		default:
		}
		fmt.Printf("|\r")
		select {
		case <-c:
			goto Wait
		default:
		}
		time.Sleep(100 * time.Millisecond)
		select {
		case <-c:
			goto Wait
		default:
		}
		fmt.Printf("-\r")
		select {
		case <-c:
			goto Wait
		default:
		}
		time.Sleep(100 * time.Millisecond)
		select {
		case <-c:
			goto Wait
		default:
		}
		fmt.Printf("/\r")
		time.Sleep(100 * time.Millisecond)
		select {
		case <-c:
			goto Wait
		default:
		}
	}
}
