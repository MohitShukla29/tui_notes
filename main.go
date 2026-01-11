package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	msg string
}

func (m model) Init() tea.Cmd {
	return nil
}
func initialModel() model {
	return model{
		msg: "🍆",
	}
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
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

	help := fmt.Sprintf("%s\n\nCtrl+N: %s, Ctrl+L: %s, Esc: %s, Ctrl+S: %s, Ctrl+Q: %s", welc, "new file", "list", "back/save", "save", "quit")
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
