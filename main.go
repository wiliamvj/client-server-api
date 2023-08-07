package main

import (
  "log"
  "net/http"

  "github.com/wiliamvj/client-server-api/src/client"
  "github.com/wiliamvj/client-server-api/src/database"
  "github.com/wiliamvj/client-server-api/src/server"
)

func main() {
  db, err := database.ConnectDB()
  if err != nil {
    log.Fatalf("Error connecting to database: %v", err)
  }
  go http.ListenAndServe(":8080", nil)
  server.HTTPServer(db)
  client.Client()
}
