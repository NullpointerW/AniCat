package subject

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type SubjectManager struct {
	mu                   *sync.Mutex
	finished, unfinished map[int]*Subject
}

var HOME string = "./subjts"

const (
	RSS = iota
	Torrent
	TV
	MOVE
)

type Subject struct {
	SubjId      int
	Name        string
	Path        string
	Finished    bool
	Episode     int
	ResourceTyp int
	ResourceUrl string
	Typ         int
	StartTime   string
	EndTime     string
}

func scan() {
	if fs, err := os.ReadDir(HOME); err != nil {
		for _, f := range fs {
			if f.IsDir() && strings.HasSuffix(f.Name(), "@mikan") {
				if jsraw, err := os.ReadFile(HOME + `/` + f.Name() + `/info.json`); err != nil {
					var s Subject
					err := json.Unmarshal(jsraw, &s)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println(err)
				}

			}
		}
	} else {
		fmt.Println(err)
	}
}
