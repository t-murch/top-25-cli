package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type LoginModel struct {
	options  *huh.Form
	method   LoginMethod
	username string
	password string
}

type LoginMethod string

func (m LoginMethod) FilterValue() string { return "" }

const (
	facebook LoginMethod = "facebook"
	google   LoginMethod = "google"
	apple    LoginMethod = "apple"
	email    LoginMethod = "email"
)

func (m LoginModel) Init() tea.Cmd {
	return m.options.Init()
}

func newLoginModel() LoginModel {
	items := []LoginMethod{
		LoginMethod("Facebook"),
		LoginMethod("Google"),
		LoginMethod("Apple"),
		LoginMethod("Email"),
	}

	newModel := LoginModel{
		options: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[LoginMethod]().
					Key("method").
					Title("Select your Spotify Login Method").
					Options(huh.NewOptions(items...)...),
				huh.NewInput().Key("username").Title("Username"),
				huh.NewInput().Key("password").Title("Password").Password(true),
			),
		),
	}

	// return LoginModel{options: list.New(items, list.NewDefaultDelegate(), 14, 20)}
	return newModel
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	h, v := docStyle.GetFrameSize()
	// 	m.options.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return nil, GoBack(AppPage(ScopePage))
		}
	}

	form, cmd := m.options.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.options = f
		cmds = append(cmds, cmd)
	}

	if m.options.State == huh.StateCompleted {
		username := m.options.GetString("username")
		cmds = append(cmds, UpdateUserName(username))
		// Quit when the form is done.
		// cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

// This view will be a tabbed view of the different login methods
func (m LoginModel) View() string {
	return docStyle.Render(m.options.View())
}
