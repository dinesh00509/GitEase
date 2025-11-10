package internals

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFFF")).Padding(0, 1)
	stepDoneStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	stepCursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	outputStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEFA")).MarginTop(1)
	successStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#32CD32")).Bold(true)
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
	divider         = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render(strings.Repeat("─", 30))
)

type step struct {
	label string
	done  bool
}

type model struct {
	cursor    int
	steps     []step
	output    string
	commiting bool
	textInput textinput.Model
}

func IntialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter your Message to Commit...."
	ti.Focus()
	ti.CharLimit = 120
	ti.Width = 40
	return model{
		cursor: 0,
		steps: []step{
			{"stage all the changes", false},
			{"show git status", false},
			{"commit the changes", false},
			{"Push to remote", false},
		},
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

// For viewing.
func (m model) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("GitEase") + "\n")
	b.WriteString(divider + "\n\n")

	for i, s := range m.steps {
		var cursor = " "
		if m.cursor == i {
			cursor = stepCursorStyle.Render("->")
		}
		status := "[ ]"
		if s.done {
			status = stepDoneStyle.Render("[✔]")
		}

		line := fmt.Sprintf("%s %s %s\n", cursor, status, s.label)
		b.WriteString(line)
	}
	b.WriteString("\n" + divider + "\n")

	if m.commiting {

		b.WriteString(stepCursorStyle.Render("Type your commit message and press Enter (ESC to cancel):") + "\n")
		b.WriteString(m.textInput.View())
	} else {
		if strings.Contains(m.output, "Error") {
			b.WriteString(errorStyle.Render(m.output))
		} else if strings.Contains(m.output, "completed") {
			b.WriteString(successStyle.Render(m.output))
		} else {
			b.WriteString(outputStyle.Render(m.output))
		}
	}
	b.WriteString("\n" + divider + "\n")
	b.WriteString(lipgloss.NewStyle().Faint(true).Render("Use ↑/↓ to navigate, Enter to run, q to quit.\n"))

	return b.String()
}

func (m model) runCurrentStep() (tea.Model, tea.Cmd) {
	m.output = fmt.Sprintf("Running: %s...\n", m.steps[m.cursor].label)
	switch m.cursor {
	case 0:
		m.output += RunGit("add", ".")
		m.output += "\nStaging completed successfully."

	case 1:
		m.output += RunGit("status")
		m.output += "Status check completed."

	case 2:
		m.commiting = true
		m.textInput.SetValue("")
		m.textInput.Focus()
		m.output = ""

	case 3:
		m.output += RunGit("push")
		m.output += "Pushed to remote branch successfully."
		m.output += "\nAll steps completed. Exiting GitFlow..."
		m.steps[m.cursor].done = true

		return m, tea.Tick(time.Second*2, func(time.Time) tea.Msg {
			return tea.QuitMsg{}
		})
	}

	if m.cursor != 2 {
		m.steps[m.cursor].done = true
	}

	return m, tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg { return nil })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.commiting {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				commitMsg := strings.TrimSpace(m.textInput.Value())
				if commitMsg != "" {
					m.output = RunGit("commit", "-m", commitMsg)
					m.steps[m.cursor].done = true
				} else {
					m.output = "Commit message cannot be empty."
				}
				m.commiting = false
				return m, nil
			case "esc":
				m.commiting = false
				m.output = "Commit Action Cancelled"
				return m, nil
			}
		}
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.steps)-1 {
				m.cursor++
			}

		case "enter":
			return m.runCurrentStep()
		}
	}
	return m, nil
}
