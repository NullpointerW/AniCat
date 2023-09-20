package cmd

import "encoding/json"

type Cmd struct {
	Cmd string          `json:"cmd"`
	Raw json.RawMessage `json:"raw"`
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
