package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.HiddenBorder()).
	BorderForeground(lipgloss.Color("240"))

type TorrItem struct {
	Name       string `json:"name"`
	Size       string `json:"size"`
	UpdateTime string `json:"uptime"`
}

type Subj struct {
	Sid     int    `json:"sid"`
	Typ     string `json:"type"`
	Name    string `json:"name"`
	Episode string `json:"epi"`
	Status  string `json:"status"`
	Compl   string `json:"compl"`
}

func (m *model) tableUpdate(msg tea.Msg, istorr bool) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mod = text
			return m, cmd
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if istorr {
				m.ninput = "add " + getArg(m.history[len(m.history)-1]) + "-i " + m.table.SelectedRow()[0]
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
		torrls []TorrItem
		subjls []Subj
	)
	t := "torrls"
	err := json.Unmarshal([]byte(raw), &torrls)
	if err != nil || len(torrls) == 0 || torrls[0].Size == "" {
		err = json.Unmarshal([]byte(raw), &subjls)
		if err != nil || len(subjls) == 0 || subjls[0].Sid == 0 {
			t = "unkonwn"
		} else {
			t = "subjls"
		}
	}
	switch t {
	case "torrls":
		columns := []table.Column{
			{Title: "index", Width: 4},
			{Title: "name", Width: 10},
			{Title: "size", Width: 4},
			{Title: "uptime", Width: 4},
		}
		var tr []table.Row
		for i, torr := range torrls {
			tr = append(tr, []string{strconv.Itoa(i), torr.Name, torr.Size})
		}
		return newTable(columns, tr), true, nil
	case "subjls":
		columns := []table.Column{
			{Title: "sid", Width: 4},
			{Title: "type", Width: 4},
			{Title: "name", Width: 10},
			{Title: "epi", Width: 4},
			{Title: "status", Width: 4},
			{Title: "compl", Width: 4},
		}
		var tr []table.Row
		for _, subj := range subjls {
			tr = append(tr, []string{strconv.Itoa(subj.Sid), subj.Typ, subj.Name, subj.Episode, subj.Status, subj.Compl})
		}
		return newTable(columns, tr), false, nil

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
