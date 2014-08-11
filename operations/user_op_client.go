// Package operations contains an rpc client-server pair for maintenance operations on the user server.
package operations

import (
	"net/rpc"

	"github.com/polyglottis/platform/user"
)

type Client struct {
	c *rpc.Client
}

// NewClient creates an rpc client for maintenance operations on the user server.
func NewClient(addr string) (*Client, error) {
	c, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{c: c}, nil
}

// Dump returns a dump of the whole language database.
func (c *Client) Dump() ([]*user.Account, error) {
	var d Dump
	err := c.c.Call("OpRpcServer.Dump", false, &d)
	if err != nil {
		return nil, err
	}
	return []*user.Account(d), nil
}
