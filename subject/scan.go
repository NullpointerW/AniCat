package subject

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var HOME string = "./subjts"

func Scan() {
	if fs, err := os.ReadDir(HOME);err == nil {
		for _, f := range fs {
			if f.IsDir() && strings.HasSuffix(f.Name(), "@mikan") {
				if jsraw, err := os.ReadFile(HOME + `/` + f.Name() + `/info.json`); err != nil {
					var s Subject
					err := json.Unmarshal(jsraw, &s)
					if err != nil {
						fmt.Println(err)
					}
					if s.Finished {
						Manager.Add(s.SubjId, &s, true)
					} else {
						Manager.Add(s.SubjId, &s, false)
						go func(s *Subject) {}(&s)
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
