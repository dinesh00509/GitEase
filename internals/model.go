package internals

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Step struct {
	Label string
	Done  bool
}

type Model struct {
	Cursor              int
	Steps               []Step
	Output              string
	Committing          bool
	TextInput           textinput.Model
	BranchMode          bool
	NewBranch           bool
	PullBranch          bool
	PullFromOtherBranch bool
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter here..."
	ti.CharLimit = 120
	ti.Width = 50

	return Model{
		Cursor: 0,
		Steps: []Step{
			{"Stage all changes", false},
			{"Show Git status", false},
			{"Commit changes", false},
			{"Push to remote", false},
			{"Create new branch", false},
			{"Switch to another branch", false},
			{"Pull from current branch", false},
			{"Pull from other branch", false},
			{"list all branches", false},
			{"show the current branch", false},
			{"check the logs", false},
			{" check the logs with reflog", false},
		},
		TextInput: ti,
	}
}
