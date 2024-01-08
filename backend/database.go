package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// docker run --name forum -e POSTGRES_USER=forumadmin -e POSTGRES_PASSWORD=gossiping -e POSTGRES_DB=forum -p 5432:5432 -d postgres

func NewPostgressStore() (*PostgresStore, error) {
	connStr := "user=forumadmin dbname=forum password=gossiping sslmode=disable" 

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createTables()
}

func (s *PostgresStore) createTables() error {
	query :=
		`CREATE TABLE IF NOT EXISTS users(
			username VARCHAR(50),
			email VARCHAR(50),
			encryptedpw VARCHAR(100),
			created TIMESTAMP,

			CONSTRAINT pk_account PRIMARY KEY (username)
		);
		CREATE TABLE IF NOT EXISTS threads(
			threadid int GENERATED ALWAYS AS IDENTITY,
			title VARCHAR(80),
			username VARCHAR(50),
			tag VARCHAR(50),
			created TIMESTAMP,

			CONSTRAINT pk_threads PRIMARY KEY (threadid),
			CONSTRAINT fk_threads_u FOREIGN KEY
				(username) REFERENCES users(username)
				ON DELETE CASCADE
		);	
		CREATE TABLE IF NOT EXISTS posts(
			postid INT GENERATED ALWAYS AS IDENTITY,
			threadid INT,
			title VARCHAR(80),
			username c(50),
			content TEXT,
			created TIMESTAMP,
			likes int,

			CONSTRAINT pk_posts PRIMARY KEY (postid),
			CONSTRAINT fk_posts FOREIGN KEY 
				(threadid) REFERENCES threads(threadid)
				ON DELETE CASCADE,
			CONSTRAINT fk_posts_u FOREIGN KEY
				(username) REFERENCES users(username)
    			ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS comments(
			commentid INT GENERATED ALWAYS AS IDENTITY,
			postid INT,
			username VARCHAR(50),
			content TEXT,
			created TIMESTAMP,

			CONSTRAINT pk_comments PRIMARY KEY (commentid),
			CONSTRAINT fk_comments FOREIGN KEY
				(postid) REFERENCES posts(postid)
    			ON DELETE CASCADE,
			CONSTRAINT fk_comments_u FOREIGN KEY
				(username) REFERENCES users(username)
    			ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS likes(
			username VARCHAR(50),
			postid INT,

			CONSTRAINT fk_likes_u FOREIGN KEY
				(username) REFERENCES users(username)
				ON DELETE CASCADE,
			CONSTRAINT fk_likes_id FOREIGN KEY
				(postid) REFERENCES posts(postid)
				ON DELETE CASCADE
		);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CheckExistingAcc(acc *CreateAccountRequest) error {
	check := (`SELECT * FROM users WHERE email = $1`)
	exist := s.db.QueryRow(check, acc.Email)
	if exists := exist.Scan(); exists != sql.ErrNoRows {
		return fmt.Errorf("an account with this email already exists")
	}

	name := (`SELECT * FROM users WHERE username = $1`)
	dup := s.db.QueryRow(name, acc.UserName)
	if duplicate := dup.Scan(); duplicate != sql.ErrNoRows {
		return fmt.Errorf("an account with this user name already exists")
	}

	return nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := (`
	INSERT INTO users 
	(username, email, encryptedpw, created)
	VALUES 
	($1, $2, $3, $4)
	`)

	_, err := s.db.Query( //Exec and LastInsertId not supported by this psql driver
		query,
		acc.UserName,
		acc.Email,
		acc.EncryptedPW,
		acc.Created)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CreateThread(t *Thread) error {
	query := (`
	INSERT INTO threads 
	(title, username, tag, created)
	VALUES 
	($1, $2, $3, $4) 
	RETURNING threadid
	`)

	row, err := s.db.Query(
		query,
		t.Title,
		t.UserName,
		t.Tag,
		t.Created)

	if err != nil {
		return err
	}

	return retrieveID(row, &t.ThreadID)
}

func (s *PostgresStore) CreatePost(p *Post) error {
	query := (`
	INSERT INTO posts 
	(threadid, title, username, content, created)
	VALUES 
	($1, $2, $3, $4, $5)
	RETURNING postid
	`)

	row, err := s.db.Query(
		query,
		p.ThreadID,
		p.Title,
		p.UserName,
		p.Content,
		p.Created)

	if err != nil {
		return err
	}

	return retrieveID(row, &p.PostID)
}

func (s *PostgresStore) CreateComment(c *Comment) error {
	query := (`
	INSERT INTO comments 
	(postid, username, content, created)
	VALUES 
	($1, $2, $3, $4)
	RETURNING commentid
	`)

	row, err := s.db.Query(
		query,
		c.PostID,
		c.UserName,
		c.Content,
		c.Created)

	if err != nil {
		return err
	}

	return retrieveID(row, &c.CommentID)
}

func (s *PostgresStore) GetLatestThreads(tag string)([]*Thread, error) {
	if tag == "latest" {
		query := (`
		SELECT * FROM threads
		ORDER BY created DESC
		`)
		rows, err := s.db.Query(query, )
		if err != nil {
			return nil, err
		}

		return ScanThreads(rows)

	} else {
		query := (`
		SELECT * FROM threads
		WHERE tag = $1
		ORDER BY created DESC
		`)
		rows, err := s.db.Query(query, tag)
		if err != nil {
			return nil, err
		}
	
		return ScanThreads(rows)
	}
}

func (s *PostgresStore) GetThreadByThreadID(id int) ([]*Thread, error) {
	query := (`
	SELECT * FROM threads
	WHERE threadid = $1
	`)

	row, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	
	return ScanThreads(row)
}

// func (s *PostgresStore) GetPostsByThreadID(id int) ([]*Post, error) {
// 	query := (`
// 	SELECT * FROM posts
// 	WHERE threadid = $1
// 	ORDER BY created DESC
// 	`)

// 	rows, err := s.db.Query(query, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ScanPosts(rows)
// }

func (s *PostgresStore) GetPostsByThreadID(id int, user string) ([]*UserLikedPost, error) {
	query := (`
	SELECT * FROM posts 
		LEFT JOIN likes 
		ON posts.postid = likes.postid 
		WHERE threadid = $1 AND likes.username = $2
	
	`)

	rows, err := s.db.Query(query, id, user)
	if err != nil {
		return nil, err
	}

	return ScanLikedPosts(rows)
}

func (s *PostgresStore) GetPostByPostID(id int) ([]*UserLikedPost, error) {
	query := (`
	SELECT * FROM posts
	WHERE postid = $1
	`)

	row, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	
	return ScanLikedPosts(row)
}

func (s *PostgresStore) GetCommentsByPostID(id int) ([]*Comment, error) {
	query := (`
	SELECT * FROM comments
	WHERE postid = $1
	ORDER BY created DESC
	`)

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	return ScanComments(rows)
}


func (s *PostgresStore) GetThreadPosts(id int, user string) (map[string]interface{}, error) {
	thread, err := s.GetThreadByThreadID(id)
	if err != nil {
		return nil, err
	}

	posts, err := s.GetPostsByThreadID(id, user)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m["parent"] = thread
	m["child"] = posts

	return m, nil
}

func (s *PostgresStore) GetPostComments(id int) (map[string]interface{}, error) {
	post, err := s.GetPostByPostID(id)
	if err != nil {
		return nil, err
	}

	threadid := post[0].ThreadID
	thread, err := s.GetThreadByID(threadid)
	if err != nil {
		return nil, err
	}

	comments, err := s.GetCommentsByPostID(id)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m["root"] = thread
	m["parent"] = post
	m["child"] = comments

	return m, nil
}

func (s *PostgresStore) GetThreadsByUser(userName string) ([]*Thread, error) {
	query := (`
	SELECT * FROM threads
	WHERE username = $1
	ORDER BY created DESC
	`)

	rows, err := s.db.Query(query, userName)
	if err != nil {
		return nil, err
	}

	return ScanThreads(rows)
}

func (s *PostgresStore) GetPostsByUser(userName string) ([]*Post, error) {
	query := (`
	SELECT * FROM posts
	WHERE username = $1
	ORDER BY created DESC
	`)

	rows, err := s.db.Query(query, userName)
	if err != nil {
		return nil, err
	}

	return ScanPosts(rows)
}

func (s *PostgresStore) GetCommentsByUser(userName string) ([]*Comment, error) {
	query := (`
	SELECT * FROM comments
	WHERE username = $1
	ORDER BY created DESC
	`)

	rows, err := s.db.Query(query, userName)
	if err != nil {
		return nil, err
	}

	return ScanComments(rows)
}

func (s *PostgresStore) GetThreadByID(id int) ([]*Thread, error) {
	query := (`
	SELECT * FROM threads
	WHERE threadid = $1
	`)

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	return ScanThreads(rows)
}

func (s *PostgresStore) GetCommentByID(id int) ([]*Comment, error) {
	query := (`
	SELECT * FROM comments
	WHERE commentid = $1
	`)

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	return ScanComments(rows)
}

func (s *PostgresStore) GetUser(typ string, id int) (*string, error) {
	tquery := (`
	SELECT username FROM threads
	WHERE threadid = $1 
	`)
	pquery := (`
	SELECT username FROM posts
	WHERE postid = $1
	`)
	cquery := (`
	SELECT username FROM comments
	WHERE commentid = $1
	`)

	var Q string
	switch {
		case typ == "thread":
			Q = tquery
		case typ == "post":
			Q = pquery
		case typ == "comment":
			Q = cquery
	}

	row, err := s.db.Query(Q, id)
	if err != nil {
		return nil, err
	}

	user := new(string)
	if row.Next() {
		err := row.Scan(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *PostgresStore) GetAccUploads(userName string) (map[string]interface{}, error) {
	threads, err := s.GetThreadsByUser(userName)
	if err != nil {
		return nil, err
	}

	posts, err := s.GetPostsByUser(userName)
	if err != nil {
		return nil, err
	}

	comments, err := s.GetCommentsByUser(userName)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m["threads"] = threads
	m["posts"] = posts
	m["comments"] = comments

	return m, nil
}



func (s *PostgresStore) Delete(typ string, id int) error {
	tquery := (`
	DELETE FROM threads
	WHERE threadid = $1 
	`)
	pquery := (`
	DELETE FROM posts
	WHERE postid = $1
	`)
	cquery := (`
	DELETE FROM comments
	WHERE commentid = $1
	`)

	var Q string
	switch {
		case typ == "thread":
			Q = tquery
		case typ == "post":
			Q = pquery
		case typ == "comment":
			Q = cquery
	}

	_, err := s.db.Query(Q, id)
	if err != nil {
		return err
	}

	fmt.Println("deleted")
	return nil
}

func (s *PostgresStore) Update(input1, input2, typ string, id int) error {
	tquery := (`
	UPDATE threads
	SET title = $1, tag = $2
	WHERE threadid = $3
	`)
	pquery := (`
	UPDATE posts
	SET title = $1, content = $2
	WHERE postid = $3
	`)
	cquery := (`
	UPDATE comments
	SET content = $1
	WHERE commentid = $2
	`)

	var err error
	switch typ{
		case "thread":
			_, err = s.db.Query(tquery, input1, input2, id)
		case "post":
			_, err = s.db.Query(pquery, input1, input2, id)
		case "comment":
			_, err = s.db.Query(cquery, input1, id)
	}

	if err != nil {
		return err
	}

	fmt.Println("updated")
	return nil
}

func (s *PostgresStore) Like(username string, postid int, like bool) error {
	likequery := (`
	SELECT (CASE likes.username WHEN $2 THEN true ELSE false END) as liked FROM likes
		RIGHT JOIN posts 
		ON posts.postid = likes.postid 
		WHERE threadid = $1
	`)

	unlikequery := (`
	DELETE FROM likes WHERE
	username = $1 AND postid = $2;
	`)

	var err error
	if like {
		_,err = s.db.Query(likequery, username, postid)
	} else {
		_,err = s.db.Query(unlikequery, username, postid)
	}
	if err != nil {
		return err
	}

	updatequery := (`
	UPDATE posts
	SET likes = (
		SELECT COUNT(username)
		FROM likes
		WHERE postid = $1)
	`)

	_, uperr := s.db.Query(updatequery, postid)
	if uperr != nil {
		return uperr
	}

	return nil
}

// helpers
func retrieveID (r *sql.Rows, mem *int) error {
	for r.Next() {
		err2 := r.Scan(mem)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func (s *PostgresStore) GetAccountByUserName(userName string) (*Account, error) {
	query := (`
	SELECT * FROM users
	WHERE userName = $1`)
	row, err1 := s.db.Query(query, userName)

	if err1 != nil {
		return nil, err1
	}

	for row.Next() {
		return ScanAccount(row)
	}

	return nil, fmt.Errorf("invalid username")
}

func (s *PostgresStore) GetAccountByEmail(email string) (*Account, error) {
	query := (`
	SELECT * FROM users
	WHERE email = $1`)
	row, err1 := s.db.Query(query, email)

	if err1 != nil {
		return nil, err1
	}

	for row.Next() {
		return ScanAccount(row)
	}

	return nil, fmt.Errorf("invalid email")
}

func ScanAccount(row *sql.Rows) (*Account, error) {
	acc := new(Account)
	err := row.Scan(
		&acc.UserName,
		&acc.Email,
		&acc.EncryptedPW,
		&acc.Created)
	return acc, err
}

func ScanT(row *sql.Rows) (*Thread, error) {
	t := new(Thread)
	err := row.Scan(
		&t.ThreadID,
		&t.Title,
		&t.UserName,
		&t.Tag,
		&t.Created)
	return t, err
}

func ScanThreads(row *sql.Rows) ([]*Thread, error) {
	threads := []*Thread{}
	for row.Next() {
		t, err := ScanT(row)
		if err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}	
	return threads, nil
}

func ScanP(row *sql.Rows) (*Post, error) {
	p := new(Post)
	err := row.Scan(
		&p.PostID,
		&p.ThreadID,
		&p.Title,
		&p.UserName,
		&p.Content,
		&p.Created,)
	return p, err
}

func ScanPosts(row *sql.Rows) ([]*Post, error) {
	posts := []*Post{}
	for row.Next() {
		t, err := ScanP(row)
		if err != nil {
			return nil, err
		}
		posts = append(posts, t)
	}	
	return posts, nil
}

//should try to extract a list posts liked by current user and compare on frontend in map

func ScanLikedP(row *sql.Rows) (*UserLikedPost, error) {
	p := new(UserLikedPost)
	err := row.Scan(
		&p.PostID,
		&p.ThreadID,
		&p.Title,
		&p.UserName,
		&p.Content,
		&p.Created,
		&p.Liked)
	if err != nil {
		return nil, err
	}

	return p, err
}

func ScanLikedPosts(row *sql.Rows) ([]*UserLikedPost, error) {
	posts := []*UserLikedPost{}
	for row.Next() {
		t, err := ScanLikedP(row)
		if err != nil {
			return nil, err
		}
		posts = append(posts, t)
	}	
	return posts, nil
}

func ScanC(row *sql.Rows) (*Comment, error) {
	c := new(Comment)
	err := row.Scan(
		&c.CommentID,
		&c.PostID,
		&c.UserName,
		&c.Content,
		&c.Created)
	return c, err
}

func ScanComments(row *sql.Rows) ([]*Comment, error) {
	comments := []*Comment{}
	for row.Next() {
		c, err := ScanC(row)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}	
	return comments, nil
}

// seeding database
func SeedData(s Database) *Account {

	acc, err := NewAccount("dummyUser2", "dummy2@email.com", "dummyPassword2")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	t, err := NewThread(acc.UserName, "sampleThread2", "sampleTag2")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreateThread(t); err != nil {
		log.Fatal(err)
	}

	p, err := NewPost(t.ThreadID, acc.UserName, "samplePostTitle2", "samplePostContent2")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreatePost(p); err != nil {
		log.Fatal(err)
	}

	c, err := NewComment(p.PostID, acc.UserName, "sampleCommentContent2")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreateComment(c); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("seeded database with: userName[%v], threadID[%v], postID[%v], commentID[%v]\n", acc.UserName, t.ThreadID, p.PostID, c.CommentID)

	return acc
}
