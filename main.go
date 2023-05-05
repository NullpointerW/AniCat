package main

import (
	// "bufio"
	// "fmt"
	// "os"

	// CFG "github.com/NullpointerW/mikanani/conf"
	"os"
	"os/signal"
	"syscall"

	"github.com/NullpointerW/mikanani/download/detection"
	"github.com/NullpointerW/mikanani/subject"
)

func init() {

}

func main() {
	// fmt.Println("path: " + CFG.SubjPath)
	// fmt.Printf("host:%v \n", CFG.Proxy)
	// fmt.Printf("env:%v \n", CFG.Env)
	// r := bufio.NewReader(os.Stdin)
	// fmt.Println("请输入多行文本（以回车结束）：")
	// l, err := r.ReadString('\n')
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("read " + l)
	subject.Scan()
	go subject.StartManagement()
	go detection.Detect()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
