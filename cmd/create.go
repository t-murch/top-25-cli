package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/t-murch/top-25-cli/pkg/common"
)

/* COBRA */
func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "short description",
	Long:  ".",

	Run: func(cmd *cobra.Command, args []string) {
		if f, err := tea.LogToFile("debug.log", "help"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer func() {
				err = f.Close()
				if err != nil {
					log.Fatal(err)
				}
			}()
		}

		p := tea.NewProgram(newMainModel())
		// p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Oh poop, we have an Error: %v", err)
			os.Exit(1)
		}

		// Wait for termination signal (Ctrl+C)
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
	},
}

/* BUBBLETEA */
var (
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	tipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#190")).Italic(true)
	endingMsgStyle = lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("170")).Bold(true)
)

type sessionState int

const (
	welcomeView sessionState = iota
	scopeView
	loginView
	playlistsView
)

type AppPage int

const (
	WelcomePage   AppPage = 0
	ScopePage     AppPage = 1
	LoginPage     AppPage = 2
	PlaylistsPage AppPage = 3
)

type GoTo struct {
	Page AppPage
}

func GoNext(page AppPage) tea.Cmd {
	return func() tea.Msg {
		return GoTo{Page: page}
	}
}

func GoBack(page AppPage) tea.Cmd {
	return func() tea.Msg {
		return GoTo{Page: page}
	}
}

type MainModel struct {
	User   *common.User
	scopes ScopesModel
	login  LoginModel
	state  sessionState
	width  int
	height int
}

type (
	errMsg error
)

func newMainModel() *MainModel {
	return &MainModel{
		state:  welcomeView,
		scopes: newScopesModel(),
		login:  newLoginModel(),
		User:   common.NewUser(),
	}
}

func (m MainModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		switch m.state {
		case welcomeView:
		// m.welcome.Update(msg)
		case scopeView:
		// m.scope.Update(msg)
		case loginView:
			_, cmd := m.login.Update(msg)
			cmds = append(cmds, cmd)
		}

	case GoTo:
		m.state = sessionState(msg.Page)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		//
		// case errMsg:
		// 	m.err = msg
		// 	return m, nil
	}

	switch m.state {
	case welcomeView:
		wv := greetingModel()
		_, cmd := wv.Update(msg)
		cmds = append(cmds, cmd)
	case scopeView:
		// sv := newScopesModel()
		// _, cmd := sv.Update(msg)
		_, cmd := m.scopes.Update(msg)
		cmds = append(cmds, cmd)
	case loginView:
		_, cmd := m.login.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var currentView string

	switch m.state {
	case welcomeView:
		currentView = greetingModel().View()
	case scopeView:
		currentView = m.scopes.View()
	case loginView:
		currentView = m.login.View()
	}

	return currentView
}
