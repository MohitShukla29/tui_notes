package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	vaultDir string
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

type model struct {
	textInput          textinput.Model
	createInputVisible bool
	currentFile        *os.File
	noteTextarea       textarea.Model
	list               list.Model
	showingList        bool
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m model) Init() tea.Cmd {
	return nil
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory", err)
	}
	vaultDir = fmt.Sprintf("%s/.tui_notes", homeDir)
}
func initialModel() model {
	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatal(err)
	}

	ti := textinput.New()
	ti.Placeholder = "Write the command you want to execute!"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	ta := textarea.New()
	ta.Placeholder = "Write your note here"
	ta.Focus()

	noteList := listFiles()
	finalList := list.New(noteList, list.NewDefaultDelegate(), 0, 0)
	finalList.Title = "Created notes 📚"
	finalList.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("254")).
		Padding(0, 1)

	return model{
		textInput:          ti,
		createInputVisible: false,
		noteTextarea:       ta,
		list:               finalList,
	}
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-5)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+n":
			m.createInputVisible = true
			return m, nil
		case "ctrl+l":
			m.showingList = true
			return m, nil
		case "ctrl+s":
			if m.currentFile == nil {
				break
			}
			if err := m.currentFile.Truncate(0); err != nil {
				fmt.Println("Error saving the file")
				return m, nil
			}
			if _, err := m.currentFile.Seek(0, 0); err != nil {
				fmt.Println("Error saving the file")
				return m, nil
			}
			if _, err := m.currentFile.WriteString(m.noteTextarea.Value()); err != nil {
				fmt.Println("Error saving the file")
				return m, nil
			}
			if err := m.currentFile.Close(); err != nil {
				fmt.Println("cannot close the file")
			}
			m.currentFile = nil
			m.noteTextarea.SetValue("")
			return m, nil
		case "enter":
			if m.currentFile != nil {
				break
			}
			filename := m.textInput.Value()
			filepath := fmt.Sprintf("%s/%s.md", vaultDir, filename)
			if filepath != "" {
				if _, err := os.Stat(filepath); err == nil {
					return m, nil
				}

				f, err := os.Create(filepath)
				if err != nil {
					log.Fatalf("%v", err)
				}
				m.currentFile = f
				m.createInputVisible = false
				m.textInput.SetValue("")
			}
			return m, nil
		}

	}
	if m.createInputVisible == true {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	if m.currentFile != nil {
		m.noteTextarea, cmd = m.noteTextarea.Update(msg)
	}
	if m.showingList {
		m.list, cmd = m.list.Update(msg)
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

	if m.currentFile != nil {
		View = m.noteTextarea.View()
	}
	if m.showingList {
		View = m.list.View()
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

func listFiles() []list.Item {
	items := make([]list.Item, 0)
	entries, err := os.ReadDir(vaultDir)
	if err != nil {
		log.Fatal("Error in reading files")
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			modTime := info.ModTime().Format("2006-01-02 15:04")

			items = append(items, item{
				title: entry.Name(),
				desc:  fmt.Sprintf("Last modified at:%s", modTime),
			})
		}
	}
	return items
}
