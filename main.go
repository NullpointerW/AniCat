package main

import (
	"bufio"
	"fmt"
	CFG "github.com/NullpointerW/mikanani/conf"
	"os"
)

func init() {

}

func main() {
	fmt.Println("path: " + CFG.SubjPath)
	fmt.Printf("host:%v \n", CFG.Proxy)
	fmt.Printf("env:%v \n", CFG.Env)
	r := bufio.NewReader(os.Stdin)
	fmt.Println("请输入多行文本（以回车结束）：")
	l, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println("read " + l)
}
