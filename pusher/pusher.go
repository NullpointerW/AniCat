package pusher

type Payload struct {
	SubjectId    int    // $id$
	SubjectName  string // $name$
	DownLoadName string // $dlname$
	Size         int    // $size$
}

type Pusher interface {
	Push(p Payload)
	Template(t string)
}

var Mock Pusher
