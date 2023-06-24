package main

import (
	"net/http"

	"github.com/olahol/melody"
)

func main() {
	m := melody.New()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
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

		// 登録したらDBからデータを一括取得

		// DBに登録されたデータをブロードキャスト
		m.Broadcast(msg)
	})

	http.ListenAndServe(":5001", nil)
}
