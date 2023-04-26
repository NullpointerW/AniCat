package conf

import (
	"flag"
	"fmt"
	"log"
	"testing"
)

type multiValue []string

func (mv *multiValue) String() string {
	return fmt.Sprintf("%v", *mv)
}

func (mv *multiValue) Set(value string) error {
	*mv = append(*mv, value)
	return nil
}

var (
	SubjPath string
	Proxy    multiValue
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	flag.StringVar(&SubjPath, "p", "./subject", "subjects directory path")
	flag.Var(&Proxy, "h", "http proxy host")
	testing.Init()
	flag.Parse()
}