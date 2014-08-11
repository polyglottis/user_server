package operations

import (
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/rpc"
	"github.com/polyglottis/user_server/database"
)

type OpRpcServer struct {
	db *database.DB
}

func NewOpServer(db *database.DB, addr string) *rpc.Server {
	return rpc.NewServer("OpRpcServer", &OpRpcServer{db}, addr)
}

type Dump []*user.Account

func (s *OpRpcServer) Dump(nothing bool, d *Dump) error {
	lines, err := s.db.Dump()
	if err != nil {
		return err
	}
	*d = lines
	return nil
}
