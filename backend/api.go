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
	//"golang.org/x/tools/go/analysis/passes/nilfunc"
)

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/home", makeHTTPHandleFunc(s.handleJWT))
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/signup", makeHTTPHandleFunc(s.handleSignup))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

	//router.HandleFunc("/feed", makeHTTPHandleFunc(s.handleFeed))
	router.HandleFunc("/feed/{tag}", makeHTTPHandleFunc(s.handleFeed))
	router.HandleFunc("/newthread", makeHTTPHandleFunc(s.handleNewThread))

	router.HandleFunc("/threadposts/{threadid}", makeHTTPHandleFunc(s.handleGetThreadPosts))
	router.HandleFunc("/postcomments/{postid}", makeHTTPHandleFunc(s.handleGetPostComments))

	router.HandleFunc("/delete/{type}/{id}", makeHTTPHandleFunc(s.handleDelete))
	//router.HandleFunc("/comment/{commentid}", makeHTTPHandleFunc(s.handleComments))

	log.Println("JSON API server running on port", s.listenAddr)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:" + frontend + "*"},
		AllowCredentials: true,
		Debug: false,
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions, http.MethodPatch},
	})
	handler := c.Handler(router)
	http.ListenAndServe(s.listenAddr, handler)
}

func (s *APIServer) handleJWT(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return WriteJSON(w, http.StatusMethodNotAllowed, nil)
	}

	if JWTAuth(w, r, s.database) == nil {
		return WriteJSON(w, http.StatusOK, "Welcome Back!")
	}

	return WriteJSON(w, http.StatusUnauthorized, nil)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return WriteJSON(w, http.StatusMethodNotAllowed, nil)
	}

	req := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.database.GetAccountByEmail(req.Email)
	switch {
	case err != nil:
		return WriteJSON(w, http.StatusUnauthorized, nil)
	case !acc.ValidatePassword(req.Password):
		return WriteJSON(w, http.StatusUnauthorized, nil)
	}

	token, err := createJWT(acc)
	if err != nil {
		return err
	}
	setCookie(w, r, "jwtToken", token)
	setCookie(w, r, "userName", acc.UserName)

	return WriteJSON(w, http.StatusOK, ServerResponse{Resp: "succesful login"})
}


func (s *APIServer) handleSignup(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return WriteJSON(w, http.StatusMethodNotAllowed, nil)
	}

	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	exist := s.database.CheckExistingAcc(req)
	if exist != nil {
		return WriteJSON(w, http.StatusConflict, ServerResponse{Resp: exist.Error()})
	}

	acc, err := NewAccount(req.UserName, req.Email, req.Password)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ServerResponse{Resp: "failed to create account, please try again"})
	}

	if err := s.database.CreateAccount(acc); err != nil {
		return WriteJSON(w, http.StatusBadRequest, ServerResponse{Resp: err.Error()})
	}
	
	return WriteJSON(w, http.StatusOK, ServerResponse{Resp: "succesful signup"})
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	err := JWTAuth(w, r, s.database)
	if err != nil {
		return WriteJSON(w, http.StatusUnauthorized, ServerResponse{Resp: "please log in again."})
	}

	cookie, cookerr := r.Cookie("userName")
	if cookerr != nil {
		return err
	}
	username := cookie.Value

	acc, err := s.database.GetAccUploads(username)
	if err != nil {
		return err
	}
	acc["username"]=username

	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIServer) handleFeed(w http.ResponseWriter, r *http.Request) error {
	err := JWTAuth(w, r, s.database)
	if err != nil {
		return WriteJSON(w, http.StatusUnauthorized, ServerResponse{Resp: "please log in again."})
	}

	if r.Method != "GET" {
		return WriteJSON(w, http.StatusMethodNotAllowed, nil)
	}

	tag := mux.Vars(r)["tag"]
	threads, err := s.database.GetLatestThreads(tag)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ServerResponse{Resp: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, threads)
}

func (s *APIServer) handleGetThreadPosts(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return WriteJSON(w, http.StatusMethodNotAllowed, nil)
	}

	threadid, err := strconv.Atoi(mux.Vars(r)["threadid"])
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ServerResponse{Resp: err.Error()})	
	}

	posts, err := s.database.GetThreadPosts(threadid)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ServerResponse{Resp: err.Error()})	
	}

	return WriteJSON(w, http.StatusOK, posts)
}

func (s *APIServer) handleNewThread(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return WriteJSON(w, http.StatusMethodNotAllowed, nil)
	}

	err := JWTAuth(w, r, s.database)
	if err != nil {
		return WriteJSON(w, http.StatusUnauthorized, ServerResponse{Resp: err.Error()})
	}

	req := new(Thread)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	cookie, cookerr := r.Cookie("userName")
	if cookerr != nil {
		return err
	}
	username := cookie.Value

	thread, err := NewThread(req.Title, username, req.Tag)
	if err != nil {	
		return err
	}

	err = s.database.CreateThread(thread)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, thread)
}

func (s *APIServer) handleGetPostComments(w http.ResponseWriter, r *http.Request) error {
	//check for JWT
	return nil
}

func (s *APIServer) handleComments(w http.ResponseWriter, r *http.Request) error {
	//check for JWT
	return nil
}

func (s *APIServer) handleDelete(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "DELETE" {
		return fmt.Errorf(r.Method, "method not allowed for deleting")
	}

	err := JWTAuth(w, r, s.database)
	if err != nil {
		return WriteJSON(w, http.StatusUnauthorized, ServerResponse{Resp: err.Error()})
	}

	typ := mux.Vars(r)["type"]
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}
	fmt.Println(typ, id)

	delErr := s.database.Delete(typ, id)
	if delErr != nil {
		return WriteJSON(w, http.StatusBadRequest, ServerResponse{Resp: delErr.Error()})
	}

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

func getPath(r *http.Request, pathInput string) (string, error) {
	inputStr := mux.Vars(r)[pathInput]
	
	
	return inputStr, nil
}

func (a *Account) ValidatePassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.EncryptedPW), []byte(pw))

	return err == nil
}
