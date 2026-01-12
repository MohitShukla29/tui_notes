package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textInput          textinput.Model
	createInputVisible bool
}

func (m model) Init() tea.Cmd {
	return nil
}
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Write the command you want to execute!"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return model{
		textInput:          ti,
		createInputVisible: false,
	}
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+n":
			m.createInputVisible = true
			return m, nil
		case "enter":
			return m, nil
		}

	}
	if m.createInputVisible == true {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}
func (m model) View() string {
	welcome := "Welcome to tui Notes"
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("205")).
		PaddingRight(2).
		PaddingLeft(2)
	welc := style.Render(welcome)
	View := ""
	if m.createInputVisible {
		View = m.textInput.View()
	}
	help := fmt.Sprintf("%s\n\n%s\n\nCtrl+n: %s, Ctrl+L: %s, Esc: %s, Ctrl+S: %s, Ctrl+Q: %s", welc, View, "new file", "list", "back/save", "save", "quit")
	return help
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	fmt.Println("Welcome to tui")
}
