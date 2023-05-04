package subject

import (
	"log"
	"sync"
)

var (
	Create chan string
	Delete chan int
)

func init() {
	Create = make(chan string, 1024)
	Delete = make(chan int, 1024)
}

var Manager = SubjectManager{
	mu:  new(sync.Mutex),
	sto: make(map[int]*Subject),
}

type SubjectManager struct {
	mu  *sync.Mutex
	sto map[int]*Subject
}

func (m SubjectManager) Add(s *Subject) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sto[s.SubjId] = s
}

func (m SubjectManager) Remove(s *Subject) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sto, s.SubjId)
}

func (m SubjectManager) GetSubject(sid int) *Subject {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sto[sid]
}

func StartManagement() {
	for {
		select {
		case n := <-Create:
			err := CreateSubject(n)
			if err != nil {
				log.Println(err)
			}
		case i := <-Delete:
			s := Manager.GetSubject(i)
			if s != nil {
				s.exit()
				Manager.Remove(s)
				err := rmFolder(s)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

}
