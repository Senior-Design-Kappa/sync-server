package logic

import (
	"github.com/sync-server/backend"
	"github.com/sync-server/config"
)

type Logic interface {
}
type logic struct {
	backend backend.Backend
	Config  config.Config
}

func NewLogic(conf config.Config, backend backend.Backend) (Logic, error) {
	l := &logic{
		backend: backend,
		Config:  conf,
	}
	return l, nil
}
