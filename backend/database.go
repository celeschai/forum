package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"os"

	_ "github.com/lib/pq"
)

//to run locally: docker run --name forum_local -e POSTGRES_USER=forumadmin -e POSTGRES_PASSWORD=gossiping -e POSTGRES_DB=forum_containerised -p5432:5432 -d postgres

func NewPostgressStore() (*PostgresStore, error) {
	connStr :=
		"user=" + os.Getenv("DB_USER") +
			" host=" + os.Getenv("DB_HOST") + //remove this line to connect psql container from local machine
			" dbname=" + os.Getenv("DB_NAME") +
			" password=" + os.Getenv("DB_PASSWORD") +
			" sslmode=" + os.Getenv("DB_SSL_MODE") +
			" port=" + os.Getenv("DB_PORT")

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
			username VARCHAR(50),
			content TEXT,
			created TIMESTAMP,

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
		);`

	_, err := s.db.Exec(query)
	return err
}

// new account
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

// new content
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

// feed
func (s *PostgresStore) GetLatestThreads(tag string) ([]*Thread, error) {
	if tag == "latest" {
		query := (`
		SELECT * FROM threads
		ORDER BY created DESC
		`)
		rows, err := s.db.Query(query)
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

func (s *PostgresStore) GetPostByID(id int) ([]*Post, error) {
	query := (`
	SELECT * FROM posts
	WHERE postid = $1
	`)

	row, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	return ScanPosts(row)
}

func (s *PostgresStore) GetPostsByThreadID(id int) ([]*Post, error) {
	query := (`
	SELECT * FROM posts
	WHERE threadid = $1
	ORDER BY created DESC
	`)

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	return ScanPosts(rows)
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

// parentchild
func (s *PostgresStore) GetThreadPosts(id int, user string) (map[string]interface{}, error) {
	thread, err := s.GetThreadByID(id)
	if err != nil {
		return nil, err
	}

	posts, err := s.GetPostsByThreadID(id)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m["parent"] = thread
	m["child"] = posts

	return m, nil
}

func (s *PostgresStore) GetPostComments(id int) (map[string]interface{}, error) {
	post, err := s.GetPostByID(id)
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

// user
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

// accounts
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
	switch typ {
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

// helpers
func retrieveID(r *sql.Rows, mem *int) error {
	for r.Next() {
		err2 := r.Scan(mem)
		if err2 != nil {
			return err2
		}
	}
	return nil
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
		&p.Created)
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
func SeedData(s Database) {
	acc, err := NewAccount("admin", "admin@email.com", "adminPassword")

	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	t, err := NewThread("Best Korean Food", acc.UserName, "University Town")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreateThread(t); err != nil {
		log.Fatal(err)
	}

	tID := strconv.Itoa(t.ThreadID)
	p, err := NewPost(tID, acc.UserName, "Hwangs", "Amazing Jigae!")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreatePost(p); err != nil {
		log.Fatal(err)
	}

	pID := strconv.Itoa(p.PostID)
	c, err := NewComment(pID, acc.UserName, "Yummy!")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreateComment(c); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("seeded database with: userName[%v], threadID[%v], postID[%v], commentID[%v]\n", acc.UserName, t.ThreadID, p.PostID, c.CommentID)
}
