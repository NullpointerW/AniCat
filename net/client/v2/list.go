package main

import (
	"encoding/json"
	"fmt"
	// "fmt"
	"os"
	"strconv"
	"strings"

	N "github.com/NullpointerW/anicat/net"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m *model) torrlsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				if m.list.FilterState() == list.Unfiltered {
					m.istrls = false
					m.mod = text
					m.textInput.Focus()
					return m, nil
				}
			case "enter":
				m.istrls = false
				// panic("DEBUG_PANIC: " + m.history[len(m.history)-1])
				idx := strings.Split(m.list.SelectedItem().(item).title, ".")[0]
				cmd := "add " + getArg(m.history[len(m.history)-1]) + " -i " + idx
				m.sendCmd(cmd)
				return m, m.spinner.Tick
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) rsslsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				if m.list.FilterState() == list.Unfiltered {
					m.istrls = false
					m.mod = text
					m.textInput.Focus()
					return m, nil
				}
			case "enter":
				m.ischdls = true
				m.loadChidlist(m.list.SelectedItem().FilterValue())
				// blur
				// m.chdlist.Select(-1)
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) rssChildlsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				if m.chdlist.FilterState() == list.Unfiltered {
					m.ischdls = false
					return m, nil
				}
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.chdlist, cmd = m.chdlist.Update(msg)
	return m, cmd
}

func (m *model) listView() string {
	if m.ischdls {
		return docStyle.Render(m.chdlist.View())
	} else {
		return docStyle.Render(m.list.View())
	}
}

type lsirtyp int

const (
	torr = lsirtyp(iota)
	rss
	er
)

func lsiReturnTyp(raw string) lsirtyp {
	var (
		torrls []N.TorrItem
	)
	err := json.Unmarshal([]byte(raw), &RssGroup)
	if len(RssGroup) != 0 && err == nil {
		return rss
	}

	err = json.Unmarshal([]byte(raw), &torrls)
	if len(torrls) != 0 && err == nil && torrls[0].Name != "" {
		return torr
	}

	return er
}

func (m *model) NewTorrlist(raw string) {
	var (
		torrls []N.TorrItem
		items  []list.Item
	)
	json.Unmarshal([]byte(raw), &torrls)
	for i, t := range torrls {
		it := item{}
		it.title = strconv.Itoa(i+1) + "." + t.Name
		it.desc = t.UpdateTime + " | " + t.Size
		items = append(items, it)
	}
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	wd, hd := docStyle.GetFrameSize()
	// panic(fmt.Sprintln("width", w-wd, "heigh", h-hd))
	m.list = list.New(items, list.NewDefaultDelegate(), w-wd, h-hd)
	m.list.Title = "torrents"
}

var RssGroup map[string][]N.TorrItem

func (m *model) NewRsslist(raw string) {
	var items []list.Item
	// json.Unmarshal([]byte(raw), &RssGroup)
	for group, its := range RssGroup {
		it := item{}
		total := strconv.Itoa(len(its))
		up := its[0].UpdateTime
		it.title = group
		it.desc = fmt.Sprintf("total: %s | lastUpdate: %s", total, up)
		items = append(items, it)
	}
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	wd, hd := docStyle.GetFrameSize()
	// panic(fmt.Sprintln("width", w-wd, "heigh", h-hd))
	m.list = list.New(items, list.NewDefaultDelegate(), w-wd, h-hd)
	m.list.Title = "subtitle groups"
}

func (m *model) loadChidlist(group string) {
	var (
		chdItems []list.Item
		its      []N.TorrItem = RssGroup[group]
	)
	for _, itd := range its {
		it := item{}
		it.title = itd.Name
		it.desc = fmt.Sprintf("updateTime: %s | size: %s", itd.UpdateTime, itd.Size)
		chdItems = append(chdItems, it)
	}
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	wd, hd := docStyle.GetFrameSize()
	m.chdlist = list.New(chdItems, list.NewDefaultDelegate(), w-wd, h-hd)
	m.chdlist.Title = "items"
}
