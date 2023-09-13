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
	jsonfileName = "meta-data#%s.json"
)

// qbt tag generation template
const (
	QbtTag_prefix = "anicat@subj-"
	QbtTag        = QbtTag_prefix + "%d" // identity for each resource
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
	reg0v_epi  = `\[(\d{2})[vV]`  // [02v1] // \[(\d+)\w*\]
	reg1_epi   = `\[(\d+)\]`      // [02]
	reg2_epi   = `- ?(\d+)`       // - 02
	reg3_epi   = `\[(\d+)[集话話]\]` // [02集]
	reg4_epi   = `第(\d+)[话話集]`    // 第02集
	regSpecial = `#\s*(\d+)`      // #02
)

const (
	reg0_coll = `\[[^\]]*\d+(\.\d+)\s*[Gg][Bb][^\]]*\]`
	_         = `\d+-\d+`
	reg2_coll = `\d+-(\d+)`
)

const (
	reg0_part = `第(\d)クール`
	reg1_part = `第(\d)部分`
	reg2_part = `(?i)Part\s*(\d)`

	reg_part2 = `後半|后半`
)

var coll_regs = []string{reg0_coll}

var epi_regs = []string{reg0v_epi, reg1_epi, reg2_epi, reg3_epi, reg4_epi}

var part_regs = []string{reg0_part, reg1_part, reg2_part}

const (
	CHSSubStationReg = "(?i)简体|簡體|简中|簡中|chs"
	CHTSubStationReg = "(?i)繁体|繁體|正體中文|正体中文|繁中|cht"
)
