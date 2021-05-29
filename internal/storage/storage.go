package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect(connString string) *mongo.Client {
	var err error

	client, err := mongo.NewClient(options.Client().ApplyURI(connString))
	if err != nil {
		fmt.Println("error trying to create new client for mongodb")
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("error trying to connect to mongodb")
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("error trying to Ping to mongodb connection")
		panic(err)
	}

	fmt.Println("mongodb connected!")
	return client
}
