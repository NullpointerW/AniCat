package main

import (
	"encoding/json"
	"strconv"

	N "github.com/NullpointerW/anicat/net"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		switch msg.String() {
		case "ctrl+c":
			{
				return m, tea.Quit
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

func (m model) View2() string {
	return docStyle.Render(m.list.View())
}

type lsirtyp int

const (
	torr = lsirtyp(iota)
	rss
	t
)

func lsiReturnTyp(raw string) lsirtyp {
	var (
		rssGroups []N.RssGroup
		torrls    []N.TorrItem
	)
	json.Unmarshal([]byte(raw), &rssGroups)
	if len(rssGroups) != 0 {
		return rss
	}
	json.Unmarshal([]byte(raw), &torrls)
	if len(torrls) != 0 {
		return torr
	}
	return t
}

func (m *model) NewTorrlist(raw string) {
	var (
		torrls []N.TorrItem
		items  []list.Item
	)
	json.Unmarshal([]byte(raw), &torrls)
	for i, t := range torrls {
		it := item{}
		it.title = strconv.Itoa(i) + "." + t.Name
		it.desc = t.UpdateTime + "|" + t.Size
		items = append(items, it)
	}
	m.list.Title = "Torrent List"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
}
