package rename

const (
	chsSubStationReg = "(?i)简体|簡體|简中|簡中|chs"
	chtSubStationReg = "(?i)繁体|繁體|正體中文|正体中文|繁中|cht"
)

const (
	epiReg0v   = `\[(\d{2})[vV]`        // [02v1] // \[(\d+)\w*\]
	epiReg1    = `\[(\d+)\]`            // [02]
	epiReg2    = `- ?(\d+)`             // - 02
	epiReg2spec    = `\[(\d+)(?:\s+END)?\]` // [02 END]
	epiReg3    = `\[(\d+)[集话話]\]`       // [02集]
	epiReg4    = `第(\d+)[话話集]`          // 第02集
	epiReg5    = `\s+(\d+)`             // xxx 02
	specialReg = `#\s*(\d+)`            // #02
)

var epiRegs = []string{epiReg0v, epiReg1, epiReg2, epiReg2spec,epiReg3, epiReg4, epiReg5}
