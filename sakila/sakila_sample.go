package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	dbsql "github.com/rrgmc/debefix-db/v2/sql"
	"github.com/rrgmc/debefix-db/v2/sql/postgres"
	"github.com/rrgmc/debefix/v2"
)

var useDB = flag.Bool("use-db", false, "use db")

func main() {
	flag.Parse()

	ctx := context.Background()

	err := importFixtures(ctx)
	if err != nil {
		panic(err)
	}
}

func importFixtures(ctx context.Context) error {
	var db *sql.DB
	var err error

	if *useDB {
		db, err = sql.Open("pgx",
			fmt.Sprintf("postgres://postgres:password@%s:%s/%s?sslmode=disable", "localhost", "5478", "sakila"))
		if err != nil {
			return err
		}
	} else {
		fmt.Println("not using DB, dumping only")
	}

	// wrap query interface, so we can print the output statements
	var sqlQueryInterface dbsql.QueryInterface

	if *useDB {
		// will send an INSERT SQL for each row to the db, taking table dependency in account for the correct order.
		sqlQueryInterface = dbsql.NewSQLQueryInterface(db)
	}

	insertCount := 0

	data, err := Data()
	if err != nil {
		return err
	}

	debugQI := dbsql.NewDebugQueryInterface(nil)

	_, err = debefix.Resolve(ctx, data,
		postgres.ResolveFunc(dbsql.QueryInterfaceFunc(func(ctx context.Context, tableID debefix.TableID, query string, returnFieldNames []string, args ...any) (map[string]any, error) {
			insertCount++
			if sqlQueryInterface != nil {
				return sqlQueryInterface.Query(ctx, tableID, query, returnFieldNames, args...)
			}
			return debugQI.Query(ctx, tableID, query, returnFieldNames, args...)
		})))

	fmt.Printf("INSERTED: %d\n", insertCount)

	return nil
}
