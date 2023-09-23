package cmd

import "encoding/json"

type Cmd struct {
	Cmd string          `json:"cmd"`
	Arg string          `json:"arg"`
	Raw json.RawMessage `json:"raw"`
}

type Add_ struct {
	MustContain    string `json:"mustContain"`
	MustNotContain string `json:"mustNotContain"`
	UseRegexp      bool   `json:"useRegexp"`
	Group          string `json:"group"`
	Feed           string `json:"feed"`
	FeedName       string `json:"feedName"`
	Index          string `json:"index"`
}

type Lsi_ struct {
	SearchList bool `json:"searchList"`
}
