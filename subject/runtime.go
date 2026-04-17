package subject

import (
	"context"

	"github.com/NullpointerW/anicat/downloader/builtin"
	"github.com/NullpointerW/anicat/downloader/rss"
	util "github.com/NullpointerW/anicat/utils"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

// RuntimeState holds all goroutine-lifecycle and channel fields for a running Subject.
// These fields are never persisted to JSON.
//
// Separating runtime state from the persisted Subject record makes the boundary between
// "what lives on disk" and "what lives in memory" explicit, and prevents accidental
// serialisation of channels or function values.
type RuntimeState struct {
	// Exit cancels the context that drives this subject's goroutine.
	// Call it to request a graceful shutdown.
	Exit context.CancelFunc

	// Exited is closed when the subject's goroutine has fully stopped.
	// Readers (e.g. Jellyfin helper, detector) can select on this to know
	// the goroutine is gone before accessing shared state.
	Exited chan struct{}

	// OperationChan carries in-band control messages (e.g. Rename) to the
	// run goroutine while it is alive.
	OperationChan chan Operate

	// --- qBittorrent downloader ---

	// PushChan receives completed-download events from the detector goroutine.
	PushChan chan qbt.Torrent

	// --- builtin downloader ---

	// MonitorChan receives MonitoredTorrent objects that need to be tracked.
	// MonitorBuiltin reads from this channel.
	MonitorChan chan builtin.MonitoredTorrent

	// PushChanBuiltin receives fully-completed MonitoredTorrent objects from
	// MonitorBuiltin once a download is done.
	PushChanBuiltin chan builtin.MonitoredTorrent

	// RssReader is the stateful RSS feed reader used by the builtin downloader.
	// It tracks which GUIDs have already been processed.
	RssReader *rss.Reader

	// FinishedTorrentNameList is the thread-safe progress list shown to the CLI.
	FinishedTorrentNameList *util.ListView[builtin.TorrentProgress]

	// TorrentMonitor provides TTL-cached download progress for the builtin downloader.
	TorrentMonitor *builtin.TorrentProgressMonitor
}
