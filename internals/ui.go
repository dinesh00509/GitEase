package internals

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FFFF")).
			Background(lipgloss.Color("#1E1E1E")).
			Padding(0, 2).
			MarginBottom(1).
			Underline(true)

	subtleDivider = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444")).
			Render(strings.Repeat("─", 50))

	stepPendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	stepActiveStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	stepDoneStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF7F")).Bold(true)

	outputBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#555555")).
			Padding(1, 2).
			MarginTop(1).
			Width(65)

	successText = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF7F")).Bold(true)
	errorText   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	hintText    = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Faint(true).MarginTop(1)
)

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("GitEase – Git CLI Assistant") + "\n")
	b.WriteString(subtleDivider + "\n\n")

	for i, s := range m.Steps {
		var line string
		if s.Done {
			line = stepDoneStyle.Render(fmt.Sprintf("✔ %s", s.Label))
		} else if i == m.Cursor {
			line = stepActiveStyle.Render(fmt.Sprintf("▶ %s", s.Label))
		} else {
			line = stepPendingStyle.Render(fmt.Sprintf("• %s", s.Label))
		}
		b.WriteString(line + "\n")
	}
	b.WriteString("\n" + subtleDivider + "\n\n")

	if m.Committing {
		b.WriteString(stepActiveStyle.Render("Commit Message:") + "\n")
		b.WriteString(m.TextInput.View())
	} else if m.BranchMode {
		if m.NewBranch {
			b.WriteString(stepActiveStyle.Render("Enter new branch name:") + "\n")
		} else {
			b.WriteString(stepActiveStyle.Render("Enter branch name to switch:") + "\n")
		}
		b.WriteString(m.TextInput.View())
	} else if m.PullBranch {
		if m.PullFromOtherBranch {
			b.WriteString(stepActiveStyle.Render("Enter branch name to pull from:") + "\n")
			b.WriteString(m.TextInput.View())
		} else {
			b.WriteString(stepActiveStyle.Render("Press Enter to pull from current branch") + "\n")
			b.WriteString(m.TextInput.View())
		}
	} else {
		if m.Output == "" {
			b.WriteString(outputBox.Render("Waiting for command..."))
		} else if strings.Contains(strings.ToLower(m.Output), "error") {
			b.WriteString(outputBox.Render(errorText.Render(m.Output)))
		} else if strings.Contains(strings.ToLower(m.Output), "completed") ||
			strings.Contains(strings.ToLower(m.Output), "success") {
			b.WriteString(outputBox.Render(successText.Render(m.Output)))
		} else {
			b.WriteString(outputBox.Render(m.Output))
		}
	}

	b.WriteString("\n\n" + subtleDivider + "\n")
	b.WriteString(hintText.Render("↑/↓ navigate • Enter run • q quit • ESC cancel input") + "\n")

	return b.String()
}
