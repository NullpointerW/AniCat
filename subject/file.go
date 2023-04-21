package subject

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	CFG "github.com/NullpointerW/mikanani/conf"
)

const (
	folderSuffix = "@mikan"
	jsonfileName = "info.json"
)

var HOME string = CFG.SubjPath

func Scan() {
	home := trimPath(HOME)
	if fs, err := os.ReadDir(home); err == nil {
		for _, f := range fs {
			if f.IsDir() && strings.HasSuffix(f.Name(), folderSuffix) {
				if jsraw, err := os.ReadFile(home + `/` + f.Name() + `/` + jsonfileName); err != nil {
					var s Subject
					err := json.Unmarshal(jsraw, &s)
					if err != nil {
						fmt.Println(err)
					}
					Manager.Add(s.SubjId, &s, s.Finished)
					if !s.Finished {
						// TODO handler
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

func initFolder(subject *Subject) error {
	var folderPath string

	folderPath = trimPath(HOME)
	folderPath += "/" + strconv.Itoa(subject.SubjId) + folderSuffix

	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}

	if ap, err := filepath.Abs(folderPath); err == nil {
		subject.Path = ap
	} else {
		return err
	}

	b, _ := json.Marshal(*subject)
	err = os.WriteFile(folderPath+"/"+jsonfileName, b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func trimPath(n string) string {
	return strings.TrimSuffix(strings.TrimSuffix(CFG.SubjPath, "\\"), "/")
}
