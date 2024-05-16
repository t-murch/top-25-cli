package cmd

//
// func initialModel() model {
// 	date := time.Now().Format("02 Jan 2006")
//
// 	ti := textinput.New()
// 	ti.Placeholder = fmt.Sprintf("Top 25 Playlist Name - %s\n", date)
// 	ti.Focus()
// 	ti.CharLimit = 50
// 	ti.Width = 50
//
// 	return model{
// 		textInput: ti,
// 		err:       nil,
// 	}
// }
//
// func (m model) View() string {
// 	// The header
// 	s := "What would you like the name of your Top 25 Playlist to be?\n\n"
// 	s += "Press enter to use the default name.\n\n"
// 	s += m.textInput.View() + "\n\n"
// 	s += "Press (esc) to quit.\n"
//
// 	return s
// }
