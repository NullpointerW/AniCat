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
	FolderSuffix = "@anicat"
	jsonfileName = "meta-data.json"
)

// qbt tag generation template
const (
	QbtTag_prefix = "anicat@subj-"
	QbtTag        = QbtTag_prefix + "%d"
)

const (
	CoverFN = "folder.jpg" // adapt infuse
)

const (
	reg0 = `[sS](\d+)`
	reg1 = `(?i)season(\d+)`
	zhreg = `第(.+?)季`
)

var sregs = []string{reg0, reg1}
