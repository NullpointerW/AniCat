package subject

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	CFG "github.com/NullpointerW/anicat/conf"
	DL "github.com/NullpointerW/anicat/download"
	"github.com/NullpointerW/anicat/download/rss"
	TORR "github.com/NullpointerW/anicat/download/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/util"

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
	sd, err := util.ParseShortTime(subject.StartTime)
	if err != nil {
		return err
	}

	floderN := strings.ReplaceAll(util.FileSeparatorConv(subject.Name), "/", " ")
	folderPath += "/" + floderN + " (" + sd + ")"

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}

	if ap, err := filepath.Abs(folderPath); err == nil {
		subject.Path = util.FileSeparatorConv(ap)
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
