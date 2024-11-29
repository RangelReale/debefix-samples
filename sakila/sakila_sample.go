package main

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"flag"
// 	"fmt"
// 	"path/filepath"
// 	"runtime"
//
// 	_ "github.com/jackc/pgx/v5/stdlib"
// 	dbsql "github.com/rrgmc/debefix/db/sql"
// 	"github.com/rrgmc/debefix/db/sql/postgres"
// 	"github.com/rrgmc/debefix/v2"
// )
//
// var useDB = flag.Bool("use-db", false, "use db")
//
// func main() {
// 	flag.Parse()
//
// 	ctx := context.Background()
//
// 	err := importFixtures(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// }
//
// func importFixtures(ctx context.Context) error {
// 	var db *sql.DB
// 	var err error
//
// 	if *useDB {
// 		db, err = sql.Open("pgx",
// 			fmt.Sprintf("postgres://postgres:password@%s:%s/%s?sslmode=disable", "localhost", "5478", "sakila"))
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		fmt.Println("not using DB, dumping only")
// 	}
//
// 	curDir, err := currentSourceDirectory()
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// wrap query interface, so we can print the output statements
// 	var sqlQueryInterface dbsql.QueryInterface
//
// 	if *useDB {
// 		// will send an INSERT SQL for each row to the db, taking table dependency in account for the correct order.
// 		sqlQueryInterface = dbsql.NewSQLQueryInterface(db)
// 	}
//
// 	insertCount := 0
//
// 	debugQI := dbsql.NewDebugQueryInterface(nil)
//
// 	_, err = postgres.GenerateDirectory(ctx, filepath.Join(curDir, "fixtures"),
// 		dbsql.QueryInterfaceFunc(func(ctx context.Context, databaseName, tableName string, query string, returnFieldNames []string, args ...any) (map[string]any, error) {
// 			insertCount++
// 			if sqlQueryInterface != nil {
// 				return sqlQueryInterface.Query(ctx, databaseName, tableName, query, returnFieldNames, args...)
// 			}
// 			return debugQI.Query(ctx, databaseName, tableName, query, returnFieldNames, args...)
// 		}),
// 		debefix.WithLoadProgress(func(filename string) {
// 			fmt.Printf("Loading file %s...\n", filename)
// 		}),
// 		debefix.WithGenerateResolveCheck(true),
// 		debefix.WithResolveProgress(func(tableID, databaseName, tableName string) {
// 			fmt.Printf("Importing table %s...\n", tableName)
// 		}))
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	fmt.Printf("INSERTED: %d\n", insertCount)
//
// 	return nil
// }
//
// func currentSourceDirectory() (string, error) {
// 	_, filename, _, ok := runtime.Caller(0)
// 	if !ok {
// 		return "", errors.New("unable to get the current filename")
// 	}
// 	return filepath.Dir(filename), nil
// }
