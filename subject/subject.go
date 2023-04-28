package subject

import (
	"context"
	"fmt"
	"strconv"
	"time"

	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	CC "github.com/NullpointerW/mikanani/crawl/cover"
	IC "github.com/NullpointerW/mikanani/crawl/information"
	RC "github.com/NullpointerW/mikanani/crawl/resource"
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
	// used while `ResourceTyp` is `Torrent`
	TorrentHash string `json:"torrentHash"`
	// manager use ctxcancel func to exit gorountine running the current subject.
	// when delete a subject manager,should run the cancelfunc and if a gorountine is runing
	// for this subject it will exit.
	// Context is hold by subject-running gorountine
	// while subject-running gorountine exit actively func should be called
	exit context.CancelFunc `json:"-"`
	// before detection-gorountine push to subject,Check if this channel is closed.
	// before exit exited channel should  be closed
	exited chan struct{} `json:"-"`
	// While detection-gorountine detected that the resource download of the subject is completed
	// it will send downLoad message to subject-running gorountine
	// received and push to terminal
	PushChan chan qbt.Torrent `json:"-"`
	// The anime series of this project has already ended and all episodes have been downloaded.
	// while init,if this flag is false then there is noneed to start a gorountine to run it
	// exit actively flag should  be set to true
	Terminate bool `json:"terminate"`
}

// The tag used when adding a torrent with qbt
// can be used to monitor the download status of resources
// related to this subject file.
func (s *Subject) QbtTag() string {
	return fmt.Sprintf(QbtTag, s.SubjId)
}

func CreateSubject(n string) error {
	subject := new(Subject)

	tips, err := IC.InfoScrape(n)
	if err != nil {
		return err
	}

	sid, _ := strconv.Atoi(tips[IC.SubjId])
	if Manager.GetSubject(sid) != nil {
		return errs.Custom("%w:sid: ", errs.ErrSubjectAlreadyExisted, sid)
	}
	subject.SubjId = sid

	subject.Name = tips[IC.SubjName]

	if subject.Episode, _ = strconv.Atoi(tips[IC.SubjEpisode]); subject.Episode > 1 {
		subject.Typ = TV
	} else {
		subject.Typ = MOVIE
	}

	if subject.Typ == TV {
		subject.StartTime = tips[IC.SubjStartTime]
		if et, e := tips[IC.SubjectEndTime]; e {
			n := time.Now()
			eti, err := util.ParseTime(et)
			if err != nil {
				return err
			}
			subject.EndTime = et
			subject.Finished = n.After(eti) || n.Equal(eti)
		}
	} else { // if movie finished
		subject.StartTime = tips[IC.SubjMoveStartTime]
		subject.Finished = true
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

	cp := subject.Path + "/" + CoverFN
	err = CC.DOUBANCoverScraper.Scrape(cp, n)
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
	u, isrss, err := RC.Scrape(n)
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
