package main

import (
	"database/sql"
	"net/http"
	"time"
)

type Account struct {
	UserName    string    `json:"username"`
	Email       string    `json:"email"`
	EncryptedPW string    `json:"-"`
	Created     time.Time `json:"created"`
}

// http requests and responses
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Resp any `json:"resp"`
}

type CreateAccountRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// handling API & connection
type APIServer struct {
	listenAddr string
	database   Database
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

// content
type Thread struct {
	ThreadID int       `json:"id"`
	Title    string    `json:"title"`
	UserName string    `json:"username"`
	Tag      string    `json:"tag"`
	Created  time.Time `json:"created"`
}

type Post struct {
	PostID   int       `json:"id"`
	ThreadID int       `json:"threadid"`
	Title    string    `json:"title"`
	UserName string    `json:"username"`
	Content  string    `json:"content"`
	Created  time.Time `json:"created"`
}

type Comment struct {
	CommentID int `json:"id"`
	PostID    int `json:"postid"`
	//ThreadID  int       `json:"threadid"`
	//Title     string    `json:"title"`
	UserName string    `json:"username"`
	Content  string    `json:"content"`
	Created  time.Time `json:"created"`
}

// // database
// type method string

type Database interface {
	CreateAccount(*Account) error //check if account exists
	GetAccountByUserName(string) (*Account, error)
	GetAccountByEmail(string) (*Account, error)
	//UpdateAccount(*Account) error
	//DeleteAccountByID(int) error

	GetLatestThreads(string) ([]*Thread, error)
	CreateThread(*Thread) error
	// CreateThread(*Thread) error
	// DeleteThread(*Thread) error - maybe cannot delete threads
	GetThreadPosts(int) (map[string]interface{}, error)

	//storePosts(*Post, method) error
	CreatePost(*Post) error
	// DeletePost(*Post) error
	// EditPost(*Post) error +++
	//GetComments() ([]*Comment, error)

	//storeComments(*Comment, method) error
	CreateComment(*Comment) error
	// DeleteComment(*Comment) error
	// EditComment(*Comment) error +++
}

// check connection: nc -vz localhost 5432
type PostgresStore struct {
	db *sql.DB
}
