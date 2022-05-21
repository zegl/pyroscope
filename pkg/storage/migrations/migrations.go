package migrations

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dgraph-io/badger/v2"
)

var migrations = []migration{
	migrateDictionaryKeys,
}

// I suppose we might need to expand this next time we do a migration
type Migrations struct {
	main  *badger.DB
	dicts *badger.DB
}

func New(main *badger.DB, dicts *badger.DB) *Migrations {
	return &Migrations{
		main:  main,
		dicts: dicts,
	}
}

type migration func(*Migrations) error

const dbVersionKey = "db-version"

func (s *Migrations) Migrate() error {
	ver, err := s.dbVersion()
	if err != nil {
		return err
	}
	switch {
	case ver == len(migrations):
		return nil
	case ver > len(migrations):
		return fmt.Errorf("db version %d: future versions are not supported", ver)
	}
	for v, m := range migrations[ver:] {
		if err = m(s); err != nil {
			return fmt.Errorf("migration %d: %w", v, err)
		}
	}
	return s.setDbVersion(len(migrations))
}

// dbVersion returns the number of migrations applied to the storage.
func (s *Migrations) dbVersion() (int, error) {
	var version int
	err := s.main.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(dbVersionKey))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			version, err = strconv.Atoi(string(val))
			return err
		})
	})
	if errors.Is(err, badger.ErrKeyNotFound) {
		return 0, nil
	}
	return version, err
}

func (s *Migrations) setDbVersion(v int) error {
	return s.main.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(&badger.Entry{
			Key:   []byte(dbVersionKey),
			Value: []byte(strconv.Itoa(v)),
		})
	})
}
