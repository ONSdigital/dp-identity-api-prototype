package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/server"
	"github.com/gorilla/mux"
)

var bindAddr = ":23700"
var mockSessionID = "12345-abcde-67890-fghij"
var mockUserPassword = "one two three four"

type MockedUser struct {
	Email         string
	Password      string
	Migrated      bool
	EncryptionKey string
}

var MockedUsers = map[string]*MockedUser{
	"test@test.com": {
		Email:         "test@test.com",
		Password:      "one two three four",
		Migrated:      false,
		EncryptionKey: "encryption key for test@test.com",
	},
	"test-identity-api@email.com": {
		Email:         "test-identity-api@email.com",
		Password:      "one two three four",
		Migrated:      false,
		EncryptionKey: "encryption key for test-identity-api@email.com",
	},
}

func main() {
	r := mux.NewRouter()
	r.Path("/login").Methods("POST").HandlerFunc(loginHandler)
	// r.Path("/users").Methods("POST").HandlerFunc(userHandler)
	r.Path("/validate").Methods("POST").HandlerFunc(validateSessionHandler)
	r.Path("/migrate").Methods("POST").HandlerFunc(migrateUserHandler)

	log.Debug("Starting server", log.Data{
		"bind_addr": bindAddr,
	})
	s := server.New(bindAddr, r)
	if err := s.ListenAndServe(); err != nil {
		log.Error(err, nil)
		os.Exit(2)
	}
}

type LoginFields struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	ID string `json:"id"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	log.DebugR(r, "Handling login request", nil)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err, nil)
		return
	}
	defer r.Body.Close()

	var loginCreds LoginFields
	err = json.Unmarshal(data, &loginCreds)
	if err != nil {
		log.Error(err, nil)
		w.WriteHeader(400)
		w.Write([]byte("Unrecognised login data"))
		return
	}

	log.DebugR(r, "Checking whether login is valid", log.Data{
		"email": loginCreds.Email,
	})

	var user *MockedUser
	var ok bool

	if user, ok = MockedUsers[loginCreds.Email]; !ok {
		log.DebugR(r, "Email address not recognised", log.Data{
			"email": loginCreds.Email,
		})
		w.WriteHeader(401)
		w.Write([]byte("Email address not recognised"))
		return
	}

	if loginCreds.Password != user.Password {
		log.DebugR(r, "Login attempt failed due to incorrect password", log.Data{
			"email": user.Email,
		})
		w.WriteHeader(401)
		w.Write([]byte("Unable to login because password is incorrect"))
		return
	}

	if user.Migrated == false {
		log.DebugR(r, "User is not migrated", log.Data{
			"email": user.Email,
		})
		w.WriteHeader(303)
		w.Write([]byte("Unable to login because user is not migrated"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("i-am-an-access-token"))
}

// func userHandler(w http.ResponseWriter, r *http.Request) {

// }

type ValidSessionResponse struct {
	Email    string `json:"email"`
	Type     string `json:"type"`
	Password string `json:"password"`
}

func validateSessionHandler(w http.ResponseWriter, r *http.Request) {
	log.DebugR(r, "Validating session ID", nil)
	log.DebugR(r, r.Header.Get("Content-type"), nil)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err, nil)
		return
	}
	defer r.Body.Close()

	var session Session
	err = json.Unmarshal(data, &session)
	if err != nil {
		log.Error(err, nil)
		w.WriteHeader(400)
		w.Write([]byte("Unrecognised data structure"))
		return
	}

	if session.ID == "" {
		w.WriteHeader(400)
		w.Write([]byte("Must contain a session ID"))
		return
	}

	if session.ID != mockSessionID {
		w.WriteHeader(401)
		return
	}

	response := ValidSessionResponse{
		Email:    "crispin.merriman@gmail.com",
		Type:     "ADMIN",
		Password: mockUserPassword,
	}
	json, err := json.Marshal(response)

	log.DebugR(r, "Returning valid access token", nil)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func migrateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.DebugR(r, "Migrating user", nil)
	err := r.ParseForm()
	if err != nil {
		log.Error(err, nil)
	}

	username := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println("USERNAME:", username, "PASSWORD:", password)

	if user, ok := MockedUsers[username]; ok {
		log.Debug("User found", nil)
		user.Migrated = true
		user.Password = password
		w.Write([]byte(user.EncryptionKey))
		return
	}
	log.Debug("User not found", nil)
	w.WriteHeader(http.StatusNotFound)
}
