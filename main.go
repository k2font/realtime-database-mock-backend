package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/olahol/melody"
)

type Location struct {
	Name string
}

func main() {
	// Ginのセットアップ
	r := gin.Default()

	// dotenvのセットアップ
	godotenv.Load()

	// MongoDBに接続する
	mongo := MongoClient{passwd: os.Getenv("MONGODB_ATLAS_PASSWD")}
	client, err := mongo.New()
	if err != nil {
		log.Fatal(err)
	}

	// mongoの疎通確認
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Println("Connected to MongoDB!")

	// Melodyのセットアップ
	m := melody.New()

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	r.GET("close", func(c *gin.Context) {
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
	m.HandleConnect(ConnectHandler(client))

	// WebSocket切断時の処理
	m.HandleDisconnect(DisconnectHandler())

	// メッセージ受信時の処理
	m.HandleMessage(MessageHandler(client, m))

	// Ginの起動
	r.Run(":5001")
}
