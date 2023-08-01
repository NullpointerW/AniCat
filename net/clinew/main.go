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

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// "github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
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
	p := tea.NewProgram(initialModel())
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

type model struct {
	welcome   bool
	history   []string
	spinner   spinner.Model
	textInput textinput.Model
	table     table.Model
	torrls    bool
	mod       mode
	ninput    string
	recvtb    bool
	// cmd       string
	recv bool
	err  error
}

func initialModel() *model {
	ti := textinput.New()
	ti.Placeholder = "anicat-cliv2"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 70

	// magenta red `>`
	ti.Prompt = wordwrap.String("\x1B[38;2;249;38;114m>\x1B[0m", 0)

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
			reply = "\n" + reply
			m.history = append(m.history, reply)
			m.textInput.Focus()
			m.mod = text
			return m, nil
		default:
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
mod:
	switch m.mod {
	case tb:
		// first into
		if m.recvtb {
			atb, istorr, err := NewTable(reply)
			if err != nil {
				// ls
			} else {
				m.torrls = istorr
				m.table = atb
				m.recvtb = false
			}
		}
		return m.tableUpdate(msg, m.torrls)
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
			if t := getCmdTyp(input); t == Ls || t == Lsi { // tb
				m.recvtb = true
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
		default:
			last := strings.Join(m.history, "\n")
			return last + "\n" +
				m.textInput.View()
		}
	}
}
