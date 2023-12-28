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
			username varchar(50),
			email varchar(50),
			encryptedpw varchar(100),
			created timestamp,

			CONSTRAINT pk_account PRIMARY KEY (username)
		);
		CREATE TABLE IF NOT EXISTS threads(
			threadid int GENERATED ALWAYS AS IDENTITY,
			title varchar(80),
			username varchar(50),
			tag varchar(50),
			created timestamp,

			CONSTRAINT pk_threads PRIMARY KEY (threadid),
			CONSTRAINT fk_threads_u FOREIGN KEY
				(username) REFERENCES users(username)
				ON DELETE CASCADE
		);	
		CREATE TABLE IF NOT EXISTS posts(
			postid INT GENERATED ALWAYS AS IDENTITY,
			threadid INT,
			title VARCHAR(80),
			username varchar(50),
			content TEXT,
			created timestamp,

			CONSTRAINT pk_posts PRIMARY KEY (postid),
			CONSTRAINT fk_posts FOREIGN KEY 
				(threadid) REFERENCES threads(threadid),
			CONSTRAINT fk_posts_u FOREIGN KEY
				(username) REFERENCES users(username)
    			ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS comments(
			commentid INT GENERATED ALWAYS AS IDENTITY,
			postid INT,
			username varchar(50),
			content TEXT,
			created timestamp,

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

func (s *PostgresStore) CheckExistingAcc(acc *Account) error {
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

func (s *PostgresStore) GetLatestThreads(pg int)([]*Thread, error) {
	query := (`
	SELECT * FROM threads
	ORDER BY created DESC
	LIMIT 10
	OFFSET (10*$1)
	`)

	rows, err := s.db.Query(query, pg)
	if err != nil {
		return nil, err
	}

	threads := []*Thread{}
	for rows.Next() {
		t, err := ScanThread(rows)
		if err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, nil
}

// helpers
func retrieveID(r *sql.Rows, mem any) error {
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

func ScanThread(row *sql.Rows) (*Thread, error) {
	t := new(Thread)
	err := row.Scan(
		&t.ThreadID,
		&t.Title,
		&t.UserName,
		&t.Tag,
		&t.Created)
	return t, err
}

// seeding database
func SeedData(s Database) *Account {

	acc, err := NewAccount("dummyUser", "dummy@email.com", "dummyPassword")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	t, err := NewThread(acc.UserName, "sampleThread", "sampleTag")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreateThread(t); err != nil {
		log.Fatal(err)
	}

	p, err := NewPost(t.ThreadID, acc.UserName, "samplePostTitle", "samplePostContent")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreatePost(p); err != nil {
		log.Fatal(err)
	}

	c, err := NewComment(p.PostID, acc.UserName, "sampleCommentContent")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreateComment(c); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("seeded database with: userName[%v], threadID[%v], postID[%v], commentID[%v]\n", acc.UserName, t.ThreadID, p.PostID, c.CommentID)

	return acc
}
