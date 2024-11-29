package main

import (
	"context"
	"fmt"

	debefix_mongodb "github.com/rrgmc/debefix-mongodb/v2"
	data2 "github.com/rrgmc/debefix-samples/v2/data"
	"github.com/rrgmc/debefix/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	fmt.Println("loading data...")
	data := data2.Data()
	if data.Err() != nil {
		panic(data.Err())
	}

	// spew.Dump(data)

	fmt.Println("connecting to mongodb...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27097"))
	if err != nil {
		panic(err)
	}

	err = debefix.ResolveCheck(ctx, data)
	if err != nil {
		panic(err)
	}

	_, err = debefix_mongodb.Resolve(ctx, client.Database("debefix"), data)
	if err != nil {
		panic(err)
	}
}
