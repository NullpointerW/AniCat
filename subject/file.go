package subject

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	CFG "github.com/NullpointerW/anicat/conf"
	DL "github.com/NullpointerW/anicat/download"
	"github.com/NullpointerW/anicat/download/rss"
	TORR "github.com/NullpointerW/anicat/download/torrent"
	"github.com/NullpointerW/anicat/errs"
	util "github.com/NullpointerW/anicat/utils"

	// "strconv"
	"strings"
)

var HOME string = CFG.Env.SubjPath

func Scan() {
	home := trimPath(HOME)
	if fs, err := os.ReadDir(home); err == nil {
		for _, f := range fs {
			if f.IsDir() {
				log.Println("scan:found folder:" + home + string(os.PathSeparator) + f.Name())
				fs, err := os.ReadDir(home + `/` + f.Name())
				if err != nil {
					log.Println(err)
					continue
				}
				for _, ff := range fs {
					isf := !ff.IsDir()
					if isf && util.IsJsonFile(ff.Name()) && strings.Contains(ff.Name(), "meta-data#") {
						if jsraw, err := os.ReadFile(home + `/` + f.Name() + `/` + ff.Name()); err == nil {
							var s Subject
							err := json.Unmarshal(jsraw, &s)
							if err != nil {
								log.Println(err)
								continue
							}
							s.runtimeInit(true)
						} else {
							log.Println(err)
						}
					}
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
	sd, err := util.ParseShort02Time(strings.ReplaceAll(subject.FolderTime, " ", ""))
	if err != nil {
		return err
	}

	folderPath += "/" + subject.FolderName + " (" + sd + ")"

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}

	if ap, err := filepath.Abs(folderPath); err == nil {
		subject.Path = util.FileSeparatorConv(ap)
	}

	return
}

func RmFolder(s *Subject) error {
	fs, err := os.ReadDir(s.Path)
	if err != nil {
		return err
	}
	var c int
	for _, f := range fs {
		isf := !f.IsDir()
		if isf && util.IsJsonFile(f.Name()) && strings.Contains(f.Name(), "meta-data#") {
			c++
		}
	}
	if c > 1 {
		return os.Remove(s.Path + "/" + fmt.Sprintf("meta-data#%s.json", s.Season))
	} else {
		return rmFolder(s)
	}
}

func rmFolder(s *Subject) error {
	return os.RemoveAll(s.Path)
}

func (s *Subject) writeJson() (err error) {
	b, _ := json.Marshal(*s)
	fldrp := s.Path
	jsfn := fmt.Sprintf(jsonfileName, s.Season)
	err = os.WriteFile(fldrp+"/"+jsfn, b, os.ModePerm)
	return err
}

func (s *Subject) RmRes() error {
	wrap := errs.ErrWrapper{}
	if s.ResourceTyp == RSS {
		wrap.Handle(func() error {
			categ := s.QbtTag()
			return DL.Qbt.RmCategoies(categ)
		})

		wrap.Handle(func() error {
			return rss.RmRss(s.RssPath())
		})
	}
	wrap.Handle(func() error {
		return TORR.DelTorrs(s.Path)
	})
	wrap.Handle(func() error {
		return TORR.DelTag(s.QbtTag())
	})
	return wrap.Error()
}

func trimPath(n string) string {
	return strings.TrimSuffix(strings.TrimSuffix(CFG.Env.SubjPath, "\\"), "/")
}

func FindLastSeason(p string) (int, error) {
	fs, err := os.ReadDir(p)
	if err != nil {
		return -1, err
	}
	var max int
	for _, f := range fs {
		isf := !f.IsDir()
		if isf && util.IsJsonFile(f.Name()) && strings.Contains(f.Name(), "meta-data#") {
			s := strings.Split(f.Name(), ".")[0]
			s = strings.ReplaceAll(s, "meta-data#", "")
			cmp, err := strconv.Atoi(s)
			if err != nil {
				return -1, err
			}
			if cmp > max {
				max = cmp
			}
		}
	}
	return max, nil
}
