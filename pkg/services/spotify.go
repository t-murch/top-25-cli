package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
)

var (
	AUTH_URL            = "https://accounts.spotify.com/authorize"
	SCOPES              = "playlist-modify-private playlist-modify-public playlist-read-private user-read-email user-read-private user-top-read"
	PlaylistDescription = fmt.Sprintf("Built by I Miss My Top 25 by https://github.com/t-murch - %s", time.Now().Format("2006-Jan-02"))
)

/*
* Program inputs needed:
* - Login Strategy
*   - Login info; email/username and password
*   - Suggest using environment vars
*   - Or we can request the sensitive info to be set at ENV VARS we specify for the app to read.
*   $(TOP_25_CLI_USERNAME)
*   $(TOP_25_CLI_PASSWORD)
* */

func GrantAuthForUser(strategy string) {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Failed to launch Playwright: %v \n", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(false)})
	if err != nil {
		log.Fatalf("Failed to launch browser instance: %v \n", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v \n", err)
	}

	if _, err = page.Goto(buildAccessCodeRequest()); err != nil {
		log.Fatalf("could not goto: %v \n", err)
	}

	/**
				* AUTH OPTIONS
		* - Google: button data-testId="google-login"
		* - Facebook: button data-testId="facebook-login"
	  * - - FB Login Page
	  * - - - Email/phone: id="email" name="email"
	  * - - - Pass: id="pass" name="pass"
	  * - - - Submit: id="loginbutton" name="login"
		* - Apple: button data-testId="apple-login"
		* - Email/Pass
		* - - Email: input data-testid="login-username"
		* - - Pass: input data-testid="login-password"
		* - - Submit: button data-testid="login-button"
				* */
	if err = page.GetByTestId("facebook-login").Click(); err != nil {
		log.Fatalf("could not goto: %v \n", err)
	}

	// page.Pause()

	if strategy == "facebook" {
		username, usernameSet := os.LookupEnv("TOP_25_CLI_USERNAME")
		password, passwordSet := os.LookupEnv("TOP_25_CLI_PASSWORD")

		if !usernameSet || !passwordSet {
			log.Fatalf("Unable to grant auth due to lack of Username or Password Environment Vars set \n")
		}

		loginField := page.Locator("#email")
		passwordField := page.Locator("#pass")

		loginField.Fill(username)
		passwordField.Fill(password)

		page.Locator("#loginbutton").Click()

		log.Print("Input login details and hit enter successfully. ")
	}

	// entries, err := page.Locator(".athing").All()
	// if err != nil {
	// 	log.Fatalf("could not get entries: %v", err)
	// }

	// for i, entry := range entries {
	// 	title, err := entry.Locator("td.title > span > a").TextContent()
	// 	if err != nil {
	// 		log.Fatalf("could not get text content: %v", err)
	// 	}
	// 	fmt.Printf("%d: %s\n", i+1, title)
	// }

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}

	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}

func buildAccessCodeRequest() string {
	data := url.Values{}
	data.Set("client_id", os.Getenv("SPOT_CLIENT_ID"))
	data.Set("response_type", "code")
	// data.Set("client_secret", os.Getenv("SPOT_CLIENT_SECRET"))
	// TODO: Change this once working.
	data.Set("redirect_uri", "http://localhost:8080/cli/callback")
	// data.Set("grant_type", "authorization_code")
	data.Set("show_dialog", "true")

	tokenUrl := AUTH_URL

	fmt.Printf("auth url: %v", tokenUrl+"?"+data.Encode())
	return tokenUrl + "?" + data.Encode()
}

func GenerateRandomBytes(num int) ([]byte, error) {
	bytes := make([]byte, num)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func GenerateRandomString(num int) (string, error) {
	bytes, err := GenerateRandomBytes(num)
	return base64.URLEncoding.EncodeToString(bytes), err
}
