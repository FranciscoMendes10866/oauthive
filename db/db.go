package db

import (
	"sync"

	"github.com/go-rel/rel"
	"github.com/go-rel/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseAdapter interface {
	GetClient() rel.Repository
	Close() error
}

type databaseAdapter struct {
	repo      rel.Repository
	adapter   rel.Adapter
	once      sync.Once
	closeOnce sync.Once
	initErr   error
}

var (
	instance *databaseAdapter
	once     sync.Once
)

func New(dsn string) (DatabaseAdapter, error) {
	var err error
	once.Do(func() {
		instance = &databaseAdapter{}
		instance.adapter, err = sqlite3.Open(dsn)
		if err != nil {
			instance.initErr = err
			return
		}
		instance.repo = rel.New(instance.adapter)
	})
	if instance.initErr != nil {
		return nil, instance.initErr
	}
	return instance, nil
}

func (d *databaseAdapter) GetClient() rel.Repository {
	return d.repo
}

func (d *databaseAdapter) Close() error {
	var err error
	d.closeOnce.Do(func() {
		if d.adapter != nil {
			err = d.adapter.Close()
		}
	})
	return err
}
