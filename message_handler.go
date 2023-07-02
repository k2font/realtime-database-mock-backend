package main

import (
	"context"
	"log"

	"github.com/olahol/melody"
	"go.mongodb.org/mongo-driver/mongo"
)

func MessageHandler(client *mongo.Client, m *melody.Melody) func(s *melody.Session, msg []byte) {
	return func(s *melody.Session, msg []byte) {
		// DBにデータを登録する
		collection := client.Database("city").Collection("locations")
		data := Location{string(msg)}

		insertResult, err := collection.InsertOne(context.Background(), data)
		if err != nil {
			log.Fatal(err)
		}

		log.Default().Println("Insert a single document: ", insertResult.InsertedID)

		// 登録したらDBからデータを一括取得

		// DBに登録されたデータをブロードキャスト
		m.Broadcast(msg)
	}
}
