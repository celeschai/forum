package main

import (
	//"fmt"
	"time"
	"strconv"

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

func NewThread(title, username, tag string) (*Thread, error) {
	return &Thread{
		Title:    title,
		UserName: username,
		Tag:      tag,
		Created:  time.Now().UTC(),
	}, nil
}

func NewPost(threadid, username, title, content string) (*Post, error) {
	id, err := strconv.Atoi(threadid)
	if err != nil {
		return nil, err
	}

	return &Post{
		ThreadID: id,
		Title:    title,
		UserName: username,
		Content:  content,
		Created:  time.Now().UTC(),
	}, nil
}

func NewComment(postid string, username string, content string) (*Comment, error) {
	id, err := strconv.Atoi(postid)
	if err != nil {
		return nil, err
	}
	
	return &Comment{
		PostID:   id,
		UserName: username,
		Content:  content,
		Created:  time.Now().UTC(),
	}, nil
}
