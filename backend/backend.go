package backend

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sync-server/config"
)

type Backend interface {
}

type backend struct {
	Config config.Config
	db     *sql.DB
}

func NewBackend(conf config.Config) (Backend, error) {
	db, err := sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		fmt.Printf("error: could not open db connection (%+v)\n", err)
	}
	b := &backend{
		Config: conf,
		db:     db,
	}
	return b, nil
}
