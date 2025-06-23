package db

import (
	"database/sql"
	"os"
)

const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title   VARCHAR(256) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX scheduler_data_index ON scheduler(date);`

var db *sql.DB

func Init(dbFile string) error {
	var err error

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if !install {
		return nil
	}

	_, err = db.Exec(schema)
	return err
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
