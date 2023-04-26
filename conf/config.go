package conf

import ()

var Env Environment

type Environment struct {
	Qbt struct {
		Host         string `json:"host"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		LocalConnect bool   `json:"localed"`
	} `json:"qbt_config"`
}
