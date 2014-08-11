// Package user-dump contains the user-dump executable.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/polyglottis/user_server/operations"
)

var operationsAddr = flag.String("op-tcp", ":17492", "TCP address of operations RPC server")

func main() {
	flag.Parse()

	c, err := operations.NewClient(*operationsAddr)
	if err != nil {
		log.Fatalln(err)
	}

	users, err := c.Dump()
	if err != nil {
		log.Fatalln(err)
	}

	for _, u := range users {
		fmt.Println(u)
	}
}
