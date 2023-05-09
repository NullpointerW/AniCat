package subject

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	CFG "github.com/NullpointerW/mikanani/conf"
)

var HOME string = CFG.SubjPath

func Scan() {
	home := trimPath(HOME)
	if fs, err := os.ReadDir(home); err == nil {
		for _, f := range fs {
			if f.IsDir() && strings.HasSuffix(f.Name(), FolderSuffix) {
				log.Println("scan:found subj:" + home + `/` + f.Name() + `/` + jsonfileName)
				if jsraw, err := os.ReadFile(home + `/` + f.Name() + `/` + jsonfileName); err == nil {
					var s Subject
					err := json.Unmarshal(jsraw, &s)
					if err != nil {
						log.Println(err)
					}
					s.runtimeInit(true)
				} else {
					log.Println(err)
				}
			}
		}
	} else {
		log.Println(err)
	}
}

// Initialize the content library in OS file for the `subject` in the memory.
// Path can be used to monitor the download status of resources
// apply to RSS and Torrent type.
// If initialization is successful, write the path to Subject.Path.
func initFolder(subject *Subject) (err error) {
	var folderPath string

	folderPath = trimPath(HOME)
	folderPath += "/" + strconv.Itoa(subject.SubjId) + FolderSuffix

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}

	if ap, err := filepath.Abs(folderPath); err == nil {
		subject.Path = ap
	}

	return
}

func rmFolder(s *Subject) error {
	return os.RemoveAll(s.Path)
}

func (s *Subject) writeJson() (err error) {
	b, _ := json.Marshal(*s)
	fldrp := s.Path
	err = os.WriteFile(fldrp+"/"+jsonfileName, b, os.ModePerm)
	return err
}

func trimPath(n string) string {
	return strings.TrimSuffix(strings.TrimSuffix(CFG.SubjPath, "\\"), "/")
}
