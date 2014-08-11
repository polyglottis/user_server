package server

import (
	"os"
	"testing"

	userRpc "github.com/polyglottis/platform/user/rpc"
	"github.com/polyglottis/platform/user/test"
	"github.com/polyglottis/rpc"
	"github.com/polyglottis/user_server/database"
	"github.com/polyglottis/user_server/operations"
)

var mainAddr = ":1234"
var opAddr = ":2345"
var testDB = "user_test.db"

func TestServer(t *testing.T) {

	db, err := database.Open(testDB)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(testDB)

	main := New(db, mainAddr)
	op := operations.NewOpServer(db, opAddr)
	p := rpc.NewServerPair("User Test Server", main, op)

	err = p.RegisterAndListen()
	if err != nil {
		t.Fatal(err)
	}

	go p.Accept()

	opc, err := operations.NewClient(opAddr)
	if err != nil {
		t.Fatal(err)
	}

	users, err := opc.Dump()
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 0 {
		t.Fatal("Database should be empty")
	}

	c, err := userRpc.NewClient(mainAddr)
	if err != nil {
		t.Fatal(err)
	}

	a, err := c.NewAccount(test.Account)
	if err != nil {
		t.Fatal(err)
	}

	b, err := c.GetAccount(test.Account.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !a.Equals(b) {
		t.Fatal("Accounts should coincide, but they don't: %+v != %+v", a, b)
	}

	users, err = opc.Dump()
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 1 {
		t.Fatal("DB should contain one account at this point")
	}

	if !a.Equals(users[0]) {
		t.Fatal("Accounts should coincide, but they don't: %+v != %+v", a, b)
	}
}
