package subject

import (
	"context"
	"fmt"
	"net/url"
	// "reflect"
	"strings"
	"time"

	CFG "github.com/NullpointerW/anicat/conf"
	DL "github.com/NullpointerW/anicat/downloader"
	"github.com/NullpointerW/anicat/downloader/builtin"
	"github.com/NullpointerW/anicat/downloader/rss"
	TORR "github.com/NullpointerW/anicat/downloader/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	eslog "github.com/NullpointerW/anicat/pkg/log"
	P "github.com/NullpointerW/anicat/pusher"
	"github.com/NullpointerW/anicat/pusher/email"
	util "github.com/NullpointerW/anicat/utils"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

// runtimeInit before goroutine handle it init inner channels and ctx func
func (s *Subject) runtimeInit(reload bool) {
	s.Exited = make(chan struct{})
	if s.Terminate {
		close(s.Exited)
		Mgr.Add(s)
		if s.BuiltinDownload{
			s.initializeFinishedTorrentNameList()
		}
		return
	}
	c := context.Background()
	ctx, exit := context.WithCancel(c)
	s.Exit = exit
	if s.Pushed == nil {
		s.Pushed = make(map[string]string)
	}
	s.OperationChan = make(chan Operate)
	Mgr.Add(s)
	if s.BuiltinDownload {
		if s.TorrentUrls == nil {
			s.TorrentUrls = make(map[string]RssFileOptStrage)
		}
		if s.TorrentFinishedUrls == nil {
			s.TorrentFinishedUrls = make(map[string]struct{})
		}
		if reload {
			s.initializeFinishedTorrentNameList()
		} else {
			s.FinihsedTorrentNameList = util.NewListView([]builtin.TorrentProgress(nil))
		}
		s.PushChanBuiltin = make(chan builtin.MonitoredTorrent, 1024)
		s.DetctchanBuiltin = make(chan builtin.MonitoredTorrent, 1024)
		m:=builtin.NewTorrentProgressMonitor(time.Second*15)
		go builtin.DetectBuiltin(s.DetctchanBuiltin, s.PushChanBuiltin, ctx,m)
		go s.runWithBuiltinDownloader(ctx, reload)

	} else {
		s.PushChan = make(chan qbt.Torrent, 1024)
		go s.run(ctx, reload)
	}
	go JellyfinMetaDataHelper(s.Path, s.FolderName, s.Exited)
}

// runWithBuiltinDownloader handles the downloading process for a given subject
// using built-in downloaders. It supports both torrent and RSS resource types.
// The function performs the following tasks:
//  1. Adds a wait group counter for the download process.
//  2. Logs a debug message if the subject is set to reload.
//  3. Sets up a ticker to periodically check for updates.
//  4. For torrent resources, it initializes the appropriate torrent seeker based
//     on the URL scheme and starts the download process using the built-in downloader.
//  5. For RSS resources, it initializes an RSS reader and resumes any pending downloads.
//  6. Enters a loop to handle various operations such as renaming, pushing updates,
//     and periodic updates based on the ticker.
//  7. Exits gracefully when the context is done, logging the exit event.
func (s *Subject) runWithBuiltinDownloader(ctx context.Context, reload bool) {
	Mgr.wg.Add(1)
	defer Mgr.wg.Done()
	if reload {
		log.Debug(log.Struct{"sid", s.SubjId}, "subject reload")
		//s.checkDL()
	}
	t := time.NewTicker(30 * time.Minute)
	if s.ResourceTyp == Torrent {
		var seeker builtin.TorrentSeeker
		u, err := url.Parse(s.ResourceUrl)
		if err != nil {
			log.Error(log.Struct{"err", err}, "parse torrentUrl failed")
			s.Exit()
		}
		scheme := strings.ToLower(u.Scheme)
		switch {
		case scheme == "magnet":
			seeker = &MagnetUrlSeeker{}
		case scheme == "http" || scheme == "https":
			seeker = nil
		default:
			log.Error(log.Struct{"err", fmt.Errorf("unexpected scheme %q", u.Scheme)}, "parse torrentUrl failed")
			s.Exit()
		}
		var fop builtin.FileOption
		switch s.Typ {
		case TV:
			fop = FilePath{FileName: &TorrFileOpt{s}, DirPath: s.Path}
		case MOVIE:
			fop = FilePath{FileName: new(MovieFileOpt), DirPath: s.Path}
		}
		fmt.Printf("TorrFileOpt: %+v \n", fop)
		t, err := builtin.DefaultDownLoader.Download(s.ResourceUrl, fop, seeker)
		if err != nil {
			log.Error(log.Struct{"err", err}, "download torrentResource failed")
			s.Exit()
		}
		s.builtinDownload(builtin.MonitoredTorrent{TorrentInfo: builtin.TorrentInfo{Torrent: t}, Url: s.ResourceUrl})
	}
	if s.ResourceTyp == RSS && reload {
		var ff rss.FilterFunc
		if s.Filter != nil {
			ff = s.Filter.Filter()
		}
		s.RssReader = rss.NewReader(s.ResourceUrl, s.RssGuids, ff)
	}
	err := s.resumeRssDownload()
	if err != nil {
		log.Error(log.Struct{"err", err, "subj", s.SubjId}, "resume download failed")
		s.Exit()
	}
	for {
		select {
		case o := <-s.OperationChan:
			switch o.op {
			case Rename:
				err := s.Rename(o.arg.(string))
				if err != nil {
					log.Error(log.Struct{"sid", s.SubjId, "err", err}, "rename failed")
				}
			}
		case torr := <-s.PushChanBuiltin:
			if s.ResourceTyp == Torrent {
				s.TorrentFinishedUrls[torr.Rename] = struct{}{}
			} else {
				s.TorrentFinishedUrls[torr.Url] = struct{}{}
			}
			s.FinihsedTorrentNameList.Append(builtin.TorrentProgress{Percentage: 100, Name: torr.Rename})
			err := s.writeJson()
			if err != nil {
				log.Error(log.Struct{"sid", s.SubjId, "err", err}, "write json failed")
			}
			err = s.pushBuiltin(torr, email.Poster)
			if err != nil {
				log.Error(log.Struct{"sid", s.SubjId, "err", err}, "push process failed")
			}
		case <-ctx.Done():
			log.Debug(log.Struct{"sid", s.SubjId}, "runner exited")
			exit(s)
			return
		case <-t.C:
			log.Debug(log.Struct{"sid", s.SubjId}, "subject update mission started")
			s.readRssAndDownload()
			err := s.update()
			if err != nil {
				log.Error(log.Struct{"sid", s.SubjId, "err", err}, "update mission failed")
			}

		}
	}
}

func (s *Subject) run(ctx context.Context, reload bool) {
	Mgr.wg.Add(1)
	defer Mgr.wg.Done()
	if reload {
		log.Debug(log.Struct{"sid", s.SubjId}, "subject reload")
		//s.checkDL()
	}
	t := time.NewTicker(30 * time.Minute)
	for {
		select {
		case o := <-s.OperationChan:
			switch o.op {
			case Rename:
				err := s.Rename(o.arg.(string))
				if err != nil {
					log.Error(log.Struct{"sid", s.SubjId, "err", err}, "rename failed")
				}
			}
		case torr := <-s.PushChan:
			err := s.push(torr, email.Poster)
			if err != nil {
				log.Error(log.Struct{"sid", s.SubjId, "err", err}, "push process failed")
			}
		case <-ctx.Done():
			log.Debug(log.Struct{"sid", s.SubjId}, "runner exited")
			exit(s)
			return
		case <-t.C:
			log.Debug(log.Struct{"sid", s.SubjId}, "subject update mission started")
			err := s.update()
			if err != nil {
				log.Error(log.Struct{"sid", s.SubjId, "err", err}, "update mission failed")
			}
		}
	}

}

func (s *Subject) update() error {
	wrap := errs.ErrWrapper{}
	wrap.Handle(func() error {
		return s.FetchInfo()
	})
	wrap.Handle(func() error {
		return s.writeJson()
	})
	wrap.Handle(func() error {
		if s.BuiltinDownload {
			return s.checkDLWithBuiltin()
		}
		return s.checkDL()
	})
	return wrap.Error()
}

func exit(s *Subject) {
	err := s.writeJson()
	if err != nil {
		log.Error(log.Struct{"sid", s.SubjId, "err", err}, "write json failed while exited")
	}
	close(s.Exited)
	if !s.BuiltinDownload {
		close(s.PushChan)
	}
	close(s.OperationChan)
	Mgr.Sync()
}

func (s *Subject) checkDLWithBuiltin() (err error) {
	if s.ResourceTyp == Torrent {
		if len(s.TorrentFinishedUrls) > 0 {
			log.Info(log.Struct{"sname", s.Name, "resType", "Torrent"}, "compiled,exited now")
			s.terminate()
		}
	} else {
		fin, err := s.ElapsedfromFinishedTime(util.Day * 2)
		if err != nil {
			return fmt.Errorf("checkDLWithBuiltin: %w", err)
		}
		if fin && len(s.TorrentFinishedUrls) == len(s.TorrentUrls) {
			log.Info(log.Struct{"sname", s.Name, "resType", "RSS"}, "compiled,exited now")
			s.terminate()
		}
	}
	return nil
}
func (s *Subject) checkDL() (err error) {
	if s.ResourceTyp == Torrent {
		log.Debug(log.Struct{"sid", s.SubjId, "type", "torrent"}, "start check DL")
		compl, err := TORR.DLcompl(s.TorrentHash)
		if err != nil {
			return err
		} else if compl {
			log.Debug(log.Struct{"sid", s.SubjId, "type", "torrent"}, "DL fin terminate now")
			s.terminate()
			return err
		}
	} else if s.ResourceTyp == RSS && s.Finished {
		log.Debug(log.Struct{"sid", s.SubjId, "type", "rss"}, "start check DL")
		if s.Typ == TV && s.EndTime != "" {
			log.Debug(log.Struct{"sid", s.SubjId, "resType", "TV"}, "epi fin")
			e, err := util.ParseTime(s.EndTime, util.YMDParseLayout)
			log.Debug(log.Struct{"sid", s.SubjId, "resType", "TV"}, "epi endtime is ", util.ParseTimeStr(e))
			if err != nil {
				return err
			}
			if time.Since(e) >= util.Day {
				log.Debug(log.Struct{"sid", s.SubjId, "resType", "TV"}, "The time elapsed since the end of the anime is more than 1 day. ")
				goto checkSync
			}
			log.Debug(log.Struct{"sid", s.SubjId, "resType", "TV"}, "The time elapsed DAY between the end of the anime and nowtime is ",
				time.Since(e).Hours()/24)
		} else if s.Typ == MOVIE {
			goto checkSync
		}
	}
	return
checkSync:
	sync, err := s.RssDLSynced()
	if err != nil {
		return err
	}
	if sync {
		s.terminate()
	}
	return
}

// RssDLSynced called only when the subject epis is fin
func (s *Subject) RssDLSynced() (bool, error) {
	arts, err := rss.GetMatchedArts(s.RssPath())
	if err != nil {
		return false, nil
	}
	tlen := len(arts)
	if tlen == 0 {
		log.Warn(log.Struct{"sid", s.SubjId, "resType", "RSS"}, "there is no arts matched,check the rss match rule!")
		return true, nil
	}
	log.Debug(log.Struct{"sid", s.SubjId, "rssTotalLen", tlen})
	c := len(s.RssTorrents)
	log.Debug(log.Struct{"sid", s.SubjId, "series", tlen, "localSeries", c, "cmplSeries", c})
	return c >= tlen, nil
}

func (s *Subject) push(torr qbt.Torrent, pusher P.Pusher) error {
	if s.ResourceTyp == Torrent {
		if torr.Hash == s.TorrentHash {
			var err error
			if s.Typ == TV {
				err = renameTorr(s, torr)
				if err != nil {
					goto term
				}
				epi := "S" + s.Season + "E01-" + fmt.Sprintf("%02d", s.Episode)
				err = pusher.Push(P.Payload{
					SubjectId:    s.SubjId,
					SubjectName:  s.Name,
					DownLoadName: torr.Name,
					Size:         torr.Size,
					Episode:      epi,
				})
			} else {
				err = pusher.Push(P.Payload{
					SubjectId:    s.SubjId,
					SubjectName:  s.Name,
					DownLoadName: torr.Name,
					Size:         torr.Size,
					Episode:      "MOVIE",
				})
			}
		term:
			s.terminate()
			return err
		}
		return nil
	}
	// RSS
	if s.Pushed == nil {
		s.Pushed = make(map[string]string)
	}
	if s.RssTorrents == nil {
		s.RssTorrents = map[string]struct{}{}
	}

	// perf: skip rename process
	if _, e := s.RssTorrents[torr.Hash]; e {
		log.Debug(log.Struct{"sid", s.SubjId, "torrfn", torr.Name, "torrHash", torr.Hash}, "skip rename")
		return nil
	}

	s.RssTorrents[torr.Hash] = struct{}{}
	if s.Typ == TV {
		var se = ""
		if checkSingleVideo(torr) {
			rename, err := renameTV(s, torr.Name)
			if err != nil {
				if CFG.Env.BgmiLog {
					CFG.BgmiLogger.Infof(eslog.Struct{"sid", s.SubjId, "name", s.Name}, "episode update(unnamed): %s", torr.Name)
				}
				return err
			}
			se = util.TrimExtensionAndGetEpi(rename)
			if th, e := s.Pushed[se]; e {
				merr := errs.MultiErr{}
				dumpliErr := fmt.Errorf("%w: origin_name=%s,rename=%s", errs.ErrItemAlreadyPushed, torr.Name, rename)
				merr.Add(dumpliErr)
				if CFG.Env.DropOnDuplicate && th != torr.Hash {
					log.Warn(log.Struct{"sid", s.SubjId, "torrfn", torr.Name, "torrHash", torr.Hash, "size", torr.Size}, "delete dumplicate file")
					merr.Add(DL.Qbt.DelTorrentsFs(torr.Hash))
				}
				return merr.Err()
			}
			err = DL.Qbt.RenameFile(torr.Hash, torr.Name, rename)
			if err != nil {
				return err
			}
			s.Pushed[se] = torr.Hash
			if CFG.Env.BgmiLog {
				CFG.BgmiLogger.Infof(eslog.Struct{"sid", s.SubjId, "name", s.Name}, "episode update: %s", rename)
			}
		} else {
			log.Info(log.Struct{"sid", s.SubjId, "torrfn", torr.Name, "torrHash", torr.Hash}, "not a videoFile,may external subtitles")
			ok, rn, err := renameSubRssTorr(s, torr)
			log.Error(log.Struct{"err", err}, "rename RssTorr with subtitles failed")
			if !ok {
				return nil
			}
			se = util.TrimExtensionAndGetEpi(rn)
		}
		mErr := errs.MultiErr{}
		err := pusher.Push(P.Payload{
			SubjectId:    s.SubjId,
			SubjectName:  s.Name,
			DownLoadName: torr.Name,
			Size:         torr.Size,
			Episode:      se,
		})
		mErr.Add(err)
		mErr.Add(s.writeJson())
		// FIXME if rss contains sp ,may be exit early
		episNum := s.Episode
		if episNum != 0 && len(s.Pushed) >= episNum {
			log.Info(log.Struct{"sid", s.SubjId, "resType", "RSS"}, "compiled,exited now")
			s.terminate()
		}
		return mErr.Err()
	} else { //Movie
		if _, e := s.Pushed[torr.Hash]; e {
			return fmt.Errorf("%w: name=%s", errs.ErrItemAlreadyPushed, torr.Name)
		}
		mErr := errs.MultiErr{}
		s.Pushed[torr.Hash] = ""
		err := pusher.Push(P.Payload{
			SubjectId:    s.SubjId,
			SubjectName:  s.Name,
			DownLoadName: torr.Name,
			Size:         torr.Size,
			Episode:      "Movie",
		})
		mErr.Add(err)
		mErr.Add(s.writeJson())
		return mErr.Err()
	}
}

func (s *Subject) terminate() {
	log.Debug(log.Struct{"sid", s.SubjId, "resType", s.ResourceTyp.String()}, "exited")
	s.Terminate, s.Finished = true, true
	//FIXME: maybe ignore writejson case writing while exit context
	err := s.writeJson()
	if err != nil {
		log.Info(log.Struct{"err", err}, "write json failed")
	}
	s.Exit()
}

// readRssAndDownload reads RSS feeds and downloads the corresponding files.
// It performs the following steps:
// 1. Checks if the resource type is Torrent, and if so, returns immediately.
// 2. Reads the RSS feed using the RssReader.
// 3. If reading is successful, iterates over the read items.
// 4. For each item, constructs a fake filename and attempts to rename it using renameTV.
// 5. If renaming fails, logs the error and uses the original title as the renamed value.
// 6. Checks for duplicate episodes and skips them if found.
// 7. Downloads the torrent file using the DefaultDownLoader.
// 8. If the download fails, logs the error, undoes the RSS read operation, and removes the renamed entry.
// 9. If the download succeeds, calls builtinDownload with the downloaded torrent information.
// 10. Updates the TorrentUrls map with the new torrent information.
// 11. Updates the RssGuids with the current GUIDs from the RssReader.
func (s *Subject) readRssAndDownload() {
	if s.ResourceTyp == Torrent {
		return
	}
	read, ok, err := s.RssReader.Read()
	if err != nil {
		log.Error(log.Struct{"err", err}, "read rss error")
		return
	}
	if ok {
		for _, r := range read {
			fakeFn := r.Title + ".mp4"
			renamed, err := renameTV(s, fakeFn)
			if err != nil {
				log.Error(log.Struct{"err", err}, "rename failed")
				renamed = r.Title
			} else {
				renamed = strings.TrimSuffix(renamed, ".mp4")
			}
			if s.RssTorrentsName == nil {
				s.RssTorrentsName = map[string]struct{}{}
			}
			if _, ex := s.RssTorrentsName[renamed]; ex {
				log.Warn(log.Struct{"file", r.Desc}, "skip duplicate episode")
				continue
			}
			s.RssTorrentsName[renamed] = struct{}{}
			rop := RssFileOpt{Renamed: renamed}
			fop := FilePath{FileName: &rop, DirPath: s.Path}
			t, err := builtin.DefaultDownLoader.Download(r.TorrUrl, fop, nil)
			if err != nil {
				log.Error(log.Struct{"err", err}, "download failed")
				s.RssReader.Undo(r.Guid)
				delete(s.RssTorrentsName, renamed)
				continue
			}
			s.builtinDownload(builtin.MonitoredTorrent{Url: r.TorrUrl, TorrentInfo: builtin.TorrentInfo{Rename: renamed, Torrent: t}})
			s.TorrentUrls[r.TorrUrl] = RssFileOptStrage{renamed}
		}
		s.RssGuids = s.RssReader.Guids()
	}

}

func (s *Subject) pushBuiltin(torr builtin.MonitoredTorrent, pusher P.Pusher) error {
	log.Debug(log.Struct{"sid", s.SubjId, "torrName", torr.Rename}, "push builtin")
	pusher.Push(P.Payload{})
	return nil
}
