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

	N "github.com/NullpointerW/anicat/net"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
)

type model struct {
	welcome   bool
	history   []string
	textInput textinput.Model
	// reply     string	
	err error
}

func initialModel() *model {
	ti := textinput.New()
	ti.Placeholder = "anicat-cliv2"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 70

	ti.Prompt = wordwrap.String("\x1B[38;2;249;38;114m>\x1B[0m", 0)
	// ti.PromptStyle = lipgloss.NewStyle().

	// 	Foreground(lipgloss.Color("#FAFAFA")).
	// 	Background(lipgloss.Color("#7D56F4"))

	return &model{
		textInput: ti,
		err:       nil,
	}
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			input := m.textInput.Value()

			m.welcome = true
			m.textInput.Value()
			m.history = append(m.history, input)
			tcpConn.Write([]byte(input + N.CRLF))
			var reply string
			if scanok := conn.Scan(); scanok {
				reply = conn.Text()
			} else {
				reply = conn.Err().Error()
			}
			reply = "\n" + reply
			m.history = append(m.history, reply)
			m.textInput.Reset()
			return m, nil

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
		last := strings.Join(m.history, "\n")
		return last + "\n" +
			m.textInput.View()
	}
}
