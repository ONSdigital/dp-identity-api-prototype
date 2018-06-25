package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

func main() {
	r := mux.NewRouter()
	r.Path("/login").Methods("POST").HandlerFunc(loginHandler)
	// r.Path("/users").Methods("POST").HandlerFunc(userHandler)
	r.Path("/validate").Methods("POST").HandlerFunc(validateSessionHandler)

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
	log.DebugR(r, "Logging in user", nil)
	buf := new(bytes.Buffer)
	session := Session{
		ID: mockSessionID,
	}
	b, err := json.Marshal(session)
	if err != nil {
		log.Error(err, nil)
		return
	}

	err = binary.Write(buf, binary.BigEndian, &b)
	if err != nil {
		log.Error(err, nil)
		return
	}

	log.DebugR(r, "Returning valid session ID", nil)
	w.Write(b)
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
