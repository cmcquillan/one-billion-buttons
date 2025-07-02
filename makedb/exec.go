package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"slices"
	"strings"
	"time"
)

func ExecFile(dbc *sql.DB, file string) error {
	fileData, err := os.ReadFile(file)

	if len(fileData) == 0 {
		return errors.New("no sql script is available to run")
	}

	if err != nil {
		return err
	}

	sqlScript := string(fileData)

	tx, err := dbc.Begin()

	if err != nil {
		return err
	}

	_, errSql := tx.Exec(sqlScript)

	if errSql != nil {
		log.Print(errSql)
		return errSql
	}

	errSql = tx.Commit()

	if errSql != nil {
		log.Print(errSql)
		return errSql
	}

	return nil
}

func ExecDir(dbc *sql.DB, dir string) error {
	mFiles, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	log.Printf("found %d scripts to execute", len(mFiles))

	mIter := slices.Values(mFiles)
	mOrdered := slices.SortedFunc(mIter, func(a os.DirEntry, b os.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	for _, f := range mOrdered {
		start := time.Now()
		log.Printf("executing %v...", f.Name())
		errFile := ExecFile(dbc, dir+"/"+f.Name())
		dur := -time.Until(start)
		log.Printf("executed %v... %v", f.Name(), dur)

		if errFile != nil {
			return errFile
		}
	}

	return nil
}
