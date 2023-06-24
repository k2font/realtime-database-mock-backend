package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/olahol/melody"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Location struct {
	Name string
}

func main() {
	// MongoDBに接続する
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 疎通確認
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	m := melody.New()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		// DBから接続を切断
		err = client.Disconnect(
			context.Background(),
		)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connection to MongoDB closed.")
	})

	// WebSocket接続時の処理
	m.HandleConnect(func(s *melody.Session) {

	})

	// WebSocket切断時の処理
	m.HandleDisconnect(func(s *melody.Session) {

	})

	// メッセージ受信時の処理
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		// DBにデータを登録する
		collection := client.Database("city").Collection("locations")
		tokyo := Location{"Tokyo"}

		insertResult, err := collection.InsertOne(context.Background(), tokyo)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Insert a single document: ", insertResult.InsertedID)

		// 登録したらDBからデータを一括取得

		// DBに登録されたデータをブロードキャスト
		m.Broadcast(msg)
	})

	http.ListenAndServe(":5001", nil)
}
