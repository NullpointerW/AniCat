package cmd

import (
	"encoding/json"
	"github.com/NullpointerW/anicat/net/cmd/view"
)

type cTyp int

const (
	Add cTyp = iota
	AddFeed
	Remove
	Ls
	LsItems
	Status
	Stop
	Rename
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

type CommandCase struct {
	invoke func(Cmd, view.Render) (string, error)
	flag   cTyp
}

func NewCommandCase(flag cTyp, invokeFunc func(Cmd, view.Render) (string, error)) CommandCase {
	return CommandCase{invokeFunc, flag}
}

type Selector struct {
	cases []CommandCase
}

func (sl *Selector) Select(c Cmd, r view.Render) (string, error) {
	for _, _case := range sl.cases {
		if c.Cmd == _case.flag {
			return _case.invoke(c, r)
		}
	}
	return "", nil
}

func NewSelector(cases ...CommandCase) *Selector {
	return &Selector{
		cases: cases,
	}
}
