package pusher

type Payload struct {
	SubjectId    int    // ${id}
	SubjectName  string // ${name}
	DownLoadName string // ${dlname}
	Episode      string // ${epi}
	Size         int    // ${size}
}

type Pusher interface {
	Push(p Payload) error
}

type Mock struct{}

func (m Mock) Push(p Payload) error {
	return nil
}

var Mock_ Pusher = Mock{}
