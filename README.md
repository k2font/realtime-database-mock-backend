# realtime-database-mock-backend

- [Realtime Database](https://firebase.google.com/docs/database?hl=ja)をイチから再実装するプロジェクト
- backendコードを掲載するリポジトリ
  - clientコードは以下のリポジトリ
  - https://github.com/k2font/realtime-database-mock

## システム構成

**画像を挿入する**

## かんたんな仕組み
**実装中のため現時点では方針のメモ**
- クライアントからCRUD処理の通知をWebSocketのメッセージとして受け取る
  - CRUDのどの操作かにより、受け取るデータが異なる
    - Cの場合は、同時に追加するデータを受け取る
    - Uの場合は、同時に対象のドキュメントと更新するデータを受け取る
    - R/Dの場合は、対象のドキュメントを受け取る
  - 受け取るメッセージのイメージは以下の通り
    - C
      ```jsonc
      {
        "action": "create", // 操作の種類
        "document": "users", // ドキュメント名
        "data": { // コレクションに格納するデータ
          "name": "k2font",
          "age": 29
        }
      }
      ```
    - R
      ```jsonc
      {
        "action": "read", // 操作の種類
        "document": "users" // ドキュメント名
      }
      ```
    - U
      ```jsonc
      {
        "action": "update", // 操作の種類
        "document": "users", // ドキュメント名
        "data": { // 更新するデータ
          "name": "k2font",
          "age": 29
        }
      }
      ```
    - D
      ```jsonc
      {
        "action": "delete", // 操作の種類
        "document": "users" // ドキュメント名
      }
      ```
- Create処理
  - Go側で追加するコレクションデータを `interface{}` もしくは `any` で受け取り、mongoDBの指定されたドキュメント配下にデータを格納する
  - コードのイメージ
    ```go
    func (s *Server) Create(c *gin.Context) {
      // WebSocket越しにデータの受け取り
      var data interface{}
      c.BindJSON(&data)

      // dataをmongoDBに格納する
      // ...
    }
    ```
- Read処理
  - 指定したドキュメント名のコレクションを、常にWebSocketごしに渡すように実装する
- Update処理
  - Go側で更新するコレクションデータを `interface{}` もしくは `any` で受け取り、mongoDBの指定されたドキュメント配下のデータと差し替える
  - コードのイメージ
    ```go
    func (s *Server) Update(c *gin.Context) {
      // WebSocket越しにデータの受け取り
      var data interface{}
      c.BindJSON(&data)

      // dataをmongoDBに格納する
      // ...
    }
    ```
- Delete処理
  - 指定したドキュメントを削除する

## 利用している技術
- Go
- Gin
- melody(WebSocketライブラリ)
- MongoDB