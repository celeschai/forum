package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/signup", makeHTTPHandleFunc(s.handleSignup))
	router.HandleFunc("/{userid}", makeHTTPHandleFunc(s.handleUser))

	router.HandleFunc("/feed", makeHTTPHandleFunc(s.handleFeed))
	router.HandleFunc("/{threadid}", makeHTTPHandleFunc(s.handleThreads))
	router.HandleFunc("/{postid}", makeHTTPHandleFunc(s.handlePosts))
	router.HandleFunc("/{commentid}", makeHTTPHandleFunc(s.handleComments)) //using POST for privacy - GET shows in history of browser

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	//check header for cookies
	return nil
}

func (s *APIServer) handleSignup(w http.ResponseWriter, r *http.Request) error {
	//check header for cookies
	return nil
}

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	//check for JWT
	return nil
}

func (s *APIServer) handleFeed(w http.ResponseWriter, r *http.Request) error {
	//check for JWT
	return nil
}

func (s *APIServer) handleThreads(w http.ResponseWriter, r *http.Request) error {
	//check for JWT
	return nil
}

func (s *APIServer) handlePosts(w http.ResponseWriter, r *http.Request) error {
	//check for JWT
	return nil
}

func (s *APIServer) handleComments(w http.ResponseWriter, r *http.Request) error {
	//check for JWT
	return nil
}

//Helper functions
func WriteJSON(w http.ResponseWriter, status int, v any) error { //add error so it is compatible with all function signatures
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request, idtype string) (int, error) {
	idStr := mux.Vars(r)[idtype]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

func (a *Account) ValidatePassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.EncryptedPW), []byte(pw))

	return err == nil
}
