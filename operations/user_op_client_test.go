package operations

import (
	"os"
	"testing"

	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/user_server/database"
)

var file = "user_test.db"
var testAddr = ":1234"
var testAccount = &user.NewAccountRequest{
	Name:         "testUser",
	MainLanguage: language.English.Code,
	Email:        "test@test.com",
	PasswordHash: []byte("testPW"),
}

func TestClientOperationServer(t *testing.T) {

	os.Remove(file)
	db, err := database.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(file)

	op := NewOpServer(db, testAddr)
	if err != nil {
		t.Fatal(err)
	}

	err = op.RegisterAndListen()
	if err != nil {
		t.Fatal(err)
	}

	go op.Accept()

	c, err := NewClient(testAddr)
	if err != nil {
		t.Fatal(err)
	}

	d, err := c.Dump()
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 0 {
		t.Error("Database should be empty at this point")
	}

	a, err := db.NewAccount(testAccount)
	if err != nil {
		t.Fatal(err)
	}

	d, err = c.Dump()
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 1 {
		t.Fatal("Database should contain exactly one line at this point")
	}
	if !a.Equals(d[0]) {
		t.Errorf("Incorrect database dump:\n%+v\n\n%+v", a, d[0])
	}
}
