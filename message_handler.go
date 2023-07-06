package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/olahol/melody"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// bsonへのUnmarshal処理において、構造が不明なデータをマッピングする場合、
// 対象データの型をmap[string]interface{}型にするといい。
// さらにタグにdataとomitemptyにするといい。
// これは、CとRでdataプロパティの有無が変わるため。
type Message struct {
	Action   string                 `bson:"action"`
	Document string                 `bson:"document"`
	Data     map[string]interface{} `bson:"data, omitempty"`
}

func MessageHandler(client *mongo.Client, m *melody.Melody) func(s *melody.Session, msg []byte) {
	return func(s *melody.Session, msg []byte) {
		collection := client.Database("chat").Collection("message")

		// jsonデータのメッセージをbsonデータに変換する
		// https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/
		var jsondata Message
		json.Unmarshal(msg, &jsondata)

		// jsondataのActionによって処理を分岐
		switch jsondata.Action {
		case "create":
			// データをDBに登録
			insertResult, err := collection.InsertOne(context.Background(), jsondata)
			if err != nil {
				log.Fatal(err)
			}
			log.Default().Println("Insert a single document: ", insertResult.InsertedID)
		case "read":
			// データをDBから取得
		case "update":
			// DB上のデータを更新
			// 送信されたメッセージのIDを取得し、DB上のデータを更新する
			filter := bson.D{{Key: "id", Value: jsondata.Data["id"]}}
			result, err := collection.UpdateOne(context.TODO(), filter, jsondata.Data)
			if err != nil {
				log.Fatal(err)
			}
			log.Default().Println("Update a single document: ", result.UpsertedID)
		case "delete":
			// データをDBから削除
		default:
			// 何もしない
		}

		bytes, err := json.Marshal(jsondata)
		if err != nil {
			log.Fatal(err)
		}

		// DBに登録されたデータをブロードキャスト
		m.Broadcast(bytes)
	}
}
