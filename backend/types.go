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

type ServerResponse struct {
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

type PatchRequest struct {
	Input1 string `json:"input1"`
	Input2 string `json:"input2"`
}

type Database interface {
	CreateAccount(*Account) error //check if account exists
	CheckExistingAcc(acc *CreateAccountRequest) error
	GetAccountByUserName(string) (*Account, error)
	GetAccountByEmail(string) (*Account, error)
	GetAccUploads(string) (map[string]interface{}, error)
	//UpdateAccount(*Account) error
	//DeleteAccountByID(int) error

	GetLatestThreads(string) ([]*Thread, error)
	CreateThread(*Thread) error
	// CreateThread(*Thread) error
	// DeleteThread(*Thread) error - maybe cannot delete threads
	GetThreadPosts(int) (map[string]interface{}, error)
	GetThreadsByUser(userName string) ([]*Thread, error)

	//storePosts(*Post, method) error
	CreatePost(*Post) error
	// DeletePost(*Post) error
	// EditPost(*Post) error +++
	GetPostComments(id int) (map[string]interface{}, error)
	GetPostsByUser(userName string) ([]*Post, error)

	//storeComments(*Comment, method) error
	CreateComment(*Comment) error
	// DeleteComment(*Comment) error
	// EditComment(*Comment) error +++

	Delete(typ string, id int) error
	Update(input1, input2, typ string, id int) error
}

// check connection: nc -vz localhost 5432
type PostgresStore struct {
	db *sql.DB //easier change of database if need be
}

