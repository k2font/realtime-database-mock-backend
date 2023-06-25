package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/olahol/melody"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Location struct {
	Name string
}

func main() {
	// Ginのセットアップ
	r := gin.Default()

	// dotenvのセットアップ
	godotenv.Load()

	URL := "mongodb+srv://k2font:" + os.Getenv("MONGODB_ATLAS_PASSWD") + "@cluster0.yqosybw.mongodb.net/?retryWrites=true&w=majority"

	// MongoDBに接続する
	clientOptions := options.Client().ApplyURI(URL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 疎通確認
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Println("Connected to MongoDB!")

	m := melody.New()

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		// DBから接続を切断
		err = client.Disconnect(
			context.TODO(),
		)

		if err != nil {
			log.Fatal(err)
		}

		log.Default().Println("Connection to MongoDB closed.")
	})

	// WebSocket接続時の処理
	m.HandleConnect(func(s *melody.Session) {
		// WebSocket接続時に全データを取得してbroadcastする
		collection := client.Database("city").Collection("locations")

		cur, err := collection.Find(context.TODO(), bson.D{})
		if err != nil {
			log.Fatal(err)
		}

		for cur.Next(context.TODO()) {
			var result Location
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}

			m.Broadcast([]byte(result.Name))
		}
	})

	// WebSocket切断時の処理
	m.HandleDisconnect(func(s *melody.Session) {

	})

	// メッセージ受信時の処理
	m.HandleMessage(func(s *melody.Session, msg []byte) {
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
	})

	// Ginの起動
	r.Run(":5001")
}
