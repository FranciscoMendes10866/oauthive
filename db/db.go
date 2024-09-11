package db

import (
	"sync"

	"github.com/go-rel/rel"
	"github.com/go-rel/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

var (
	repo  rel.Repository
	once  sync.Once
	dbErr error
)

func Init(dsn string) (rel.Repository, error) {
	once.Do(func() {
		adapter, err := sqlite3.Open(dsn)
		if err != nil {
			dbErr = err
			return
		}

		repo = rel.New(adapter)
	})

	return repo, dbErr
}
