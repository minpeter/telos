package database

import (
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var DB *xorm.Engine

func ConnectDatabase() error {
	var err error
	DB, err = xorm.NewEngine("sqlite", "telos.db")

	if err != nil {
		return err
	}

	if err = syncDatabase(); err != nil {
		return err
	}

	return nil
}

func syncDatabase() error {
	err := DB.Sync2(new(User), new(Challenge), new(Solve))
	if err != nil {
		return err
	}

	return nil
}
