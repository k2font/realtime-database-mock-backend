package main

import (
	"context"
	"log"

	"github.com/olahol/melody"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConnectHandler(client *mongo.Client) func(s *melody.Session) {
	return func(s *melody.Session) {
		// WebSocket接続時に全データを取得してbroadcastする
		collection := client.Database("chat").Collection("message")

		// 該当のコレクションから全データを取得する
		// bson.D{}はフィルターなしを意味する
		cur, err := collection.Find(context.TODO(), bson.D{})
		if err != nil {
			log.Fatal(err)
		}

		// 次のカーソルがある限りループを繰り返してデータを取得する
		for cur.Next(context.TODO()) {
			var result Message
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}

			s.Write([]byte(result.Data["message"].(string)))
		}
	}
}
