package main

import (
  "github.com/wiliamvj/client-server-api/src/client"
  "github.com/wiliamvj/client-server-api/src/server"
)

func main() {
  go server.HTTPServer()
  client.Client()
}
