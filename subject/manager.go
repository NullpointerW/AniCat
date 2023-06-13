package subject

import (
	// "log"
	"sync"
)

type SubjC struct {
	N string
	Extra
}

type Pip struct {
	arg any
	err error
	wg  sync.WaitGroup
}

func NewPip(a any) *Pip {
	p := new(Pip)
	p.arg = a
	p.wg.Add(1)
	return p
}

func (p *Pip) Error() error {
	p.wg.Wait()
	return p.err
}

var (
	Create, Delete chan *Pip
)

func init() {
	Create = make(chan *Pip, 1024)
	Delete = make(chan *Pip, 1024)
}

var Manager = SubjectManager{
	mu:  new(sync.Mutex),
	sto: make(map[int]*Subject),
}

type SubjectManager struct {
	mu  *sync.Mutex
	sto map[int]*Subject
	// a snapshot copy from the last list()-calling make caller fast get list
	copy []Subject
}

func (m *SubjectManager) Add(s *Subject) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sto[s.SubjId] = s
	m.copy = nil
}

func (m *SubjectManager) Remove(s *Subject) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sto, s.SubjId)
	m.copy = nil
}

func (m SubjectManager) Get(sid int) *Subject {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sto[sid]
}

func (m *SubjectManager) List() (ls []Subject) {
	m.mu.Lock()
	if m.copy != nil {
		return m.copy
	}
	defer m.mu.Unlock()
	for _, ss := range m.sto {
		ls = append(ls, *ss)
	}
	m.copy = ls
	return ls
}

func StartManagement() {
	for {
		select {
		case p := <-Create:
			sc := p.arg.(SubjC)
			err := CreateSubject(sc.N, &sc.Extra)
			if err != nil {
				p.err = err
			}
			p.wg.Done()
		case p := <-Delete:
			i := p.arg.(int)
			s := Manager.Get(i)
			if s != nil {
				if !s.Terminate {
					s.Exit()
				}
				Manager.Remove(s)
				s.RmRes()
				err := rmFolder(s)
				if err != nil {
					p.err = err
				}

			}
			p.wg.Done()
		}
	}

}
