package test

import (
	"fmt"
	"testing"

	"github.com/NullpointerW/mikanani/net/cmd"
)

func TestParseCmd(t *testing.T) {
	r := cmd.Parse([]string{"mikan", "add", "via", "--mc", "简中 1080"})
	fmt.Printf("%#+v \n", r)
	fmt.Println(r.Err)
}

func TestStrParse(t *testing.T) {
	str := `mikan lsi 骄傲 威威`
	cmds := cmd.ParseArgs(str)
	r := cmd.Parse(cmds)
	fmt.Printf("%#+v \n", r)
	fmt.Println(r.Err)
}

func TestPrintHelp(t *testing.T) {
	fmt.Print(cmd.TestingString())
}
