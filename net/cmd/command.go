package cmd

import "encoding/json"

type cTyp int

const (
	Add cTyp = iota
	AddFeed
	Remove
	Ls
	LsItems
	Status
	Stop
)

type Cmd struct {
	Cmd cTyp            `json:"cmd"`
	Arg string          `json:"arg"`
	Raw json.RawMessage `json:"raw"`
}

type AddFlag struct {
	MustContain    string `json:"mustContain"`
	MustNotContain string `json:"mustNotContain"`
	UseRegexp      bool   `json:"useRegexp"`
	Group          string `json:"group"`
	Index          int    `json:"index"`
	FeedInfoName   string `json:"feedInfoName"`
}

type LsiFlag struct {
	SearchList bool `json:"searchList"`
}
