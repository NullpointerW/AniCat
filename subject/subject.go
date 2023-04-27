package subject

import (
	"fmt"
	"strconv"
	"time"

	ic "github.com/NullpointerW/mikanani/crawl/information"
	rc "github.com/NullpointerW/mikanani/crawl/resource"
	"github.com/NullpointerW/mikanani/errs"
	"github.com/NullpointerW/mikanani/util"
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

	sid, _ := strconv.Atoi(tips[ic.SubjId])
	if Manager.GetSubject(sid) != nil {
		return errs.Custom("%w:sid: ", errs.ErrSubjectAlreadyExisted, sid)
	}
	subject.SubjId = sid

	subject.Name = tips[ic.SubjName]

	if subject.Episode, _ = strconv.Atoi(tips[ic.SubjEpisode]); subject.Episode > 1 {
		subject.Typ = TV
	} else {
		subject.Typ = MOVIE
	}

	subject.StartTime = tips[ic.SubjStartTime]
	if et, e := tips[ic.SubjectEndTime]; e {
		n := time.Now()
		eti, err := util.ParseTime(et)
		if err != nil {
			return err
		}
		subject.EndTime = et
		subject.Finished = n.After(eti) || n.Equal(eti)
	}

	// for testing
	fmt.Printf("%#+v", *subject)

	err = solveResource(n, subject)
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

func solveResource(n string, subj *Subject) error {
	rc.Scrape(n)
	return nil
}
