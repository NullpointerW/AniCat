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
	"github.com/NullpointerW/mikanani/download/rss"
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
	// before exit Exited channel should  be closed
	Exited chan struct{} `json:"-"`
	// While detection-gorountine detected that the resource download of the subject is completed
	// it will send downLoad message to subject-running gorountine
	// received and push to terminal
	PushChan chan qbt.Torrent `json:"-"`
	// The anime series of this project has already ended and all episodes have been downloaded.
	// while init,if this flag is false then there is noneed to start a gorountine to run it
	// exit actively flag should  be set to true
	Terminate bool `json:"terminate"`
}

type Extra struct {
	SubtitleGroup string
	RssOption     struct {
		MustContain    string
		MustNotContain string
		UseRegex       bool
	}
}

// The tag used when adding a torrent with qbt
// can be used to monitor the download status of resources
// related to this subject file.
func (s *Subject) QbtTag() string {
	return fmt.Sprintf(QbtTag, s.SubjId)
}

func (s *Subject) RssPath() string {
	return s.QbtTag()
}

func CreateSubject(n string) error {
	subject := new(Subject)

	tips, err := IC.InfoScrape(n)
	if err != nil {
		return err
	}

	sid, _ := strconv.Atoi(tips[IC.SubjId])
	if Manager.GetSubject(sid) != nil {
		return errs.Custom("%w:sid:%d", errs.ErrSubjectAlreadyExisted, sid)
	}
	subject.SubjId = sid
	err = subject.Loadfileds(tips)
	if err != nil {
		return err
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
	// create Info-Json after init completed
	subject.writeJson()

	cp := subject.Path + "/" + CoverFN
	err = CC.DOUBANCoverScraper.Scrape(cp, n)
	if err != nil {
		return err
	}

	err = download(subject)
	if err != nil {
		return err
	}

	subject.runtimeInit(false)

	return nil
}

func (subj *Subject) Loadfileds(tips map[string]string) error {
	subj.Name = tips[IC.SubjName]
	if subj.Episode, _ = strconv.Atoi(tips[IC.SubjEpisode]); subj.Episode > 1 {
		subj.Typ = TV
	} else {
		subj.Typ = MOVIE
	}

	if subj.Typ == TV {
		subj.StartTime = tips[IC.SubjStartTime]
		if et, e := tips[IC.SubjectEndTime]; e {
			n := time.Now()
			eti, err := util.ParseTime(et)
			if err != nil {
				return err
			}
			subj.EndTime = et
			subj.Finished = n.After(eti) || n.Equal(eti)
		}
	} else { // if movie finished
		subj.StartTime = tips[IC.SubjMoveStartTime]
		subj.Finished = true
	}
	return nil
}

func (s *Subject) FetchInfo() error {
	tips, err := IC.BgmTVInfoScrape(s.SubjId)
	if err != nil {
		return nil
	}
	wrap := errs.ErrWrapper{}
	wrap.Handle(func() error {
		return s.Loadfileds(tips)
	})
	wrap.Handle(func() error {
		return s.writeJson()
	})
	return wrap.Error()
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
	} else {
		err := rss.Download(qbt.AutoDLRule{
			Enabled:       true,
			AffectedFeeds: []string{subj.ResourceUrl},
			SavePath:      subj.Path,
		}, subj.RssPath())
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO use in v2
func buildRssDLR(subj *Subject, ext Extra) (DLR qbt.AutoDLRule) {
	DLR.Enabled = true
	DLR.AffectedFeeds = []string{subj.ResourceUrl}
	DLR.SavePath = subj.Path
	DLR.UseRegex = ext.RssOption.UseRegex
	DLR.MustContain = ext.RssOption.MustContain
	DLR.MustNotContain = ext.RssOption.MustNotContain
	return
}
