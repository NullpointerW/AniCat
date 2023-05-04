package subject

// resource type
type ResourceTyp int

func (t ResourceTyp) String() string {
	if t == 0 {
		return "RSS"
	}
	return "Torrent"
}

const (
	RSS ResourceTyp = iota
	Torrent
)

type BgmiTyp int

func (t BgmiTyp) String() string {
	if t == 0 {
		return "TV"
	}
	return "MOVIE"
}

// anime type
const (
	TV BgmiTyp = iota
	MOVIE
)

// file-related
const (
	FolderSuffix = "@mikan"
	jsonfileName = "info.json"
)

// qbt tag generation template
const (
	QbtTag_prefix = "mikan@subject-"
	QbtTag        = QbtTag_prefix + "%d"
)

const (
	CoverFN = "cover.jpg"
)
