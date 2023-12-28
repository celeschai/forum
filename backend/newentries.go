package main

import (
	//"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// func InvalidInput(s interface{}) error {
// 	return fmt.Errorf("invalid input: %v", s)
// }

func NewAPIServer(listenAddr string, database Database) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		database:   database, //access to database
	}
}

func NewAccount(username, email, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		UserName:    username,
		Email:       email,
		EncryptedPW: string(encpw),
		Created:     time.Now().UTC(),
	}, nil
}

func NewThread(username string, title, tag string) (*Thread, error) {
	return &Thread{
		Title:    title,
		UserName: username,
		Tag:      tag,
		Created:  time.Now().UTC(),
	}, nil
}

func NewPost(threadid int, username string, title, content string) (*Post, error) {
	return &Post{
		ThreadID: threadid,
		Title:    title,
		UserName: username,
		Content:  content,
		Created:  time.Now().UTC(),
	}, nil
}

func NewComment(postid int, username string, content string) (*Comment, error) {
	return &Comment{
		PostID:   postid,
		UserName: username,
		Content:  content,
		Created:  time.Now().UTC(),
	}, nil
}
