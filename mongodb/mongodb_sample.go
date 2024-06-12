package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/rrgmc/debefix"
	debefix_mongodb "github.com/rrgmc/debefix-mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func currentSourceDirectory() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filepath.Dir(filename), nil
}

func main() {
	curDir, err := currentSourceDirectory()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	fmt.Println("loading data...")
	data, err := debefix.Load(debefix.NewDirectoryFileProvider(filepath.Join(curDir, "data"),
		debefix.WithDirectoryAsTag()))
	if err != nil {
		panic(err)
	}

	// spew.Dump(data)

	fmt.Println("connecting to mongodb...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27097"))
	if err != nil {
		panic(err)
	}

	resolveTags := []string{}

	err = debefix.ResolveCheck(data, debefix.WithResolveTags(resolveTags))
	if err != nil {
		panic(err)
	}

	_, err = debefix_mongodb.Resolve(ctx, client.Database("debefix"), data)
	if err != nil {
		panic(err)
	}
}
