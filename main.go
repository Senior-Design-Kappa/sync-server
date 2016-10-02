package main

import (
	"log"

	"github.com/sync-server/api"
	"github.com/sync-server/backend"
	"github.com/sync-server/config"
	"github.com/sync-server/logic"
)

func main() {
	conf := config.NewConfig()
	b := makeBackend(conf)
	l := makeLogic(conf, b)
	s := api.NewServer(conf, l)
	log.Fatal(s.ListenAndServe())
}

func makeBackend(conf config.Config) backend.Backend {
	b, err := backend.NewBackend(conf)
	if err != nil {
		log.Fatalf("error: backend layer could not be created (%+v)\n", err)
	}
	return b
}

func makeLogic(conf config.Config, backend backend.Backend) logic.Logic {
	l, err := logic.NewLogic(conf, backend)
	if err != nil {
		log.Fatalf("error: logic layer could not be created (%+v)\n", err)
	}
	return l
}
