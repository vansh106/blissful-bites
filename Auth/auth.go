package Auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "net/http"
	// "firebase.google.com/go/v4/auth"
    // "google.golang.org/api/option"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth interface {
	Login(username, password string) bool
	Signup(username, password string) bool
}

type JSONAuth struct {
	users []User
}

func NewJSONAuth() *JSONAuth {
	return &JSONAuth{}
}

func (ja *JSONAuth) LoadUsers() error {
	data, err := ioutil.ReadFile("users.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &ja.users)
	if err != nil {
		return err
	}
	return nil
}

func (ja *JSONAuth) SaveUsers() error {
	data, err := json.Marshal(ja.users)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("users.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (ja *JSONAuth) Login(username, password string) bool {
	for _, user := range ja.users {
		if user.Username == username && user.Password == password {
			return true
		}
	}
	return false
}

func (ja *JSONAuth) Signup(username, password string) bool {
	for _, user := range ja.users {
		if user.Username == username {
			return false
		}
	}
	ja.users = append(ja.users, User{Username: username, Password: password})
	if err := ja.SaveUsers(); err != nil {
		fmt.Println("Error saving user:", err)
		return false
	}
	return true
}
