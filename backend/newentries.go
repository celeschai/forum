package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func InvalidInput(s interface{}) error {
	return fmt.Errorf("invalid input: %v", s)
}

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

func NewThread(userid int, title, tag1, tag2 string) (*Thread, error) {
	return &Thread{
		Title:   title,
		UserID:  userid,
		Tag1:    tag1,
		Tag2:    tag2,
		Created: time.Now().UTC(),
	}, nil
}

func NewPost(threadid, userid int, title, content string) (*Post, error) {
	return &Post{
		ThreadID: threadid,
		Title:    title,
		UserID:   userid,
		Content:  content,
		Created:  time.Now().UTC(),
	}, nil
}

func NewComment(postid, userid int, content string) (*Comment, error) {
	return &Comment{
		PostID:  postid,
		UserID:  userid,
		Content: content,
		Created: time.Now().UTC(),
	}, nil
}
