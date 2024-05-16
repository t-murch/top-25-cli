package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/t-murch/top-25-cli/pkg/common"
)

// Spotify Green Foreground
var welcomeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#1DB954"))

type welcomeModel struct {
	err error
	// greeting string
}

func greetingModel() welcomeModel {
	return welcomeModel{}
}

func (m welcomeModel) Init() tea.Cmd {
	return nil
}

// invoke MainModel.Next() to change state
func (m welcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return nil, GoNext(ScopePage)
		}
	}

	return m, nil
}

func (m welcomeModel) View() string {
	var s string
	s += lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#1DB954")).Bold(true).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(lipgloss.Color("241")).Render("Welcome to the Top 25 CLI Tool"),
		"\n",
		lipgloss.NewStyle().Render("Press enter to continue.\n\n"),
		common.HelpStyle.Render("Press (esc) to quit.\n"),
	)
	return s
}
