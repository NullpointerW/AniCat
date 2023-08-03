package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	N "github.com/NullpointerW/anicat/net"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.HiddenBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m *model) tableUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mod = text
			m.textInput.Focus()
			return m, cmd
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.isstatls {
				m.isstatls = false
				m.mod = text
				return m, cmd
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *model) tableView() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func NewTable(raw string) (table.Model, bool, error) {
	var (
		statls []N.Stat
		subjls []N.Subj
	)
	t := "statls"
	err := json.Unmarshal([]byte(raw), &statls)
	if err != nil || len(statls) == 0 || statls[0].Size == "" {
		err = json.Unmarshal([]byte(raw), &subjls)
		if err != nil || len(subjls) == 0 || subjls[0].Sid == 0 {
			t = "unkonwn"
		} else {
			t = "subjls"
		}
	}
	switch t {
	case "statls":
		columns := []table.Column{
			{Title: "file", Width: 30},
			{Title: "size", Width: 10},
			{Title: "fin", Width: 12},
		}
		var tr []table.Row
		for _, stat := range statls {
			tr = append(tr, []string{stat.File, stat.Size, stat.Progress})
		}
		tb := newTable(columns, tr)
		// tb.Blur()
		return tb, false, nil
	case "subjls":
		columns := []table.Column{
			{Title: "sid", Width: 12},
			{Title: "type", Width: 10},
			{Title: "name", Width: 30},
			{Title: "epi", Width: 4},
			{Title: "status", Width: 10},
			{Title: "compl", Width: 10},
		}
		var tr []table.Row
		for _, subj := range subjls {
			tr = append(tr, []string{strconv.Itoa(subj.Sid), subj.Typ, subj.Name, subj.Episode,
				subj.Status, subj.Compl})
		}
		tb := newTable(columns, tr)
		// tb.Blur()
		return tb, false, nil

	default:
		return table.Model{}, false, fmt.Errorf("unkonwn raw type")
	}
}

func newTable(c []table.Column, r []table.Row) table.Model {
	t := table.New(
		table.WithColumns(c),
		table.WithRows(r),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.HiddenBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	return t
}
