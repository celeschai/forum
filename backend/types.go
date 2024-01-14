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

type CreateThreadRequest struct {
	Title string `json:"title"`
	Tag   string `json:"tag"`
}

type CreatePostRequest struct {
	ThreadID string `json:"threadid"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type CreateCommentRequest struct {
	PostID  string `json:"postid"`
	Content string `json:"content"`
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
	Likes    string    `json:"likes"`
}

type Comment struct {
	CommentID int       `json:"id"`
	PostID    int       `json:"postid"`
	UserName  string    `json:"username"`
	Content   string    `json:"content"`
	Created   time.Time `json:"created"`
}

type PatchRequest struct {
	Input1 string `json:"input1"`
	Input2 string `json:"input2"`
}

type Database interface {
	//new account
	CreateAccount(*Account) error //check if account exists
	CheckExistingAcc(*CreateAccountRequest) error
	
	//new content
	CreateThread(*Thread) error
	CreatePost(*Post) error
	CreateComment(*Comment) error
	
	//feed
	GetLatestThreads(string) ([]*Thread, error)
	GetPostsByThreadID(int) ([]*Post, error)
	GetThreadByID(int) ([]*Thread, error)
	GetPostByID(int) ([]*Post, error)
	GetCommentByID(int) ([]*Comment, error)
	GetCommentsByPostID(int)([]*Comment, error)
	
	//parentchild
	GetThreadPosts(int, string) (map[string]interface{}, error)
	GetPostComments(int) (map[string]interface{}, error)

	//user
	GetThreadsByUser(string) ([]*Thread, error)
	GetPostsByUser(string) ([]*Post, error)
	GetCommentsByUser(string) ([]*Comment, error)

	//accounts
	GetAccountByUserName(string) (*Account, error)
	GetAccountByEmail(string) (*Account, error)
	GetAccUploads(string) (map[string]interface{}, error)
	GetUser(string, int) (*string, error)
	Update(input1, input2, typ string, id int) error
	Delete(string, int) error
}

type PostgresStore struct {
	db *sql.DB //easier change of database if need be
}
