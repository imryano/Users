package userPackage

import (
	"fmt"
)

//Enum
type AccessLevel int

const (
	AccessLevelUser        AccessLevel = iota
	AccessLevelContributor AccessLevel = iota
	AccessLevelAdmin       AccessLevel = iota
)

//go:generate stringer -type=AccessLevel

//Aggregate
type User struct {
	ID              int
	Username        string
	Password        string
	Email           string
	AccessLevel     AccessLevel
	ExpectedVersion int
	Changes         []interface{}
}

func (u User) String() string {
	format := `
		ID: 			%d
		Username: 		%s
		Password: 		%s
		Email: 			%s
		AccessLevel: 		%s

		Expected Version: 	%d
		Pending Changes: 	%d
	`

	return fmt.Sprintf(format, u.ID, u.Username, u.Password, u.Email, u.AccessLevel, u.ExpectedVersion, len(u.Changes))
}

//Events
type CreateUser struct {
	Username string
	Password string
	Email    string
}

type PromoteUser struct{}

//Event Handling
func NewUserFromHistory(events []interface{}) *User {
	state := &User{}
	for _, event := range events {
		state.Transition(event)
		state.ExpectedVersion++
	}

	return state
}

func (state *User) trackChange(event interface{}) {
	state.Changes = append(state.Changes, event)
	state.Transition(event)
}

func (self *User) PromoteUser() {
	self.trackChange(PromoteUser{})
}

//Perform event on the User object
func (state *User) Transition(event interface{}) {
	switch e := event.(type) {
	case CreateUser:
		state.ID = 1
		state.Username = e.Username
		state.Password = e.Password
		state.Email = e.Email
		state.AccessLevel = AccessLevelUser

	case PromoteUser:
		if state.AccessLevel < AccessLevelAdmin {
			state.AccessLevel += 1
		}
	}
}
