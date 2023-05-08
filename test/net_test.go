package test

import (
	"fmt"
	"testing"

	"github.com/NullpointerW/mikanani/net/cmd"
)

func TestParseCmd(t *testing.T) {
	r := cmd.Parse(`mikan  rm via -g mikangrip --mn 简中 `)
	fmt.Printf("%#+v \n", r)
	fmt.Println(r.Err)
}
