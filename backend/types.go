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
	ThreadID string    `json:"threadid"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type CreateCommentRequest struct {
	PostID  string    `json:"postid"`
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

// type Post struct {
// 	PostID   int       `json:"id"`
// 	ThreadID int       `json:"threadid"`
// 	Title    string    `json:"title"`
// 	UserName string    `json:"username"`
// 	Content  string    `json:"content"`
// 	Created  time.Time `json:"created"`
// 	Likes    string    `json:"likes"`
// 	Liked    string    `json:"liked"`
// }

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

// copy over all function signatures
type Database interface {
	CreateAccount(*Account) error //check if account exists
	CheckExistingAcc(acc *CreateAccountRequest) error
	GetAccountByUserName(string) (*Account, error)
	GetAccountByEmail(string) (*Account, error)
	GetAccUploads(string) (map[string]interface{}, error)
	GetUser(typ string, id int) (*string, error)
	//UpdateAccount(*Account) error
	//DeleteAccountByID(int) error

	GetLatestThreads(string) ([]*Thread, error)
	CreateThread(*Thread) error
	GetThreadPosts(id int, user string) (map[string]interface{}, error)
	GetThreadsByUser(userName string) ([]*Thread, error)
	GetThreadByID(id int) ([]*Thread, error)

	CreatePost(*Post) error
	GetPostComments(id int) (map[string]interface{}, error)
	GetPostsByUser(userName string) ([]*Post, error)
	GetPostsByThreadID(id int) ([]*Post, error)
	GetPostByPostID(id int) ([]*Post, error)

	// ScanLikedP(row *sql.Rows) (*UserLikedPost, error)
	// ScanLikedPosts(row *sql.Rows) ([]*UserLikedPost, error)

	CreateComment(*Comment) error
	GetCommentByID(id int) ([]*Comment, error)

	Delete(typ string, id int) error
	Update(input1, input2, typ string, id int) error
	//Like(username string, postid int, like bool) error
	//IsLiked(username string) ([]*int, error)
}

// check connection: nc -vz localhost 5432
type PostgresStore struct {
	db *sql.DB //easier change of database if need be
}
