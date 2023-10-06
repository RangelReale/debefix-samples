package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/goccy/go-yaml"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	err := importData()
	if err != nil {
		panic(err)
	}
}

func importData() error {
	db, err := sql.Open("pgx",
		fmt.Sprintf("postgres://postgres:password@%s:%s/%s?sslmode=disable", "localhost", "5478", "sakila"))
	if err != nil {
		return err
	}

	curDir, err := currentSourceDirectory()
	if err != nil {
		panic(err)
	}

	err = importTable(db, "country", filepath.Join(curDir, "..", "fixtures", "base", "country.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	return nil
}

func importTable(db *sql.DB, tableName string, outputFilename string) error {
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM "%s"`, tableName))
	if err != nil {
		return err
	}
	defer rows.Close()

	table := &Table{}

	for rows.Next() {
		row, err := rowToMap(rows)
		if err != nil {
			return err
		}
		table.Rows = append(table.Rows, row)
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	data := Data{
		tableName: table,
	}

	f, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	err = enc.Encode(data)
	if err != nil {
		return err
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
