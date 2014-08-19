package database

import (
	"os"
	"testing"
	"time"

	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/platform/user/test"
)

var testDB = "user_test.db"

func TestNewAccount(t *testing.T) {
	os.Remove(testDB)

	db, err := Open(testDB)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(testDB)

	if db == nil {
		t.Fatal("Open should never return a nil db")
	}

	tester := test.NewTester(db, t)
	tester.All()
}

func TestDeleteExpiredTokens(t *testing.T) {
	os.Remove(testDB)

	db, err := Open(testDB)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(testDB)

	n := user.Name("AAA")
	token, err := db.newToken(n, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := db.ValidToken(n, token)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("Token should be valid")
	}

	time.Sleep(1*time.Second + 500*time.Millisecond)

	valid, err = db.ValidToken(n, token)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("Token should have expired")
	}

	err = db.DeleteExpiredTokens()
	if err != nil {
		t.Fatal(err)
	}

	num, err := db.db.QueryInt("select count(1) from tokens")
	if err != nil {
		t.Fatal(err)
	}
	if num != 0 {
		t.Fatal("All tokens should have expired and have been deleted!")
	}
}
