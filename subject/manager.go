package subject

import "sync"

var Manager = SubjectManager{
	mu:         new(sync.Mutex),
	finished:   make(map[int]*Subject),
	unfinished: make(map[int]*Subject),
}

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

func (m SubjectManager) Move(sid int, tofin bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var copy *Subject
	if tofin {
		copy = m.unfinished[sid]
		delete(m.unfinished, sid)
		m.finished[sid] = copy
		return
	}
	copy = m.finished[sid]
	delete(m.finished, sid)
	m.unfinished[sid] = copy
}

func (m SubjectManager) GetSubject(sid int) *Subject {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, e := m.unfinished[sid]
	if !e {
		return m.finished[sid]
	} else {
		return s
	}
}
