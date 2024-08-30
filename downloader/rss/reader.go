package rss

import (
	"github.com/mmcdole/gofeed"
	"maps"
)

type FilterFunc func(n string) bool
type Reader struct {
	parser gofeed.Parser
	feed   string
	guids  map[string]struct{}
	filter FilterFunc
}

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
	return &Reader{
		feed:   feed,
		guids:  guids,
		filter: filterFunc,
	}
}

func (r *Reader) Read() ([]Item, bool, error) {
	f, err := r.parser.ParseURL(r.feed)
	if err != nil {
		return nil, false, err
	}
	var read []Item
	for _, it := range f.Items {
		if _, ex := r.guids[it.GUID]; !ex {
			itt := Item{}
			itt.Guid = it.GUID
			itt.TorrUrl = it.Enclosures[0].URL
			itt.Title = it.Title
			itt.Desc = it.Description
			if r.filter != nil && r.filter(itt.Title) {
				read = append(read, itt)
			}
			r.guids[it.GUID] = struct{}{}
		}
	}
	return read, len(read) > 0, nil
}
func (r *Reader) ReadOne() (Item, bool, error) {
	f, err := r.parser.ParseURL(r.feed)
	if err != nil {
		return Item{}, false, err
	}
	for _, it := range f.Items {
		if _, ex := r.guids[it.GUID]; !ex {
			itt := Item{}
			itt.TorrUrl = it.Enclosures[0].URL
			itt.Title = it.Title
			itt.Desc = it.Description
			r.guids[it.GUID] = struct{}{}
			if r.filter != nil && r.filter(itt.Title) {
				return itt, true, nil
			}
		}
	}
	return Item{}, false, nil
}

func (r *Reader) Seek() ([]Item, bool, error) {
	f, err := r.parser.ParseURL(r.feed)
	if err != nil {
		return nil, false, err
	}
	var read []Item
	for _, it := range f.Items {
		if _, ex := r.guids[it.GUID]; !ex {
			itt := Item{}
			itt.TorrUrl = it.Enclosures[0].URL
			itt.Title = it.Title
			itt.Desc = it.Description
			if r.filter != nil && r.filter(itt.Title) {
				read = append(read, itt)
			}
		}
	}
	return read, len(read) > 0, nil
}

func (r *Reader) Guids() map[string]struct{} {
	var readonly = make(map[string]struct{})
	maps.Copy(readonly, r.guids)
	return readonly
}

func (r *Reader) Undo(guid string) {
	delete(r.guids, guid)
}
