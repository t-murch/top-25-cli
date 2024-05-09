package services

import (
	"fmt"
	"net/http"
	"os"
)

// PreAuthDetails stores the authorization details
type PreAuthDetails struct {
	Code  string
	State string
}

func ServerStartCmd() {
	fmt.Print("Starting the list server....\n")
	// Create a channel to communicate the authorization details
	authChannel := make(chan PreAuthDetails)

	// Start a goroutine to listen for authorization details
	go func() {
		fmt.Print("Initializing authChannel....\n")
		for {
			auth := <-authChannel
			// Once authorization details are received, handle them
			handlePreAuthInfo(auth)
		}
	}()

	// Start the web server
	http.HandleFunc("/cli/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("handler func initialized...")
		// Extract code and state parameters from the query string
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		// Send authorization details to the channel
		authChannel <- PreAuthDetails{Code: code, State: state}

		// Respond to the browser
		// w.Write([]byte("Authorization successful! You can close this window."))
	})

	fmt.Print("Just before go func....\n")
	// Start the server
	go func() {
		fmt.Print("Initializing ListenAndServe...\n")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Server error:", err)
			os.Exit(1)
		}
		fmt.Println("Web Server Started...")
	}()
}

// var ServerStartCmd = &cobra.Command{
// 	Use:   "server_start",
// 	Short: "description",
// 	Long:  ".",
//
// 	Run: func(cmd *cobra.Command, args []string) {
// 	},
// }

func handlePreAuthInfo(preAuth PreAuthDetails) {
	fmt.Printf("received /spotifyAuth request\n")
	// Here you would handle the received authorization details
	// For example, make another API call using the received code
	// For simplicity, let's just print the details for now
	fmt.Println("Received authorization code:", preAuth.Code)
	fmt.Println("Received state:", preAuth.State)
}
