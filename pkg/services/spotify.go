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
	"github.com/t-murch/top-25-cli/pkg/common"
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

func GrantAuthForUser(channels *common.Channels, strategy string) {
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

	preAuthRequestUrl := buildAccessCodeRequest()
	preAuthResponse, err := page.Goto(preAuthRequestUrl)
	if err != nil || preAuthResponse.Status() > 399 {
		log.Fatalf("could not goto: %v, status: %d, error: %v \n", preAuthRequestUrl, preAuthResponse.Status(), err)
		// I dont think I need to initialize data in the tokenCache in error state.
		stateIdentifier := parseAccessRequestParams(preAuthRequestUrl, "state")
		errMsgOrStatus := fmt.Sprint(preAuthResponse.Status())
		if err != nil {
			errMsgOrStatus = err.Error()
		}

		channels.PreAuthChannel <- common.PreAuthDetails{
			Code:      "",
			State:     stateIdentifier,
			ErrorInfo: errMsgOrStatus,
			HasError:  true,
		}
	} else {
		stateIdentifier := parseAccessRequestParams(preAuthRequestUrl, "state")
		channels.PreAuthChannel <- common.PreAuthDetails{
			Code:      "",
			State:     stateIdentifier,
			ErrorInfo: "",
			HasError:  false,
		}
		// channels.TokenStateChannel <- common.TokenState{
		// 	State:     stateIdentifier,
		// 	ErrorInfo: "",
		// 	Details:   common.AuthDetails{},
		// 	HasError:  false,
		// }
		// log.Printf("initialized new babyTokenCache key: %s \n", stateIdentifier)
	}

	// page.Pause()
	/*
			* Workflow:
			* - Go to url of `buildAccessCodeRequest()`
		*  - Select Auth Strategy
		*  - Login to that strategy
		*  - Redirect back to Spotify PreAuth Page
		*  - Click `Agree`
		*  - Await call to `/cli/callback` for code && state values
			* */

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
	if err := page.WaitForURL("https://www.facebook.com/login**"); err != nil {
		log.Fatalf("Failed to navigate to facebook auth page. Error: %v", err)
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

		// preAuthRedirectTimeout := float64(10)
		page.Locator("#loginbutton").Click()
		// page.Pause()
		// https://accounts.spotify.com/en/authorize?client_id=035b4603779340a785097df893e33728&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcli%2Fcallback&response_type=code&show_dialog=true&state=6ARzqcK_Vua3g89F3TaR_g%3D%3D
		if err := page.WaitForURL("https://accounts.spotify.com/en/authorize?**"); err != nil {
			log.Fatalf("Failed to redirect back to PreAuth for Acceptance. Error: %v \n", err)
		}

		// callbackTimeout := float64(10)
		page.GetByTestId("auth-accept").Click()
		if err := page.WaitForURL("http://localhost:8080/cli/callback**"); err != nil {
			log.Fatalf("Failed to login. Check username and password. Error: %v \n", err)
		}

		submitBtn, err := page.Locator("loginbutton").Count()
		if err != nil {
			log.Fatalf("error logging in. Error: %v \n", err)
		}

		if submitBtn < 1 {
			log.Print("Input login details and hit enter successfully. ")
		}
	}

	// page.Pause()
	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}

	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}

func buildAccessCodeRequest() string {
	myRandomString, err := GenerateRandomString(16)
	if err != nil {
		log.Fatal("Failed to generate random string. wtf... \n")
	}

	data := url.Values{}
	data.Set("client_id", os.Getenv("SPOT_CLIENT_ID"))
	data.Set("response_type", "code")
	data.Set("state", myRandomString)
	// data.Set("client_secret", os.Getenv("SPOT_CLIENT_SECRET"))
	// TODO: Change this once working.
	data.Set("redirect_uri", "http://localhost:8080/cli/callback")
	// data.Set("grant_type", "authorization_code")
	data.Set("show_dialog", "true")

	tokenUrl := AUTH_URL

	fmt.Printf("auth url: %v\n", tokenUrl+"?"+data.Encode())
	return tokenUrl + "?" + data.Encode()
}

func parseAccessRequestParams(fullUrl string, param string) string {
	parsedUrl, err := url.Parse(fullUrl)
	if err != nil {
		log.Fatalf("failed to parse our own url. should not happennnnn. err: %v \n", err)
	}

	allParams := parsedUrl.Query()

	if hasParam := allParams.Has(param); !hasParam {
		log.Fatalf("param: %s not found in url: %s \n", param, parsedUrl)
	}

	if stateParams := allParams[param]; len(stateParams) > 1 {
		log.Fatalf("unable to parse param: %s. Multiple values. url: %s \n", param, parsedUrl)
	}

	return allParams.Get(param)
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
