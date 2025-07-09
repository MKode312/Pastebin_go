package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"text_sharing/internal/lib/jwt"
	"text_sharing/internal/storage"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	db *sql.DB
}

func New(storageParh string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storageParh)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt1, err := db.Prepare(`
	    CREATE TABLE IF NOT EXISTS users(
	    id INTEGER PRIMARY KEY,
	    username TEXT NOT NULL UNIQUE,
	    password TEXT NOT NULL);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt2, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS textblocks(
		id INTEGER PRIMARY KEY,
		date TEXT,
		link TEXT NOT NULL UNIQUE,
		linkID TEXT NOT NULL UNIQUE);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt1.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt2.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(username string, password string) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	hashedPswrd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(username, hashedPswrd)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) LoginUser(username string, password string) (string, error) {
	const op = "storage.sqlite.LoginUser"

	stmt, err := s.db.Prepare("SELECT password, id FROM users WHERE username = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resPswrd string
	var resID int64

	err = stmt.QueryRow(username).Scan(&resPswrd, &resID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUserNotFound
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(resPswrd), []byte(password))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.GenerateJWTToken(resID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *Storage) SaveLink(urlToSave string, date string, objectID string) (string, error) {
	const op = "storage.sqlite.SaveLink"

	stmt, err := s.db.Prepare("INSERT INTO textblocks(link, date, linkID) VALUES(?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(urlToSave, date, objectID)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return "", fmt.Errorf("%s: %w", op, storage.ErrLinkExists)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return objectID, nil
}

func (s *Storage) GetLink(linkID string) (string, error) {
	const op = "storage.sqlite.GetLink"

	stmt, err := s.db.Prepare("SELECT link FROM textblocks WHERE linkID = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resLink string

	err = stmt.QueryRow(linkID).Scan(&resLink)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrLinkNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resLink, nil
}

func (s *Storage) GetDate(linkID string) (string, error) {
	const op = "storage.sqlite.GetLink"

	stmt, err := s.db.Prepare("SELECT date FROM textblocks WHERE linkID = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resDate string

	err = stmt.QueryRow(linkID).Scan(&resDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrLinkNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resDate, nil
}


