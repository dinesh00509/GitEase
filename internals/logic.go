package internals

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd { return nil }

func (m Model) RunCurrentStep() (tea.Model, tea.Cmd) {
	m.Output = fmt.Sprintf("Running: %s...\n", m.Steps[m.Cursor].Label)

	switch m.Cursor {

	case 0:
		m.Output += RunGit("add", ".")
		m.Output += "\nStaging completed successfully."

	case 1:
		m.Output += RunGit("status")
		m.Output += "\nStatus check completed."

	case 2:
		m.Committing = true
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		m.Output = ""

	case 3:
		m.Output += RunGit("push")
		m.Output += "\nPush completed successfully.\nAll steps done — exiting..."
		m.Steps[m.Cursor].Done = true
		return m, tea.Tick(time.Second*2, func(time.Time) tea.Msg { return tea.QuitMsg{} })

	case 4:
		m.BranchMode = true
		m.NewBranch = true
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		m.Output = ""

	case 5:
		m.BranchMode = true
		m.NewBranch = false
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		m.Output = ""

	case 6:
		m.PullBranch = true
		m.PullFromOtherBranch = false
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		m.Output = ""

	case 7:
		m.PullBranch = true
		m.PullFromOtherBranch = true
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		m.Output = ""

	case 8:
		m.Output = RunGit("branch", "-a")
		m.Output += "\nList of branches completed successfully."
		m.Steps[m.Cursor].Done = true

	case 9:
		m.Output = RunGit("branch")
		m.Output += "\nCurrent branch completed successfully."
		m.Steps[m.Cursor].Done = true

	case 10:
		m.Output = RunGit("log")
		m.Steps[m.Cursor].Done = true
	case 11:
		m.Output = RunGit("reflog")
		m.Steps[m.Cursor].Done = true
	}

	if m.Cursor != 2 && m.Cursor != 4 && m.Cursor != 5 && m.Cursor != 6 && m.Cursor != 7 {
		m.Steps[m.Cursor].Done = true
	}

	return m, tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg { return nil })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// Commit message input handling
	if m.Committing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				commitMsg := strings.TrimSpace(m.TextInput.Value())
				if commitMsg != "" {
					m.Output = RunGit("commit", "-m", commitMsg)
					m.Output += "\nCommit created successfully."
					m.Steps[m.Cursor].Done = true
				} else {
					m.Output = "Commit message cannot be empty."
				}
				m.Committing = false
				return m, nil
			case "esc":
				m.Committing = false
				m.Output = "Commit action cancelled."
				return m, nil
			}
		}
		m.TextInput, cmd = m.TextInput.Update(msg)
		return m, cmd
	}

	// Branch mode handling
	if m.BranchMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				branchName := strings.TrimSpace(m.TextInput.Value())
				if branchName == "" {
					m.Output = "Branch name cannot be empty."
				} else {
					if m.NewBranch {
						m.Output = RunGit("checkout", "-b", branchName)
						m.Output += "\nNew branch created successfully."
					} else {
						m.Output = RunGit("checkout", branchName)
						m.Output += "\nSwitched to branch successfully."
					}
					m.Steps[m.Cursor].Done = true
				}
				m.BranchMode = false
				return m, nil
			case "esc":
				m.BranchMode = false
				m.Output = "Branch action cancelled."
				return m, nil
			}
		}
		m.TextInput, cmd = m.TextInput.Update(msg)
		return m, cmd
	}

	// Pull branch handling
	if m.PullBranch {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if m.PullFromOtherBranch {
					branchName := strings.TrimSpace(m.TextInput.Value())
					if branchName == "" {
						m.Output = "Branch name cannot be empty."
						return m, nil
					}
					m.Output = RunGit("pull", "origin", branchName)
					m.Output += "\nPull completed successfully."
					m.Steps[m.Cursor].Done = true
				} else {
					// Pull from current branch
					m.Output = RunGit("pull")
					m.Output += "\nPull completed successfully."
					m.Steps[m.Cursor].Done = true
				}
				m.PullBranch = false
				m.PullFromOtherBranch = false
				return m, nil
			case "esc":
				m.PullBranch = false
				m.PullFromOtherBranch = false
				m.Output = "Pull action cancelled."
				return m, nil
			}
		}
		// Only update text input if we're pulling from another branch
		if m.PullFromOtherBranch {
			m.TextInput, cmd = m.TextInput.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	// Navigation & actions
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down":
			if m.Cursor < len(m.Steps)-1 {
				m.Cursor++
			}
		case "enter":
			return m.RunCurrentStep()
		}
	}

	return m, nil
}
