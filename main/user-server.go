// Package main contains the user-server executable.
package main

import (
	"log"

	"github.com/polyglottis/platform/config"
	"github.com/polyglottis/rpc"
	"github.com/polyglottis/user_server/database"
	"github.com/polyglottis/user_server/operations"
	"github.com/polyglottis/user_server/server"
)

func main() {
	c := config.Get()

	db, err := database.Open(c.UserDB)
	if err != nil {
		log.Fatalln(err)
	}

	main := server.New(db, c.User)
	op := operations.NewOpServer(db, c.UserOp)
	p := rpc.NewServerPair("User Server", main, op)

	err = p.RegisterAndListen()
	if err != nil {
		log.Fatalln(err)
	}
	defer p.Close()

	p.Accept()
}
