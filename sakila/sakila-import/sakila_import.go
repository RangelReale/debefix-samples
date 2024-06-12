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
		fmt.Sprintf("postgres://postgres:password@%s:%s/%s?sslmode=disable", "localhost", "5479", "sakila"))
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

	err = importTable(db, "language", nil, sdata, filepath.Join(curDir, "..", "fixtures", "base", "language.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "country", nil, sdata, filepath.Join(curDir, "..", "fixtures", "base", "country.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "city", nil, sdata, filepath.Join(curDir, "..", "fixtures", "base", "city.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "address", []string{"city"}, sdata, filepath.Join(curDir, "..", "fixtures", "base", "address.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "actor", nil, sdata, filepath.Join(curDir, "..", "fixtures", "base", "actor.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "category", nil, sdata, filepath.Join(curDir, "..", "fixtures", "base", "category.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	err = importTable(db, "film", nil, sdata, filepath.Join(curDir, "..", "fixtures", "base", "film.dbf.yaml"))
	if err != nil {
		panic(err)
	}

	// store and staff depends on each other, more complex logic would be needed to export in the correct order

	// err = importTable(db, "customer", []string{"store", "address"}, sdata, filepath.Join(curDir, "..", "fixtures", "base", "customer.dbf.yaml"))
	// if err != nil {
	// 	panic(err)
	// }

	// err = importTable(db, "staff", []string{"address", "store"}, sdata, filepath.Join(curDir, "..", "fixtures", "base", "staff.dbf.yaml"))
	// if err != nil {
	// 	panic(err)
	// }

	// err = importTable(db, "store", []string{"address", "staff"}, sdata, filepath.Join(curDir, "..", "fixtures", "base", "store.dbf.yaml"))
	// if err != nil {
	// 	panic(err)
	// }

	// err = importTable(db, "inventory", []string{"film", "store"}, sdata, filepath.Join(curDir, "..", "fixtures", "base", "inventory.dbf.yaml"))
	// if err != nil {
	// 	panic(err)
	// }

	// err = importTable(db, "rental", []string{"inventory", "customer", "staff"}, sdata, filepath.Join(curDir, "..", "fixtures", "base", "rental.dbf.yaml"))
	// if err != nil {
	// 	panic(err)
	// }
	//
	// err = importTable(db, "payment", []string{"customer", "staff", "rental"}, sdata, filepath.Join(curDir, "..", "fixtures", "base", "payment.dbf.yaml"))
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

func importTable(db *sql.DB, tableName string, deps []string, sdata []*specialData, outputFilename string) error {
	data, err := importTableData(db, tableName, sdata)
	if err != nil {
		return err
	}

	data.Tables[tableName].Config.Depends = deps

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

func importTableData(db *sql.DB, tableName string, sdata []*specialData) (*Data, error) {
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM "%s"`, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return importTableRows(db, rows, tableName, sdata, nil, nil)
}

func importTableDataFilmCategory(db *sql.DB, filmID any, sdata []*specialData, rowTags []string) (*Data, error) {
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM "%s" WHERE "film_id" = $1`, "film_category"), filmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return importTableRows(db, rows, "film_category", sdata, rowTags, func(row Row) error {
		row["film_id"] = TaggedString{
			Tag:   "!dbfexpr",
			Value: "parent:film_id",
		}
		return nil
	})
}

func importTableDataFilmActor(db *sql.DB, filmID any, sdata []*specialData, rowTags []string) (*Data, error) {
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM "%s" WHERE "film_id" = $1`, "film_actor"), filmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return importTableRows(db, rows, "film_actor", sdata, rowTags, func(row Row) error {
		row["film_id"] = TaggedString{
			Tag:   "!dbfexpr",
			Value: "parent:film_id",
		}
		return nil
	})
}

func importTableRows(db *sql.DB, rows *sql.Rows, tableName string, sdata []*specialData,
	rowTags []string, customize func(row Row) error) (*Data, error) {
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
			return nil, err
		}

		var rowConfig *RowConfig
		if rowTags != nil {
			rowConfig = &RowConfig{
				Tags: rowTags,
			}
		}

		if currentSpecialData != nil {
			if rowConfig == nil {
				rowConfig = &RowConfig{}
			}
			rowConfig.RefID = currentSpecialData.RefID[row[currentSpecialData.Fieldname]]
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

		deps := Data{
			Tables: map[string]*Table{},
		}

		if tableName == "film" {
			filmRowTags := []string{"rate_" + strcase.ToSnake(row["rating"].(string))}

			filmCategoryData, err := importTableDataFilmCategory(db, row["film_id"], sdata, filmRowTags)
			if err != nil {
				return nil, err
			}
			deps.Tables["film_category"] = filmCategoryData.Tables["film_category"]

			filmActorData, err := importTableDataFilmActor(db, row["film_id"], sdata, filmRowTags)
			if err != nil {
				return nil, err
			}
			deps.Tables["film_actor"] = filmActorData.Tables["film_actor"]

			// add tags for rating
			if rowConfig == nil {
				rowConfig = &RowConfig{}
			}

			rowConfig.Tags = append(rowConfig.Tags, filmRowTags...)
		}

		if len(deps.Tables) > 0 {
			row["deps"] = &TaggedValue{
				Tag:   "!dbfdeps",
				Value: deps,
			}
		}

		if rowConfig != nil {
			row["config"] = &TaggedValue{
				Tag:   "!dbfconfig",
				Value: *rowConfig,
			}
		}

		if customize != nil {
			err = customize(row)
			if err != nil {
				return nil, err
			}
		}

		table.Rows = append(table.Rows, row)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &Data{
		Tables: map[string]*Table{
			tableName: table,
		},
	}, nil
}

func currentSourceDirectory() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filepath.Dir(filename), nil
}
