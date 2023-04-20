package conf

import "flag"

var (
	SubjPath string
	Proxy    string
)

func init() {
	flag.StringVar(&SubjPath, "p", "./subject", "subjects directory path")
	flag.StringVar(&Proxy, "h", "", "http proxy host")
}
