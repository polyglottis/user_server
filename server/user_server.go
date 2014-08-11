// Package server defines the Polyglottis User Server.
package server

import (
	userRpc "github.com/polyglottis/platform/user/rpc"
	"github.com/polyglottis/rpc"
	"github.com/polyglottis/user_server/database"
)

// New creates the rpc user server, as required by polyglottis/user/rpc
func New(db *database.DB, addr string) *rpc.Server {
	return userRpc.NewUserServer(db, addr)
}
