package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// docker run --name forum -e POSTGRES_USER=forumadmin -e POSTGRES_PASSWORD=gossiping -e POSTGRES_DB=forum -p 5432:5432 -d postgres
// test: docker run --name test -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -e POSTGRES_DB=test -p 5434:5433 -d postgres

func NewPostgressStore() (*PostgresStore, error) {
	connStr := "user=forumadmin dbname=forum password=gossiping sslmode=disable" //change pw in command then here
	// connStr := "port=5433 user=test dbname=test password=test sslmode=disable"
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
			userid int GENERATED ALWAYS AS IDENTITY, 
			username varchar(50),
			email varchar(50),
			encryptedpw varchar(100),
			created timestamp,

			CONSTRAINT pk_account PRIMARY KEY (userid)
		);
		CREATE TABLE IF NOT EXISTS threads(
			threadid int GENERATED ALWAYS AS IDENTITY,
			title varchar(80),
			userid int,
			tag1 varchar(50),
			tag2 varchar(50),
			created timestamp,

			CONSTRAINT pk_threads PRIMARY KEY (threadid)
		);	
		CREATE TABLE IF NOT EXISTS posts(
			postid INT GENERATED ALWAYS AS IDENTITY,
			threadid INT,
			title VARCHAR(80),
			userid INT,
			content TEXT,
			created timestamp,

			CONSTRAINT pk_posts PRIMARY KEY (postid),
			CONSTRAINT fk_posts FOREIGN KEY 
				(threadid) REFERENCES thread(threadid)
		);
		CREATE TABLE IF NOT EXISTS comments(
			commentid INT GENERATED ALWAYS AS IDENTITY,
			postid INT,
			userid INT,
			content TEXT,
			created timestamp,

			CONSTRAINT pk_comments PRIMARY KEY (commentid),
			CONSTRAINT fk_comments FOREIGN KEY
				(postid) REFERENCES posts(postid)
    			ON DELETE CASCADE
		);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	check := (`
	SELECT * FROM users
	WHERE username = $1 OR email = $2`)
	
	exist := s.db.QueryRow(check, acc.UserName, acc.Email)
	if exists := exist.Scan(); exists != sql.ErrNoRows {
		return fmt.Errorf("username or email already exists, proceed to login in")
	}
	
	query := (`
	INSERT INTO users 
	(username, email, encryptedpw, created)
	VALUES 
	($1, $2, $3, $4)
	RETURNING userid
	`)

	row, err := s.db.Query( //Exec and LastInsertId not supported by this psql driver
		query,
		acc.UserName,
		acc.Email,
		acc.EncryptedPW,
		acc.Created)
	if err != nil {
		return err
	}
	
	return retrieveID(row, &acc.UserID)
}

func (s *PostgresStore) CreateThread(t *Thread) error {
	query := (`
	INSERT INTO threads 
	(title, userid, tag1, tag2, created)
	VALUES 
	($1, $2, $3, $4, $5) 
	RETURNING threadid
	`)

	row, err := s.db.Query(
		query,
		t.Title,
		t.UserID,
		t.Tag1,
		t.Tag2,
		t.Created)

	if err != nil {
		return err
	}

	return retrieveID(row, &t.ThreadID)
}

func (s *PostgresStore) CreatePost(p *Post) error {
	query := (`
	INSERT INTO posts 
	(threadid, title, userid, content, created)
	VALUES 
	($1, $2, $3, $4, $5)
	RETURNING postid
	`)

	row, err := s.db.Query(
		query,
		p.ThreadID,
		p.Title,
		p.UserID,
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
	(postid, userid, content, created)
	VALUES 
	($1, $2, $3, $4)
	RETURNING commentid
	`)

	row, err := s.db.Query(
		query,
		c.PostID,
		c.UserID,
		c.Content,
		c.Created)

	if err != nil {
		return err
	}

	return retrieveID(row, &c.CommentID)
}

func retrieveID(r *sql.Rows, mem any) error {
	for(r.Next()) {
		err2 := r.Scan(mem)
		if err2 != nil {
			return err2
		}
	}
	return nil
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
	

	t, err := NewThread(acc.UserID, "sampleThread", "sampleTag1", "sampleTag2")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreateThread(t); err != nil {
		log.Fatal(err)
	}

	p, err := NewPost(t.ThreadID, acc.UserID, "samplePostTitle", "samplePostContent")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreatePost(p); err != nil {
		log.Fatal(err)
	}

	c, err := NewComment(p.PostID, acc.UserID, "sampleCommentContent")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.CreateComment(c); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("seeded database with: userID[%v], threadID[%v], postID[%v], commentID[%v]\n", acc.UserID, t.ThreadID, p.PostID, c.CommentID)

	return acc
}

func (s *PostgresStore) GetAccountByUserID (userID int) (*Account, error) {
	query := (`
	SELECT * FROM users
	WHERE userid = $1`)
	row, err1 := s.db.Query(query, userID)

	if err1 != nil {
		return nil, err1
	}

	acc := new(Account)
	err2 := row.Scan(
		&acc.UserID,
		&acc.UserName,
		&acc.Email,
		&acc.EncryptedPW,
		&acc.Created)

	if err2 != nil {
		return nil, err2
	}

	return acc, nil
}