package db

import (
	"github.com/eaigner/jet"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

type DB struct {
	d *jet.Db
}

func New(dsn string) (*DB, error) {
	d, err := jet.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = d.Query("SELECT COUNT(1) FROM users").Run(); err != nil {
		return nil, err
	}
	return &DB{d: d}, nil
}
