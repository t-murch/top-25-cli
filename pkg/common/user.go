package common

type User struct {
	Username    string
	authDetails AuthDetails
	HasAuth     bool
}

func NewUser() *User {
	return &User{}
}

func (u *User) SetAuthDetails(details AuthDetails) {
	u.authDetails = details
	u.HasAuth = true
}
