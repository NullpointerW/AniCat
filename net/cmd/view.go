package cmd

var (
	GreenBg = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	RedBg   = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	Cyan    = string([]byte{27, 91, 51, 54, 109})
	Cls     = "\033[2J\033[H"
	Reset   = string([]byte{27, 91, 48, 109})
	// $ 
	Cursor  = "\033[?25h$"
)

const (
	usageHelp = `

	Usage:
	             mikan  <command> [anine-name]
	The commands are:

	             add    add a anine-subject
	             rm     delete a anine-subject
	             ls     show all anine-subjects   
	             `
)
