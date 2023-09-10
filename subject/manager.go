package subject

import (
	"github.com/NullpointerW/anicat/log"
	"sync"

	"github.com/NullpointerW/anicat/errs"
)

type createType int

const (
	CreateViaStr createType = iota
	CreateViaFeed
)

type SubjC struct {
	N string
	Extra
	CreateTyp createType
}

type Pip struct {
	Arg any
	err error
	wg  sync.WaitGroup
}

func NewPip(a any) *Pip {
	p := new(Pip)
	p.Arg = a
	p.wg.Add(1)
	return p
}

func (p *Pip) Error() error {
	p.wg.Wait()
	return p.err
}

var (
	Delete chan *Pip
	Create chan *Pip
)

func init() {
	Create = make(chan *Pip, 1024)
	Delete = make(chan *Pip, 1024)
}

var Manager = SubjectManager{
	mu:  new(sync.Mutex),
	sto: make(map[int]*Subject),
	// sp_sid: make(map[string]int),
}

type SubjectManager struct {
	mu  *sync.Mutex
	sto map[int]*Subject
	// sp_sid map[string]int
	// a snapshot copy from the last list()-calling make caller fast get list
	copy []Subject
}

// func (m *SubjectManager) GetSidViaSp(sp string) int {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	sid, e := m.sp_sid[sp]
// 	if !e {
// 		return -1
// 	}
// 	return sid
// }

func (m *SubjectManager) Sync() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.copy = nil
}
func (m *SubjectManager) Add(s *Subject) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sto[s.SubjId] = s
	// m.sp_sid[s.Path] = s.SubjId
	m.copy = nil
}

func (m *SubjectManager) Remove(sid int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sto, sid)
	// delete(m.sp_sid, s.Path)
	m.copy = nil
}

func (m *SubjectManager) Get(sid int) *Subject {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sto[sid]
}

func (m *SubjectManager) List() (ls []Subject) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.copy != nil {
		return m.copy
	}
	for _, ss := range m.sto {
		ls = append(ls, *ss)
	}
	m.copy = ls
	return ls
}

func (m *SubjectManager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// alloc a new map
	m.sto = make(map[int]*Subject)
}

func StartManagement() {
	for {
		select {
		case p := <-Create:
			sc := p.Arg.(SubjC)
			var (
				sid int
				err error
			)
			if sc.CreateTyp == CreateViaStr {
				sid, err = CreateSubject(sc.N, &sc.Extra)
			} else { // CreateViaFeed
				sid, err = CreateSubjectViaFeed(sc.N, sc.Extra.RssOption.Name, &sc.Extra)
			}
			if err != nil {
				p.err = err
			} else {
				// send added subjId to peer
				p.Arg = sid
			}
			p.wg.Done()
		case p := <-Delete:
			i, _ := p.Arg.(int)
			errWrap := errs.ErrWrapper{}
			// rm *
			// remove all subjects
			if i == 0 && p.Arg.(string) == "*" {
				log.Warn(nil, "rm: remove all subjects")
				merr := errs.MultiErr{}
				for _, s := range Manager.List() {
					if !s.Terminate {
						s.Exit()
					}
					errWrap.Handle(func() error {
						return s.RmRes()
					})
					errWrap.Handle(func() error {
						return RmFolder(&s)
					})
					if errWrap.Error() == nil {
						Manager.Remove(s.SubjId)
					}
					merr.Add(errWrap.Error())
					errWrap.Reset()
				}
				p.err = merr.Err()
				p.wg.Done()
				continue
			}
			s := Manager.Get(i)
			if s != nil {
				if !s.Terminate {
					s.Exit()
				}
				errWrap.Handle(func() error {
					return s.RmRes()
				})
				errWrap.Handle(func() error {
					return RmFolder(s)
				})
				if errWrap.Error() == nil {
					Manager.Remove(s.SubjId)
				}
				p.err = errWrap.Error()
			}
			p.wg.Done()
		}
	}

}
