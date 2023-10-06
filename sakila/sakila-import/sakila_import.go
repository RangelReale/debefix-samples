package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/goccy/go-yaml"
	"github.com/iancoleman/strcase"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	err := importData()
	if err != nil {
		panic(err)
	}
}

type specialData struct {
	Tablename string
	Fieldname string
	RefID     map[any]string
}

func importData() error {
	db, err := sql.Open("pgx",
		fmt.Sprintf("postgres://postgres:password@%s:%s/%s?sslmode=disable", "localhost", "5478", "sakila"))
	if err != nil {
		return err
	}

	sdLanguage, err := getSpecialData(db, "language", "language_id", "name")
	if err != nil {
		panic(err)
	}

	sdCategory, err := getSpecialData(db, "category", "category_id", "name")
	if err != nil {
		panic(err)
	}

	sdCountry, err := getSpecialData(db, "country", "country_id", "country")
	if err != nil {
		panic(err)
	}

	sdata := []*specialData{sdLanguage, sdCategory, sdCountry}

	curDir, err := currentSourceDirectory()
	if err != nil {
		panic(err)
	}

	err = importTable(db, "language", sdata, filepath.Join(curDir, "..", "fixtures", "base", "language.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "country", sdata, filepath.Join(curDir, "..", "fixtures", "base", "country.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "city", sdata, filepath.Join(curDir, "..", "fixtures", "base", "city.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "address", sdata, filepath.Join(curDir, "..", "fixtures", "base", "address.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "actor", sdata, filepath.Join(curDir, "..", "fixtures", "base", "actor.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "staff", sdata, filepath.Join(curDir, "..", "fixtures", "base", "staff.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "store", sdata, filepath.Join(curDir, "..", "fixtures", "base", "store.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "category", sdata, filepath.Join(curDir, "..", "fixtures", "base", "category.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "film", sdata, filepath.Join(curDir, "..", "fixtures", "base", "film.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "inventory", sdata, filepath.Join(curDir, "..", "fixtures", "base", "inventory.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "film_actor", sdata, filepath.Join(curDir, "..", "fixtures", "base", "film_actor.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "film_category", sdata, filepath.Join(curDir, "..", "fixtures", "base", "film_category.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "customer", sdata, filepath.Join(curDir, "..", "fixtures", "base", "customer.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	// err = importTable(db, "rental", sdata, filepath.Join(curDir, "..", "fixtures", "base", "rental.dbf.yaml"))
	// if err != nil {
	// 	panic(err)
	// }
	//
	// err = importTable(db, "payment", sdata, filepath.Join(curDir, "..", "fixtures", "base", "payment.dbf.yaml"))
	// if err != nil {
	// 	panic(err)
	// }

	return nil
}

func getSpecialData(db *sql.DB, tableName string, fieldName string, textFieldName string) (*specialData, error) {
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM "%s"`, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sdata := &specialData{
		Tablename: tableName,
		Fieldname: fieldName,
		RefID:     map[any]string{},
	}

	for rows.Next() {
		row, err := rowToMap(rows)
		if err != nil {
			return nil, err
		}

		sdata.RefID[row[fieldName]] = strcase.ToSnake(row[textFieldName].(string))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return sdata, nil
}

func importTable(db *sql.DB, tableName string, sdata []*specialData, outputFilename string) error {
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM "%s"`, tableName))
	if err != nil {
		return err
	}
	defer rows.Close()

	var currentSpecialData *specialData
	for _, s := range sdata {
		if s.Tablename == tableName {
			currentSpecialData = s
			break
		}
	}

	table := &Table{}

	for rows.Next() {
		row, err := rowToMap(rows)
		if err != nil {
			return err
		}

		if currentSpecialData != nil {
			row["_dbfconfig"] = RowConfig{
				RefID: currentSpecialData.RefID[row[currentSpecialData.Fieldname]],
			}
		} else {
			for _, s := range sdata {
				if sfd, ok := row[s.Fieldname]; ok {
					row[s.Fieldname] = TaggedString{
						Tag:   "!dbfexpr",
						Value: fmt.Sprintf("refid:%s:%s:%s", s.Tablename, s.RefID[sfd], s.Fieldname),
					}
				}
			}
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
