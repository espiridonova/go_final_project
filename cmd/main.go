package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/espiridonova/go_final_project/pkg/api"
	"github.com/espiridonova/go_final_project/pkg/db"
	"github.com/espiridonova/go_final_project/pkg/server"

	_ "modernc.org/sqlite"
)

const defaultDBFile = "scheduler.db"

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}
	if _, err := os.Stat(dbFile); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	err := db.Init(dbFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing database")
		}
	}()

	api.Init()

	err = server.Run()
	if err != nil {
		panic(err)
	}
}
