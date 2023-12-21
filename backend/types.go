package main

import (
	"database/sql"
	"net/http"
	"time"
)

type Account struct {
	//AccNum      int       `json:"accnum"`
	UserName    string    `json:"username"`
	Email       string    `json:"email"`
	UserID      int64     `json:"userid"`
	EncryptedPW string    `json:"-"`
	CreatedAt   time.Time `json:"createdAt"`
}

// http requests and responses 
type LoginRequest struct {
	UserID   int64  `json:"userid"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID int64  `json:"userid"`
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
	ThreadID int64  `json:"threadid"`
	Title    string `json:"title"`
	UserID   int64  `json:"userid"`
	Tag1     string `json:"tag1"`
	Tag2     string `json:"tag2"`
}

type Post struct {
	PostID   int64     `json:"postid"`
	ThreadID int64     `json:"threadid"`
	Title    string    `json:"title"`
	UserID   int64     `json:"userid"`
	Content  string    `json:"content"`
	Date     time.Time `json:"date"`
}

type Comment struct {
	CommentID int64     `json:"commentid"`
	PostID    int64     `json:"postid"`
	//ThreadID  int64     `json:"threadid"`
	Title     string    `json:"title"`
	UserID    int64     `json:"userid"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"`
}

// database
type Database interface {
	CreateAccount(*Account) error
	GetAccountByUserID(int64) (*Account, error)
	//GetAccountByAccNum(int) (*Account, error)      
	//UpdateAccount(*Account) error
	//DeleteAccountByID(int) error

	NewThread(*Thread) error
	DeleteThread(*Thread) error
	GetPosts() ([]*Post, error)

	NewPost(*Post)
	DeletePost(*Post) error
	EditPost(*Post) error
	GetComments() ([]*Comment, error)

	NewComment(*Comment) error
	DeleteComment(*Comment) error
	EditComment(*Comment) error
}

// docker run --name postgres -e POSTGRES_PASSWORD=winter -p 5432:5432 -d postgres
// check connection: nc -vz localhost 5432
type PostgresStore struct {
	db *sql.DB
}
