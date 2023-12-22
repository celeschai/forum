package main

import (
	"database/sql"
	"net/http"
	"time"
)

type Account struct {
	UserID      int       `json:"useridx"`
	UserName    string    `json:"username"`
	Email       string    `json:"email"`
	EncryptedPW string    `json:"-"`
	Created     time.Time `json:"created"`
}

// http requests and responses
type LoginRequest struct {
	UserID   int    `json:"userid"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID int    `json:"userid"`
	Token  string `json:"token"`
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
	ThreadID int       `json:"threadid"`
	Title    string    `json:"title"`
	UserID   int       `json:"userid"`
	Tag1     string    `json:"tag1"`
	Tag2     string    `json:"tag2"`
	Created  time.Time `json:"created"`
}

type Post struct {
	PostID   int       `json:"postid"`
	ThreadID int       `json:"threadid"`
	Title    string    `json:"title"`
	UserID   int       `json:"userid"`
	Content  string    `json:"content"`
	Created  time.Time `json:"created"`
}

type Comment struct {
	CommentID int       `json:"commentid"`
	PostID    int       `json:"postid"`
	//ThreadID  int       `json:"threadid"`
	//Title     string    `json:"title"`
	UserID    int       `json:"userid"`
	Content   string    `json:"content"`
	Created   time.Time `json:"created"`
}

// // database
// type method string

type Database interface {
	CreateAccount(*Account) (error) //check if account exists
	//GetAccountByUserID(int) (*Account, error)
	//GetAccountByAccNum(int) (*Account, error)
	//UpdateAccount(*Account) error
	//DeleteAccountByID(int) error

	CreateThread(*Thread) (error)
	// CreateThread(*Thread) error
	// DeleteThread(*Thread) error - maybe cannot delete threads
	//GetPosts() ([]*Post, error)

	//storePosts(*Post, method) error
	CreatePost(*Post) (error)
	// DeletePost(*Post) error
	// EditPost(*Post) error +++
	//GetComments() ([]*Comment, error)

	//storeComments(*Comment, method) error
	CreateComment(*Comment) (error)
	// DeleteComment(*Comment) error
	// EditComment(*Comment) error +++
}

// check connection: nc -vz localhost 5432
type PostgresStore struct {
	db *sql.DB
}
