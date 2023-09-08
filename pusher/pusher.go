package pusher

// Payload can customize template such as the following format:
//
//	 <body>
//		 <h1>bgmi_name: ${name}</h1>
//		 <h2>epi: ${epi}</h2>
//		 <ul>
//			 <li>bmgiTVid: ${id}</li>
//			 <li>file_name: ${dlname}</li>
//			 <li>file_size: ${size}</li>
//		 </ul>
//		 <img src="http://api.bgm.tv/v0/subjects/${id}/image?type=medium">
//		 <p>dl_compl</p>
//		 <p>Enjoy it,</p>
//		 <p>AniCat</p>
//		 </div>
//	 </body>
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
