package db

import (
	"log"

	"github.com/go-rel/rel"
	"github.com/go-rel/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

func Init(dsn string) rel.Repository {
	adapter, err := sqlite3.Open(dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer adapter.Close()

	return rel.New(adapter)
}
