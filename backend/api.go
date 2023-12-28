package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/signup", makeHTTPHandleFunc(s.handleSignup))
	router.HandleFunc("/account/{accnum}", makeHTTPHandleFunc(s.handleUser))

	//router.HandleFunc("/feed", makeHTTPHandleFunc(s.handleFeed))
	router.HandleFunc("/feed/{tag}", makeHTTPHandleFunc(s.handleFeed))

	router.HandleFunc("/thread/{threadid}", makeHTTPHandleFunc(s.handleThreads))
	router.HandleFunc("/post/{postid}", makeHTTPHandleFunc(s.handlePosts))
	router.HandleFunc("/comment/{commentid}", makeHTTPHandleFunc(s.handleComments))

	log.Println("JSON API server running on port", s.listenAddr)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:" + frontend},
		AllowCredentials: true,
		Debug: true,
	})
	handler := c.Handler(router)
	http.ListenAndServe(s.listenAddr, handler)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf(r.Method, "method not allowed for login")
	}

	if JWTAuth(w, r, s.database) == nil {
		return WriteJSON(w, http.StatusOK, LoginResponse{Resp: "Welcome Back!"})
	} 
	//*** needs to be incorporated into feed page, if user is logged in, no need to jump to sign in page

	req := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.database.GetAccountByEmail(req.Email)
	switch {
	case err != nil:
		return WriteJSON(w, http.StatusUnauthorized,
			LoginResponse{Resp: err.Error()})
	case !acc.ValidatePassword(req.Password):
		return WriteJSON(w, http.StatusUnauthorized,
			LoginResponse{Resp: "invalid password"})
	}

	token, err := createJWT(acc)
	if err != nil {
		return err
	}
	setCookie(w, r, "jwtToken", token)
	setCookie(w, r, "userName", acc.UserName)

	return WriteJSON(w, http.StatusOK, LoginResponse{Resp: "succesful login"})
}

func (s *APIServer) handleSignup(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf(r.Method, "method not allowed for signup")
	}
	return WriteJSON(w, http.StatusOK, LoginResponse{Resp: nil})
}

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	//check for JWT
	return nil
}

func (s *APIServer) handleFeed(w http.ResponseWriter, r *http.Request) error {
	//check for JWT
	return nil
}

func (s *APIServer) handleFilter(w http.ResponseWriter, r *http.Request) error {
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

// Helper functions
func WriteJSON(w http.ResponseWriter, status int, v any) error { //add error so it is compatible with all function signatures
	//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:" + frontend)
	//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Println(status, v)
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
