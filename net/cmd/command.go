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

type TorrItem struct {
	Name       string `json:"name"`
	Size       string `json:"size"`
	UpdateTime string `json:"uptime"`
}

type Subj struct {
	Sid     int    `json:"sid"`
	Typ     string `json:"type"`
	Name    string `json:"name"`
	Episode string `json:"epi"`
	Status  string `json:"status"`
	Compl   string `json:"compl"`
}

type Stat struct {
	File     string `json:"file"`
	Size     string `json:"size"`
	Progress string `json:"progress"`
}
