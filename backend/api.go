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

	router.HandleFunc("/feed/{tag}", makeHTTPHandleFunc(s.handleFeed))
	router.HandleFunc("/new/{type}", makeHTTPHandleFunc(s.handleNewThread))

	router.HandleFunc("/parentchild/{type}/{id}", makeHTTPHandleFunc(s.handleGetParentChild))

	router.HandleFunc("/{type}/{id}", makeHTTPHandleFunc(s.handleUserContent))

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

func (s *APIServer) handleGetParentChild(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return WriteJSON(w, http.StatusMethodNotAllowed, nil)
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ServerResponse{Resp: err.Error()})	
	}
	typ := mux.Vars(r)["type"]

	if typ == "thread" {
		posts, err := s.database.GetThreadPosts(id)
		if err != nil {
			return WriteJSON(w, http.StatusBadRequest, ServerResponse{Resp: err.Error()})	
		}
		return WriteJSON(w, http.StatusOK, posts)
	} else if typ == "post" {
		comments, err := s.database.GetPostComments(id)
		if err != nil {
			return WriteJSON(w, http.StatusBadRequest, ServerResponse{Resp: err.Error()})	
		}
		return WriteJSON(w, http.StatusOK, comments)
	}

	return WriteJSON(w, http.StatusNotFound, nil)
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

func (s *APIServer) handleUserContent(w http.ResponseWriter, r *http.Request) error {
	JWTerr := JWTAuth(w, r, s.database)
	if JWTerr != nil {
		return WriteJSON(w, http.StatusUnauthorized, nil)
	}	
	
	typ := mux.Vars(r)["type"]
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	userName, ckerr := r.Cookie("userName")
	if ckerr != nil {
		return ckerr
	}
	
	user := userName.Value
	author, err := s.database.GetUser(typ, id)
	if err != nil {
		return WriteJSON(w, http.StatusBadGateway, "something went wrong")
	}
	//Checking for correct User 
	if user != *author {
		return WriteJSON(w, http.StatusUnauthorized, "login required")
	}

	// switch typ {
	// 	case "thread": {
	// 		thread, err := s.database.GetThreadByID(id)
	// 			if err != nil { 
	// 				return WriteJSON(w, http.StatusBadGateway, "please retry") 
	// 			}
	// 			if thread[0].UserName != user { 
	// 				return WriteJSON(w, http.StatusUnauthorized, "login required")
	// 			}
	// 			if r.Method == "GET" {
	// 				return WriteJSON(w, http.StatusOK, thread[0])
	// 			}
	// 	}
	// 	case "post": {
	// 		post, err := s.database.GetPostByID(id)
	// 			if err != nil { 
	// 				return WriteJSON(w, http.StatusBadGateway, "please retry") 
	// 			}
	// 			if post[0].UserName != user { 
	// 				return WriteJSON(w, http.StatusUnauthorized, "login required")
	// 			}
	// 			if r.Method == "GET" {
	// 				return WriteJSON(w, http.StatusOK, post[0])
	// 			}
	// 	}
	// 	case "comment": {
	// 		comment, err := s.database.GetCommentByID(id)
	// 			if err != nil { 
	// 				return WriteJSON(w, http.StatusBadGateway, "please retry") 
	// 			}
	// 			if comment[0].UserName != user { 
	// 				return WriteJSON(w, http.StatusUnauthorized, "login required")
	// 			}
	// 	}
	// 	default:
	// 		return WriteJSON(w, http.StatusBadRequest, "wrong type")
	

	switch r.Method {
		case "GET":
			return s.handleGet(w, r)
		case "DELETE":
			return s.handleDelete(w, r)
		case "PATCH":
			return s.handlePatch(w, r)
		default:
			return WriteJSON(w, http.StatusMethodNotAllowed, nil)
	}
}

func (s *APIServer) handleGet(w http.ResponseWriter, r *http.Request) error {
	typ := mux.Vars(r)["type"]
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	if typ == "thread" {
		if thread, err := s.database.GetThreadByID(id); err != nil {
			return WriteJSON(w, http.StatusBadRequest, err.Error())
		} else {
			return WriteJSON(w, http.StatusOK, thread[0])
		}
	} else if typ == "post" {
		if post, err := s.database.GetPostByID(id); err != nil {
			return WriteJSON(w, http.StatusBadRequest, err.Error())
		} else {
			return WriteJSON(w, http.StatusOK, post[0])
		}
	} else if typ == "comment" {
		if comment, err := s.database.GetCommentByID(id); err != nil {
			return WriteJSON(w, http.StatusBadRequest, err.Error())
		} else {
			return WriteJSON(w, http.StatusOK, comment[0])
		}
	}

	return WriteJSON(w, http.StatusBadRequest, "wrong type")
}


func (s *APIServer) handleDelete(w http.ResponseWriter, r *http.Request) error {
	typ := mux.Vars(r)["type"]
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	delErr := s.database.Delete(typ, id)
	if delErr != nil {
		return WriteJSON(w, http.StatusBadRequest, ServerResponse{Resp: delErr.Error()})
	}

	return WriteJSON(w, http.StatusOK, "succesfully deleted")
}

func (s *APIServer) handlePatch(w http.ResponseWriter, r *http.Request) error {
	typ := mux.Vars(r)["type"]
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	req := new(PatchRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	patchErr := s.database.Update(req.Input1, req.Input2, typ, id)
	if patchErr != nil {
		return WriteJSON(w, http.StatusBadRequest, nil)
	}
	
	return WriteJSON(w, http.StatusOK, "succesfully updated")
}

// Helper functions
func WriteJSON(w http.ResponseWriter, status int, v any) error { //add error so it is compatible with all function signatures
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

func (a *Account) ValidatePassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.EncryptedPW), []byte(pw))

	return err == nil
}


