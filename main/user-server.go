// Package main contains the user-server executable.
package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/polyglottis/rpc"
	"github.com/polyglottis/user_server/database"
	"github.com/polyglottis/user_server/operations"
	"github.com/polyglottis/user_server/server"
)

var dbFile = flag.String("db", "user.db", "path to sqlite db file")
var tcpAddr = flag.String("tcp", ":14773", "TCP address of language server")
var operationsAddr = flag.String("op-tcp", ":17492", "TCP address of operations RPC server")

func main() {
	flag.Parse()

	abs, err := filepath.Abs(*dbFile)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := database.Open(abs)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("User server accessing db file %v", abs)

	main := server.New(db, *tcpAddr)
	op := operations.NewOpServer(db, *operationsAddr)
	p := rpc.NewServerPair("User Server", main, op)

	err = p.RegisterAndListen()
	if err != nil {
		log.Fatalln(err)
	}
	defer p.Close()

	p.Accept()
}
