package subject

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/NullpointerW/anicat/log"
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

// Pip is a typed request/response envelope for operations sent over Create or Delete channels.
// The caller creates a Pip, sends it on the channel, then calls Wait() to block until done.
type Pip struct {
	Arg any // input argument (SubjC for Create, int or "*" for Delete)
	Sid int // output: the assigned subject ID (valid only after successful Create)
	err error
	wg  sync.WaitGroup
}

func NewPip(a any) *Pip {
	p := new(Pip)
	p.Arg = a
	p.wg.Add(1)
	return p
}

// Error blocks until the operation completes and returns any error.
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
	copy atomic.Value
	wg   sync.WaitGroup
}

func (m *Manager) Sync() {
	m.copy.Store(([]Subject)(nil))
}
func (m *Manager) Add(s *Subject) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sto[s.SubjId] = s
	m.Sync()
}

func (m *Manager) Remove(sid int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sto, sid)
	m.Sync()
}

func (m *Manager) Get(sid int) *Subject {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sto[sid]
}

func (m *Manager) List() (ls []Subject) {
	if v := m.copy.Load(); v != nil && v.([]Subject) != nil {
		return v.([]Subject)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// double check
	if v := m.copy.Load(); v != nil && v.([]Subject) != nil  {
		return v.([]Subject)
	}
	for _, ss := range m.sto {
		ls = append(ls, *ss)
	}
	m.copy.Store(ls)
	return ls
}

func (m *Manager) clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// alloc a new map
	m.sto = make(map[int]*Subject)
}

func (m *Manager) Exit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range m.sto {
		if !s.Terminate {
			s.Exit()
		}
	}
	m.wg.Wait()
}

func (m *Manager) Range(f func(int, *Subject) bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, s := range m.sto {
		if !f(i, s) {
			return
		}
	}
}

func StartManagement() {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error(log.Struct{"panic", r}, "StartManagement: recovered from panic, loop continues")
				}
			}()
			select {
			case p := <-Create:
				sc := p.Arg.(SubjC)
				var (
					sid int
					err error
				)
				if sc.CreateTyp == CreateViaStr {
					sid, err = CreateSubject(sc.N, &sc.Extra)
				} else {
					sid, err = CreateSubjectViaFeed(sc.N, sc.Extra.RssOption.Name, &sc.Extra)
				}
				if err != nil {
					p.err = err
				} else {
					p.Sid = sid
				}
				p.wg.Done()
			case p := <-Delete:
				errWrap := errs.ErrWrapper{}
				switch v := p.Arg.(type) {
				case string:
					if v != "*" {
						p.err = fmt.Errorf("unexpected delete arg: %q", v)
						p.wg.Done()
						return
					}
					log.Warn(nil, "rm: remove all subjects")
					merr := errs.MultiErr{}
					for _, s := range Mgr.List() {
						if !s.Terminate {
							s.Exit()
						}
						errWrap.Handle(func() error { return s.RmRes() })
						errWrap.Handle(func() error { return RmFolder(&s) })
						if errWrap.Error() == nil {
							Mgr.Remove(s.SubjId)
						}
						merr.Add(errWrap.Error())
						errWrap.Reset()
					}
					p.err = merr.Err()
				case int:
					s := Mgr.Get(v)
					if s != nil {
						if !s.Terminate {
							s.Exit()
						}
						errWrap.Handle(func() error { return s.RmRes() })
						errWrap.Handle(func() error { return RmFolder(s) })
						if errWrap.Error() == nil {
							Mgr.Remove(s.SubjId)
						}
						p.err = errWrap.Error()
					}
				default:
					p.err = fmt.Errorf("delete: unexpected arg type %T", p.Arg)
				}
				p.wg.Done()
			}
		}()
	}
}
