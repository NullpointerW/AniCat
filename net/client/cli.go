package main

import (
	"bufio"
	"flag"
	"fmt"
	// "log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	N "github.com/NullpointerW/anicat/net"
	"github.com/NullpointerW/anicat/net/cmd"
)

var (
	host string
	port int
)

func init() {
	flag.StringVar(&host, "h", "localhost", "server dial host")
	flag.IntVar(&port, "p", 8080, "server dial port")
	flag.Parse()
}

func main() {
	signal := make(chan struct{})
	go waitProgress(signal)

	r := bufio.NewReader(os.Stdin)
	defer func() {
		r.ReadString('\n')
	}()
	dialadress := host + ":" + strconv.Itoa(port)
	c, err := net.Dial("tcp", dialadress)
	if err != nil {
		fmt.Println(cmd.Red, err, cmd.Reset)
		exit(r)
	}
	s := bufio.NewScanner(c)
	s.Split(N.ScanCRLF)
	buf := make([]byte, 0, 64*1024)
	s.Buffer(buf, 1024*1024)
	var f bool
	for s.Scan() {
		if f {
			signal <- struct{}{}
		}
		f = true
		fmt.Println(s.Text())
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
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/c", "cls")
			} else {
				cmd = exec.Command("clear")
			}
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
		c.Write([]byte(l + N.CRLF))
		signal <- struct{}{}
	}
	if f{
		signal <- struct{}{}
	}
	fmt.Println(cmd.Red, s.Err(), cmd.Reset)
}

func exit(r *bufio.Reader) {
	r.ReadString('\n')
	os.Exit(1)
}

func waitProgress(c chan struct{}) {
Wait:
	fmt.Print("\033[K\r")
	fmt.Print("\033[?25h")
	<-c
	st := time.Now()
	for {
		var elapsed time.Duration
		fmt.Print("\033[?25l")
		elapsed = time.Since(st)
		fmt.Printf("\\ (%0.2f s)\r", elapsed.Seconds())
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
		elapsed = time.Since(st)
		fmt.Printf("| (%0.2f s)\r", elapsed.Seconds())
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
		elapsed = time.Since(st)
		fmt.Printf("- (%0.2f s)\r", elapsed.Seconds())
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
		elapsed = time.Since(st)
		fmt.Printf("/ (%0.2f s)\r", elapsed.Seconds())
		time.Sleep(100 * time.Millisecond)
		select {
		case <-c:
			goto Wait
		default:
		}
	}
}
