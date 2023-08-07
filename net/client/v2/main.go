package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"

	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "github.com/charmbracelet/lipgloss"
	// "github.com/muesli/reflow/wordwrap"
)

var (
	host    string
	port    int
	conn    *bufio.Scanner
	tcpConn net.Conn
	err     error
)

func init() {
	flag.StringVar(&host, "h", "localhost", "server dial host")
	flag.IntVar(&port, "p", 8080, "server dial port")
	flag.Parse()
}

func main() {
	conn, tcpConn, err = connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
	mode   int
)

const (
	text = mode(iota)
	loading
	tb
	ls
)

// debug
var errmsg string

type model struct {
	welcome                   bool
	history                   []string
	spinner                   spinner.Model
	textInput                 textinput.Model
	table                     table.Model
	mod                       mode
	recv, recvtb, recvls      bool
	err                       error
	list, chdlist             list.Model
	istrls, ischdls, isstatls bool
}

func initialModel() *model {
	ti := textinput.New()
	ti.Placeholder = "anicat-cliv2"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 70

	// magenta red `>`
	// ti.Prompt = wordwrap.String("\x1B[38;2;249;38;114m>\x1B[0m", 0)
	ti.PromptStyle = ti.PromptStyle.Foreground(lipgloss.Color("#F92672"))

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &model{
		spinner:   s,
		textInput: ti,
		err:       nil,
	}
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) replyAppend(reply string) {
	reply = "\n" + reply
	m.history = append(m.history, reply)
	m.textInput.Focus()
	m.mod = text
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd   tea.Cmd
		reply string
	)

	if m.recv {
		select {
		case reply = <-send:
			m.recv = false
			if m.recvtb {
				m.mod = tb
				goto mod
			}
			if m.recvls {
				m.mod = ls
				goto mod
			}

			m.replyAppend(reply)
			return m, nil

		default:
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
mod:
	switch m.mod {
	case tb:
		// first into tb
		if m.recvtb {
			m.recvtb = false
			atb, isstat, err := NewTable(reply)
			if err != nil {
				m.replyAppend(reply)
				return m, nil
			} else { // ls
				m.isstatls = isstat
				m.table = atb
			}
		}
		return m.tableUpdate(msg)
	case ls:
		// first into ls
		if m.recvls {
			m.recvls = false
			typ := lsiReturnTyp(reply)
			if typ == er {
				// m.replyAppend(reply)
				m.replyAppend(reply)
				return m, nil
			}
			if typ == torr {
				m.istrls = true
				m.NewTorrlist(reply)
			} else {
				m.NewRsslist(reply)
			}
		}
		if m.istrls { // torrent list
			return m.torrlsUpdate(msg)
		} else if m.ischdls { // rss child
			return m.rssChildlsUpdate(msg)
		} else { // rss
			return m.rsslsUpdate(msg)
		}

	default:
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			input := m.textInput.Value()
			if input == "cls" || input == "clear" {
				m.history = nil
				m.textInput.Reset()
				return m, nil
			}
			switch t := getCmdTyp(input); t {
			case Ls, Stat: // tb
				m.recvtb = true
			case Lsi:
				m.recvls = true
			}

			m.welcome = true
			m.textInput.Value()
			m.history = append(m.history, input)

			recv <- input
			m.recv = true
			m.mod = loading
			m.textInput.Blur()
			m.textInput.Reset()

			return m, m.spinner.Tick
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if !m.welcome {
		return fmt.Sprintln("AniCat Cliv2") + m.textInput.View()
	} else {
		switch m.mod {
		case loading:
			return fmt.Sprintf("\n\n   %s Loading \n\n", m.spinner.View())
		case tb:
			return m.tableView()
		case ls:
			return m.listView()
		default:
			var last string
			if len(m.history) == 0 {
				return m.textInput.View()
			} else {
				last = strings.Join(m.history[:len(m.history)-1], "\n")
				last = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(last)
				return last + "\n" + "\n" +
				m.history[len(m.history)-1] + "\n" + "\n" +
				m.textInput.View()
			}
		}
	}
}

func (m *model) sendCmd(cmd string) {
	m.history = append(m.history, cmd)
	recv <- cmd
	m.recv = true
	m.mod = loading
	m.textInput.Blur()
	m.textInput.Reset()
}
