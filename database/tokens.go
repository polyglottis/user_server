// Package database defines the user database.
package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/rand"
)

func (db *DB) NewToken(id user.Name) (string, error) {
	return db.newToken(id, 2*time.Hour)
}

func (db *DB) newToken(id user.Name, timeToLive time.Duration) (string, error) {
	token, err := rand.Id(20)
	if err != nil {
		return "", err
	}

	now := time.Now()
	expiration := now.Add(timeToLive)

	_, err = db.db.Exec("insert into tokens values (?,?,?,?)",
		string(id), token, now.Unix(), expiration.Unix())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (db *DB) ValidToken(id user.Name, token string) (bool, error) {
	var creation, expiration int64
	err := db.db.QueryRow("select creation,expiration from tokens where id=? and token=?",
		string(id), token).Scan(&creation, &expiration)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
	}

	now := time.Now().Unix()
	if creation <= now && now < expiration {
		return true, nil
	}
	return false, nil
}

func (db *DB) DeleteToken(id user.Name, token string) error {
	_, err := db.db.Exec("delete from tokens where id=? and token=?",
		string(id), token)
	return err
}

func (db *DB) DeleteExpiredTokens() error {
	_, err := db.db.Exec("delete from tokens where expiration<?",
		time.Now().Unix())
	return err
}

func (db *DB) deleteTokensPeriodically() {
	go func() {
		for _ = range time.Tick(24 * time.Hour) {
			err := db.DeleteExpiredTokens()
			if err != nil {
				log.Println("Unable to delete expired tokens:", err)
			}
		}
	}()
}
