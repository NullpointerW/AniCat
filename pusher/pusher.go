package pusher

type Payload struct {
	SubjectId    int    // $id$
	SubjectName  string // $name$
	DownLoadName string // $dlname$
	Size         int    // $size$
}

type Pusher interface {
	Push(p Payload) error
	Template(t string)
}

type Mock struct{}

func (m Mock) Push(p Payload) error {
	return nil
}
func (m Mock) Template(t string) {

}

var Mock_ Pusher = Mock{}
