package common

type User struct {
	Username    string
	authDetails AuthDetails
	HasAuth     bool
}

// type UserUpdateValue interface {
// 	string | bool | AuthDetails
// }

type UserName struct {
	Username string
}

func NewUser() *User {
	return &User{}
}

func (u *User) SetAuthDetails(details AuthDetails) {
	u.authDetails = details
	u.HasAuth = true
}
