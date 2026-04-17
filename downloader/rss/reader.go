package rss

import (
	"maps"

	"github.com/mmcdole/gofeed"
)

type FilterFunc func(n string) bool

// Reader is a stateful RSS feed reader that tracks processed GUIDs and
// applies an optional filter function to item titles.
type Reader struct {
	parser gofeed.Parser
	feed   string
	guids  map[string]struct{}
	filter FilterFunc
}

// Item is a single RSS entry returned by the Reader.
type Item struct {
	TorrUrl string
	Desc    string
	Title   string
	Guid    string
}

func NewReader(feed string, guids map[string]struct{}, filterFunc FilterFunc) *Reader {
	if guids == nil {
		guids = make(map[string]struct{})
	}
	return &Reader{feed: feed, guids: guids, filter: filterFunc}
}

// passes returns true if item passes the filter (or no filter is set).
func (r *Reader) passes(title string) bool {
	return r.filter == nil || r.filter(title)
}

// toItem converts a gofeed.Item to our Item type.
// Assumes at least one enclosure is present.
func toItem(it *gofeed.Item) Item {
	return Item{
		Guid:    it.GUID,
		TorrUrl: it.Enclosures[0].URL,
		Title:   it.Title,
		Desc:    it.Description,
	}
}

// Read returns all new, filter-passing items since the last call.
// GUIDs are only marked as seen for items that pass the filter, so
// filtered-out items will be reconsidered on the next call if the
// filter changes.
func (r *Reader) Read() ([]Item, bool, error) {
	f, err := r.parser.ParseURL(r.feed)
	if err != nil {
		return nil, false, err
	}
	var read []Item
	for _, it := range f.Items {
		if _, seen := r.guids[it.GUID]; seen {
			continue
		}
		itt := toItem(it)
		if !r.passes(itt.Title) {
			continue
		}
		r.guids[it.GUID] = struct{}{}
		read = append(read, itt)
	}
	return read, len(read) > 0, nil
}

// ReadOne returns the first new, filter-passing item.
func (r *Reader) ReadOne() (Item, bool, error) {
	f, err := r.parser.ParseURL(r.feed)
	if err != nil {
		return Item{}, false, err
	}
	for _, it := range f.Items {
		if _, seen := r.guids[it.GUID]; seen {
			continue
		}
		itt := toItem(it)
		r.guids[it.GUID] = struct{}{}
		if r.passes(itt.Title) {
			return itt, true, nil
		}
	}
	return Item{}, false, nil
}

// Seek returns all filter-passing items without marking them as seen.
// Used for one-time lookback (e.g. finding a collection torrent on first run).
func (r *Reader) Seek() ([]Item, bool, error) {
	f, err := r.parser.ParseURL(r.feed)
	if err != nil {
		return nil, false, err
	}
	var read []Item
	for _, it := range f.Items {
		if _, seen := r.guids[it.GUID]; seen {
			continue
		}
		itt := toItem(it)
		if r.passes(itt.Title) {
			read = append(read, itt)
		}
	}
	return read, len(read) > 0, nil
}

// Guids returns a snapshot copy of the seen-GUID set.
func (r *Reader) Guids() map[string]struct{} {
	out := make(map[string]struct{}, len(r.guids))
	maps.Copy(out, r.guids)
	return out
}

// Undo removes a GUID from the seen set so it will be reconsidered on the next Read.
func (r *Reader) Undo(guid string) {
	delete(r.guids, guid)
}
