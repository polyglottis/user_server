// Package user-dump contains the user-dump executable.
package main

import (
	"fmt"
	"log"

	"github.com/polyglottis/platform/config"
	"github.com/polyglottis/user_server/operations"
)

func main() {
	conf := config.Get()

	c, err := operations.NewClient(conf.UserOp)
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
