package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/t-murch/top-25-cli/pkg/common"
)

func ServerStartCmd(channels *common.Channels) {
	fmt.Print("Starting the list server....\n")

	babyTokenCache := make(map[string]common.TokenState)

	// Start a goroutine to listen for authorization details
	go func() {
		fmt.Print("Initializing preAuthChannel....\n")
		for {
			preAuth := <-channels.PreAuthChannel
			// Once authorization details are received, handle them
			handlePreAuthInfo(babyTokenCache, channels.TokenStateChannel, preAuth)
		}
	}()

	go func() {
		for {
			tokenStatus := <-channels.TokenStateChannel
			handleTokenUpdate(tokenStatus, babyTokenCache)
		}
	}()

	// Start the web server
	http.HandleFunc("/cli/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("handler func initialized...")
		// Extract code and state parameters from the query string
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		if hasError := r.URL.Query().Has("error"); hasError {
			// update babyTokenCache[state] to error state
			channels.TokenStateChannel <- common.TokenState{
				State:     state,
				ErrorInfo: r.URL.Query().Get("error"),
				Details:   common.AuthDetails{},
				HasError:  hasError,
			}
		}

		// update the babyTokenCache with the successful `code` value
		channels.TokenStateChannel <- common.TokenState{
			State:     state,
			ErrorInfo: "",
			Details:   common.AuthDetails{},
			HasError:  false,
		}

		// Send authorization details to the channel
		channels.PreAuthChannel <- common.PreAuthDetails{Code: code, State: state}

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

func handleTokenUpdate(tokenStatus common.TokenState, babyCache map[string]common.TokenState) {
	// if we already have a tokenState for the state, update it
	// and log the difference
	if _, ok := babyCache[tokenStatus.State]; ok {
		babyCache[tokenStatus.State] = tokenStatus
		log.Printf("updating cache key: %s. updated: %+v", tokenStatus.State, tokenStatus.Details)
	} else {
		babyCache[tokenStatus.State] = tokenStatus
		log.Printf("Initializing new cache key with: %s", tokenStatus.State)
	}
}

func handlePreAuthInfo(tokenCache map[string]common.TokenState, tokenChannel chan common.TokenState, preAuth common.PreAuthDetails) {
	fmt.Printf("received /spotifyAuth request\n")
	// Here you would handle the received authorization details
	// For example, make another API call using the received code
	// For simplicity, let's just print the details for now
	fmt.Println("Received authorization code:", preAuth.Code)
	fmt.Println("Received state:", preAuth.State)

	if len(preAuth.Code) > 0 && len(preAuth.State) > 0 {
		if _, ok := tokenCache[preAuth.State]; !ok {
			log.Printf("No token found for state: %s", preAuth.State)
		} else {

			client := &http.Client{}
			accessTokenUrl := "https://accounts.spotify.com/api/token"
			formValues := url.Values{}
			formValues.Set("code", preAuth.Code)
			formValues.Set("grant_type", "authorization_code")
			formValues.Set("redirect_uri", "http://localhost:8080/cli/callback")

			req, err := http.NewRequest("POST", accessTokenUrl, strings.NewReader(formValues.Encode()))
			if err != nil {
				log.Fatalf("Failed to construct request for access token. error: %s \n", err)
			}
			req.SetBasicAuth(os.Getenv("SPOT_CLIENT_ID"), os.Getenv("SPOT_CLIENT_ACCESS_KEY"))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			response, err := client.Do(req)
			if err != nil {
				log.Fatalf("Failed to gain access token. error: %s \n", err)
			}
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatalf("Failed to read acceTokenResponse.body. err: %s", err)
			}

			// Can we log the body here?
			log.Printf("Access Token Response Body: %s", body)

			// Parse the response body
			// and tokenChannel
			var accessTokenResponse common.AuthDetails
			err = json.Unmarshal(body, &accessTokenResponse)
			if err != nil {
				log.Fatalf("error parsing response with access token. error: %s", err)
			}

			log.Printf("Access Token Response: %+v", accessTokenResponse)

			tokenChannel <- common.TokenState{
				State:     preAuth.State,
				ErrorInfo: "",
				Details:   accessTokenResponse,
				HasError:  false,
			}
		}
	}

	// if hasError := preAuth.HasError; hasError {
	// 	tokenChannel <- common.TokenState{
	// 		State:     preAuth.State,
	// 		ErrorInfo: preAuth.ErrorInfo,
	// 		Details:   common.AuthDetails{},
	// 		HasError:  hasError,
	// 	}
	// } else {
	// 	tokenChannel <- common.TokenState{
	// 		State:     preAuth.State,
	// 		ErrorInfo: "",
	// 		Details:   common.AuthDetails{},
	// 		HasError:  hasError,
	// 	}
	// }
}
