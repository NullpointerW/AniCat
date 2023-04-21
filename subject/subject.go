package subject

import (
	"fmt"
	"strconv"

	ic "github.com/NullpointerW/mikanani/crawl/information"
	"github.com/NullpointerW/mikanani/errs"
)

const (
	//resource type
	RSS = iota
	Torrent
)
const (
	//anime type
	TV = iota
	MOVIE
)

type Subject struct {
	SubjId      int    `json:"subjId"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Finished    bool   `json:"finished"`
	Episode     int    `json:"episode"`
	ResourceTyp int    `json:"resourceTyp"`
	ResourceUrl string `json:"resourceUrl"`
	Typ         int    `json:"typ"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}

func CreateSubject(n string) error {
	subject := new(Subject)

	tips, err := ic.InfoScrape(n)
	if err != nil {
		return err
	}

	sid, _ := strconv.Atoi(tips["sid"])
	if Manager.GetSubject(sid) != nil {
		return errs.Custom("subject %d already existed ", sid)
	}
	subject.SubjId = sid

	subject.Name = tips["中文名"]

	if subject.Episode, _ = strconv.Atoi(tips["话数"]); subject.Episode > 1 {
		subject.Typ = TV
	} else {
		subject.Typ = MOVIE
	}

	subject.StartTime = tips["放送开始"]

	if et, e := tips["播放结束"]; e {
		subject.EndTime = et
		subject.Finished = true
	}
	// for testing
	fmt.Printf("%#+v", *subject)
	err = solveResource(subject)
	if err != nil {
		return err
	}

	err = initFolder(subject)
	if err != nil {
		return err
	}
	Manager.Add(subject.SubjId, subject, subject.Finished)
	if subject.Finished {
		// TODO go handlerfunc()
	}

	return nil
}

func solveResource(*Subject) error {
	//TODO qb-api
	return nil
}
