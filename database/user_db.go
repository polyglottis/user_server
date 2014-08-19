// Package database defines the user database.
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // driver import

	"github.com/polyglottis/platform/database"
	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/user"
)

type DB struct {
	db *database.DB
}

func Open(file string) (*DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	userDB, err := database.Create(db, database.Schema{{
		Name: "users",
		Columns: database.Columns{{
			Field:      "id",
			Type:       "text",
			Constraint: "primary key not null",
		}, {
			Field: "mainlanguage",
			Type:  "text",
		}, {
			Field: "active",
			Type:  "boolean",
		}, {
			Field: "email",
			Type:  "text",
		}, {
			Field: "pwhash",
			Type:  "blob",
		}},
	}, {
		Name: "tokens",
		Columns: database.Columns{{
			Field:      "id",
			Type:       "text",
			Constraint: "not null",
		}, {
			Field: "token",
			Type:  "text",
		}, {
			Field: "creation",
			Type:  "integer",
		}, {
			Field: "expiration",
			Type:  "integer",
		}},
		PrimaryKey: []string{"id", "token"},
	}})
	if err != nil {
		return nil, err
	}

	newDB := &DB{
		db: userDB,
	}
	newDB.deleteTokensPeriodically()
	return newDB, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) NewAccount(r *user.NewAccountRequest) (*user.Account, error) {
	if r == nil {
		return nil, fmt.Errorf("NewAccountRequest should not be nil")
	}
	if valid, _ := user.ValidName(string(r.Name)); !valid {
		return nil, fmt.Errorf("Invalid username.")
	}

	_, err := db.db.Exec("insert into users values (?,?,?,?,?)",
		string(r.Name), string(r.MainLanguage), true, r.Email, r.PasswordHash)
	if err != nil {
		return nil, err
	}

	return user.NewAccount(r), nil
}

func (db *DB) GetAccount(name user.Name) (*user.Account, error) {
	a, err := db.scanAccount(db.db.QueryRow("select * from users where id=?", string(name)))
	if err == sql.ErrNoRows {
		a, err = db.GetAccountByEmail(string(name))
	}
	switch {
	case err == sql.ErrNoRows:
		return nil, user.AccountNotFound
	case err != nil:
		return nil, err
	default:
		return a, nil
	}
}

func (db *DB) GetAccountByEmail(email string) (*user.Account, error) {
	a, err := db.scanAccount(db.db.QueryRow("select * from users where email=?", email))
	switch {
	case err == sql.ErrNoRows:
		return nil, user.AccountNotFound
	case err != nil:
		return nil, err
	default:
		return a, nil
	}
}

func (db *DB) UpdateAccount(a *user.Account) error {
	_, err := db.db.Exec("update users set mainlanguage=?, active=?, email=?, pwhash=? where id=?",
		string(a.MainLanguage), a.Active, a.Email, a.PasswordHash, string(a.Name))
	return err
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

func (db *DB) scanAccount(s Scanner) (*user.Account, error) {
	a := new(user.Account)
	var uName, lang string
	err := s.Scan(&uName, &lang, &a.Active, &a.Email, &a.PasswordHash)
	if err != nil {
		return nil, err
	}
	a.Name = user.Name(uName)
	a.MainLanguage = language.Code(lang)
	return a, nil
}

func (db *DB) Dump() ([]*user.Account, error) {
	dump := make([]*user.Account, 0)
	rows, err := db.db.Query("select * from users")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		a, err := db.scanAccount(rows)
		if err != nil {
			return nil, err
		}
		dump = append(dump, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dump, nil
}
