package common

// Shared comms channels
type Channels struct {
	PreAuthChannel    chan PreAuthDetails
	AuthChannel       chan AuthDetails
	TokenStateChannel chan TokenState
}

func NewChannels() *Channels {
	return &Channels{
		PreAuthChannel:    make(chan PreAuthDetails),
		AuthChannel:       make(chan AuthDetails),
		TokenStateChannel: make(chan TokenState),
	}
}

// PreAuthDetails stores the authorization details
type PreAuthDetails struct {
	Code      string
	State     string
	ErrorInfo string
	HasError  bool
}

type PreAuthError struct {
	Error string
	State string
}

type AuthDetails struct {
	Access_token  string `json:"access_token"`
	Token_type    string `json:"token_type"`
	Scope         string `json:"scope"`
	Refresh_token string `json:"refresh_token"`
	Expires_in    int    `json:"expires_in"`
}

type TokenState struct {
	State     string
	ErrorInfo string
	Details   AuthDetails
	HasError  bool
}
