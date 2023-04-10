package subject

import "sync"

const (
	RSS = iota
	Torrent
	TV
	MOVIE
)

type SubjectManager struct {
	mu                   *sync.Mutex
	finished, unfinished map[int]*Subject
}

func (m SubjectManager) Add(sid int, s *Subject, fin bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if fin {
		m.finished[sid] = s
		return
	}
	m.unfinished[sid] = s
}

func (m SubjectManager) Remove(sid int, s *Subject, fin bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if fin {
		delete(m.finished, sid)
		return
	}
	delete(m.unfinished, sid)

}

func (m SubjectManager) Move(sid int, s *Subject, tofin bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if tofin {
		delete(m.unfinished, sid)
		m.finished[sid] = s
		return
	}
	delete(m.finished, sid)
	m.unfinished[sid] = s
}

type Subject struct {
	SubjId      int    `json:"subjId"`
	Name        string `json:"name"`
	Path        string
	Finished    bool
	Episode     int
	ResourceTyp int
	ResourceUrl string
	Typ         int
	StartTime   string
	EndTime     string
}
