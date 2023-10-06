package main

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/RangelReale/debefix"
	dbsql "github.com/RangelReale/debefix/db/sql"
	"github.com/RangelReale/debefix/db/sql/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	err := importFixtures()
	if err != nil {
		panic(err)
	}
}

func importFixtures() error {
	db, err := sql.Open("pgx",
		fmt.Sprintf("postgres://postgres:password@%s:%s/%s?sslmode=disable", "localhost", "5478", "sakila"))
	if err != nil {
		return err
	}

	curDir, err := currentSourceDirectory()
	if err != nil {
		panic(err)
	}

	data, err := debefix.LoadDirectory(filepath.Join(curDir, "fixtures"))
	if err != nil {
		panic(err)
	}

	// wrap query interface so we can print the output statements
	sqlQueryInterface := dbsql.NewSQLQueryInterface(db)
	wrapQueryInterface := dbsql.QueryInterfaceFunc(func(query string, returnFieldNames []string, args ...any) (map[string]any, error) {
		return sqlQueryInterface.Query(query, returnFieldNames, args...)
	})

	// will send an INSERT SQL for each row to the db, taking table dependency in account for the correct order.
	err = postgres.Resolve(wrapQueryInterface, data)
	if err != nil {
		panic(err)
	}

	return nil
}

func currentSourceDirectory() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filepath.Dir(filename), nil
}
