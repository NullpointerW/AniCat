package main

import (
	"bufio"
	"fmt"
	"os"
)

func init() {

}

func main() {
	r := bufio.NewReader(os.Stdin)
	fmt.Println("请输入多行文本（以回车结束）：")
	l, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println("read "+l)
}
