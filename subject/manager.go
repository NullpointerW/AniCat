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

var Mgr = Manager{
	sto: make(map[int]*Subject),
}

type Manager struct {
	mu  sync.Mutex
	sto map[int]*Subject
	// a snapshot copy from the last list()-calling make caller fast get list
	copy []Subject
	wg   sync.WaitGroup
}

func (m *Manager) Sync() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.copy = nil
}
func (m *Manager) Add(s *Subject) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sto[s.SubjId] = s
	m.copy = nil
}

func (m *Manager) Remove(sid int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sto, sid)
	m.copy = nil
}

func (m *Manager) Get(sid int) *Subject {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sto[sid]
}

func (m *Manager) List() (ls []Subject) {
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

func (m *Manager) clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// alloc a new map
	m.sto = make(map[int]*Subject)
}

func (m *Manager) Exit() {
	m.wg.Wait()
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
				for _, s := range Mgr.List() {
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
						Mgr.Remove(s.SubjId)
					}
					merr.Add(errWrap.Error())
					errWrap.Reset()
				}
				p.err = merr.Err()
				p.wg.Done()
				continue
			}
			s := Mgr.Get(i)
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
					Mgr.Remove(s.SubjId)
				}
				p.err = errWrap.Error()
			}
			p.wg.Done()
		}
	}

}
