package subject

import (
	"context"
	"fmt"
	"strconv"
	"time"

	ic "github.com/NullpointerW/mikanani/crawl/information"
	rc "github.com/NullpointerW/mikanani/crawl/resource"
	detn "github.com/NullpointerW/mikanani/download/detection"
	"github.com/NullpointerW/mikanani/download/torrent"
	"github.com/NullpointerW/mikanani/errs"
	"github.com/NullpointerW/mikanani/util"
)

// Subject as basic obj of each bgmi
// programm will load it from OS file and manage thme
// while some of fileds have updating it will be refresh to OS file
type Subject struct {
	SubjId      int         `json:"subjId"`
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Finished    bool        `json:"finished"`
	Episode     int         `json:"episode"`
	ResourceTyp ResourceTyp `json:"resourceTyp"`
	ResourceUrl string      `json:"resourceUrl"`
	Typ         BgmiTyp     `json:"typ"`
	StartTime   string      `json:"startTime"`
	EndTime     string      `json:"endTime"`
	// used when `ResourceTyp` is `Torrent`
	TorrentHash string `json:"torrentHash"`
	// manager use ctxcancel func to exit gorountine running the current subject.
	// when delete a subject manager,should run the cancelfunc and if a gorountine is runing
	// for this subject it will exit.
	// Context is hold by subject-running gorountine
	// while subject-running gorountine exit actively exit should be called
	exit context.CancelFunc
	// While detection-gorountine detected that the resource download of the subject is completed
	// it will send downLoad message to subject-running gorountine
	// received and push to terminal
	PushChan chan detn.BgmiDLInfo
}

// The tag used when adding a torrent with qbt
// can be used to monitor the download status of resources
// related to this subject file.
func (s *Subject) QbtTag() string {
	return fmt.Sprintf(QbtTag, s.SubjId)
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

	err = download(subject)
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
	u, isrss, err := rc.Scrape(n)
	if err != nil {
		return err
	}
	subj.ResourceUrl = u
	if isrss {
		subj.ResourceTyp = RSS
	} else {
		subj.ResourceTyp = Torrent
	}
	return nil
}

func download(subj *Subject) error {
	if subj.ResourceTyp == Torrent {
		h, err := torrent.Add(subj.ResourceUrl, subj.Path, subj.QbtTag())
		if err != nil {
			return err
		}
		subj.TorrentHash = h
	}
	return nil
}
