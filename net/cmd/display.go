package cmd

import (
	
)

var (
	GreenBg  = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	RedBg    = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	Cyan     = string([]byte{27, 91, 51, 54, 109})
	YellowBg = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	Cls      = "\033[2J\033[H"
	Reset    = string([]byte{27, 91, 48, 109})
	// $
	Cursor = "\033[?25h$"
)



const (
	usageHelp = "\n   Usage:\n         " +
		"mikan  <command> [argument(s)]\n   " +
		"The commands are:\n\n         " +
		"add [name] [-g -i -mc ...]   add a anine-subject\n         " +
		"rm [subjid]                  delete a anine-subject\n         " +
		"ls                           show all anine-subjects\n         " +
		"lsi [name]                   show anine resource list\n         " +
		"lsg [name]                   show anine subtitleGroup list (rss type)\n"
	addCMDUsageHelp = "\n   Usage:\n         " +
		"mikan add [name] [arguments]\n   " +
		"The arguments are:\n\n         " +
		"-mn                          the substring that the torrent name must not contain (rss auto download rule)\n         " +
		"-mc                          the substring that the torrent name must contain (rss auto download rule)\n         " +
		"-rg                          enable regex mode in \"-mc\" and \"-mn\"\n         " +
		"-g,--group                   specified  subtitleGroup (rss type)\n         " +
		"-i,--index                   specified  index from torrents list (torr type)\n"
)

// just for test
func TestingString() (text string) {
	text = "\n   Usage:\n         " +
		"mikan add [name] [arguments]\n   " +
		"The arguments are:\n\n         " +
		"-mn                          the substring that the torrent name must not contain (rss auto download rule)\n         " +
		"-mc                          the substring that the torrent name must contain (rss auto download rule)\n         " +
		"-rg                          enable regex mode in \"-mc\" and \"-mn\"\n         " +
		"-g,--group                   specified  subtitleGroup (rss type)\n         " +
		"-i,--index                   specified  index from torrents list (torr type)\n"
	return
}
