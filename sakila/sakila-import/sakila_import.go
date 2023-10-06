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

	err = importTable(db, "language", filepath.Join(curDir, "..", "fixtures", "base", "language.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "country", filepath.Join(curDir, "..", "fixtures", "base", "country.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "city", filepath.Join(curDir, "..", "fixtures", "base", "city.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "address", filepath.Join(curDir, "..", "fixtures", "base", "address.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "actor", filepath.Join(curDir, "..", "fixtures", "base", "actor.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "staff", filepath.Join(curDir, "..", "fixtures", "base", "staff.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "store", filepath.Join(curDir, "..", "fixtures", "base", "store.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "category", filepath.Join(curDir, "..", "fixtures", "base", "category.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "inventory", filepath.Join(curDir, "..", "fixtures", "base", "inventory.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "film_actor", filepath.Join(curDir, "..", "fixtures", "base", "film_actor.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "film_category", filepath.Join(curDir, "..", "fixtures", "base", "film_category.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "customer", filepath.Join(curDir, "..", "fixtures", "base", "customer.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	// err = importTable(db, "rental", filepath.Join(curDir, "..", "fixtures", "base", "rental.dbf.yaml"))
	// if err != nil {
	// 	panic(err)
	// }
	//
	// err = importTable(db, "payment", filepath.Join(curDir, "..", "fixtures", "base", "payment.dbf.yaml"))
	// if err != nil {
	// 	panic(err)
	// }

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
