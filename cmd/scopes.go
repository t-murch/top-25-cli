package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

/*
 * Agreement for scopes:
 * You agree that top-25-cli-app will be able to:
 * - View your Spotify account data
 * - Your name, username, profile picture, Spotify followers, and public playlists.
 */

type ScopesModel struct {
	err        error
	acceptance *huh.Form
}

func newScopesModel() ScopesModel {
	return ScopesModel{
		acceptance: huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Key("accept").
					Title("You agree to provide top-25-cli-app with the above permissions?").
					Affirmative("Yes").
					Negative("No (exit)"),
			),
		),
	}
}

func (m ScopesModel) Init() tea.Cmd {
	return m.acceptance.Init()
}

func (m ScopesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// case "enter":
		// 	return nil, GoNext(LoginPage)
		case "esc":
			return nil, GoBack(AppPage(WelcomePage))
		}
	}

	var cmds []tea.Cmd

	// form processing
	form, cmd := m.acceptance.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.acceptance = f
		cmds = append(cmds, cmd)
	}

	if m.acceptance.GetBool("accept") {
		return nil, GoNext(AppPage(LoginPage))
	}

	// handle form submission

	return m, tea.Batch(cmds...)
}

func (m ScopesModel) View() string {
	s := "In order to display and save your playlist data, we need certain permissions to your Spotify account.\n\n"
	// s += " will be able to:\n"
	s += "- View your Spotify account data\n"
	s += "- Your name, username, profile picture, Spotify followers, and public playlists.\n\n"
	// We will insert a toggle here to allow the user to accept or reject the scopes.
	s += m.acceptance.View() + "\n\n"
	// s += "Press (enter) to accept or (esc) to reject.\n"
	return s
}
