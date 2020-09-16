package models

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = errors.New("models: email address already in use")
var ErrInvalidCredentials = errors.New("models: invalid user credentials")

type Database struct {
	*sql.DB
}

func (db *Database) GetStory(id int) (*Story, error) {
	query := "SELECT id, title, content, created FROM stories WHERE id = $1"
	row := db.QueryRow(query, id)
	story := &Story{}

	err := row.Scan(&story.ID, &story.Title, &story.Content, &story.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return story, nil
}

func (db *Database) GetLatestStories() (Stories, error) {
	query := "SELECT id, title, content, created FROM stories ORDER BY created DESC LIMIT 10"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stories := Stories{}

	for rows.Next() {
		story := &Story{}

		err := rows.Scan(&story.ID, &story.Title, &story.Content, &story.Created)
		if err != nil {
			return nil, err
		}

		stories = append(stories, story)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stories, nil
}

func (db *Database) InsertStory(title, content string) (int, error) {
	query := "INSERT INTO stories (title, content, created) VALUES ($1, $2, NOW()) RETURNING id"

	var id int

	err := db.QueryRow(query, title, content).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) InsertUser(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, password, created) VALUES($1, $2, $3, NOW())`

	_, err = db.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (db *Database) VerifyUser(email, password string) (int, error) {
	stmt := "SELECT id, password FROM users WHERE email = $1"

	var storedPassword []byte
	var id int

	row := db.QueryRow(stmt, email)
	err := row.Scan(&id, &storedPassword)

	if err == sql.ErrNoRows {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(storedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return id, nil
}
