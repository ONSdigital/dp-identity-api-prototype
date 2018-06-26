package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/server"
	"github.com/gorilla/mux"
)

var bindAddr = ":23700"
var mockAccessToken = "12345-abcde-67890-fghij"
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

var MockedAccessTokens = map[string]string{
	"12345-abcde-67890-fghij": "test-identity-api@email.com",
}

func main() {
	r := mux.NewRouter()
	r.Path("/login").Methods("POST").HandlerFunc(loginHandler)
	// r.Path("/users").Methods("POST").HandlerFunc(userHandler)
	r.Path("/validate").Methods("POST").HandlerFunc(validateAccessTokenHandler)
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

type AccessToken struct {
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
	w.Write([]byte(mockAccessToken))
}

// func userHandler(w http.ResponseWriter, r *http.Request) {

// }

type ValidAccessTokenResponse struct {
	Email         string `json:"email"`
	EncryptionKey string `json:"encryptionKey"`
	// Type     string `json:"type"`
}

func validateAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.DebugR(r, "Validating access token", nil)
	// log.DebugR(r, r.Header.Get("Content-type"), nil)

	// data, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	log.Error(err, nil)
	// 	return
	// }
	// defer r.Body.Close()§

	err := r.ParseForm()
	if err != nil {
		log.Error(err, nil)
	}

	accessToken := AccessToken{
		ID: r.FormValue("id"),
	}
	// err = json.Unmarshal(data, &accessToken)
	// if err != nil {
	// 	log.Error(err, nil)
	// 	w.WriteHeader(400)
	// 	w.Write([]byte("Unrecognised data structure"))
	// 	return
	// }

	if accessToken.ID == "" {
		w.WriteHeader(400)
		w.Write([]byte("Must contain an access token"))
		return
	}

	if accessToken.ID != mockAccessToken {
		w.WriteHeader(401)
		return
	}

	var email string
	var ok bool
	if email, ok = MockedAccessTokens[accessToken.ID]; !ok {
		w.WriteHeader(401)
		w.Write([]byte("Invalid access token"))
		return
	}

	var user *MockedUser
	if user, ok = MockedUsers[email]; !ok {
		w.WriteHeader(500)
		w.Write([]byte("Error finding email address for this access token"))
		return
	}

	response := ValidAccessTokenResponse{
		Email:         user.Email,
		EncryptionKey: user.EncryptionKey,
	}
	json, err := json.Marshal(response)

	log.DebugR(r, "Returning valid encryption key", nil)
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
