package cmd

type Command struct {
	N    string
	Opt  Option
	Flag struct {
		SubtitleGroup  string
		MustContain    string
		MustNotContain string
		useRegex       bool
	}
}
type Option int

const (
	Add Option = iota
	Del
	Ls
)

func Parse(cmd string) (reply Command) {
	return Command{}
}
