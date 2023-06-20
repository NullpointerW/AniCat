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
	reg0  = `[sS](\d+)`       // s2
	reg1  = `(?i)season(\d+)` // season2
	zhreg = `第(.+?)季`         // 第二季
)

var sregs = []string{reg0, reg1}

const (
	reg0v_epi = `]\[(\d{2})[vV]` // [02v1]
	reg1_epi  = `\[(\d+)\]`      // [02]
	reg2_epi  = `- ?(\d+)`    // - 02
)

var epi_regs = []string{reg0v_epi, reg1_epi, reg2_epi}
