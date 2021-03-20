package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jasonlvhit/gocron"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	pb "github.com/haquenafeem/basic-microservice/proto/todo"
)

func main() {
	gocron.Every(1).Second().Do(getAllTodos)

	<-gocron.Start()
}

func getAllTodos() {
	client, err := getDBClient()
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := client.Database("TODO").Collection("todo")

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	defer cur.Close(ctx)
	var results []*pb.Todo

	for cur.Next(ctx) {
		var elem pb.Todo
		err := cur.Decode(&elem)
		if err != nil {
			panic(err)
		}

		results = append(results, &elem)
	}

	fmt.Println(results)
}

func getDBClient() (*mongo.Client, error) {
	cs := "mongodb://mongo:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cs))
	if err != nil {
		return nil, err
	}

	// defer client.Disconnect(ctx)

	return client, nil
}
