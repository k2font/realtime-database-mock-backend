package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	passwd string
}

func (m MongoClient) New() (*mongo.Client, error) {
	// MongoDB AtlasのURLを作成
	URL := "mongodb+srv://k2font:" + m.passwd + "@cluster0.yqosybw.mongodb.net/?retryWrites=true&w=majority"

	clientOptions := options.Client().ApplyURI(URL)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	return client, err
}
