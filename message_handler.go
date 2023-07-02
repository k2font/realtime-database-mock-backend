package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/olahol/melody"
	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	Action   string                 `bson:"action"`
	Document string                 `bson:"document"`
	Data     map[string]interface{} `bson:"data, omitempty"`
}

func MessageHandler(client *mongo.Client, m *melody.Melody) func(s *melody.Session, msg []byte) {
	return func(s *melody.Session, msg []byte) {
		// DBにデータを登録する
		collection := client.Database("chat").Collection("message")
		// jsonデータのメッセージをbsonデータに変換する
		// https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/
		var jsondata Message
		json.Unmarshal(msg, &jsondata)
		insertResult, err := collection.InsertOne(context.Background(), jsondata)
		if err != nil {
			log.Fatal(err)
		}

		log.Default().Println("Insert a single document: ", insertResult.InsertedID)

		// DBに登録されたデータをブロードキャスト
		m.Broadcast([]byte(jsondata.Data["message"].(string)))
	}
}
