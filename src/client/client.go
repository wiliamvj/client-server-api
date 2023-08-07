package client

import (
  "context"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "os"
  "time"
)

type Quotation struct {
  Dolar string
}

func Client() {
  ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
  defer cancel()

  req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
  if err != nil {
    panic(err)
  }

  res, err := http.DefaultClient.Do(req)
  if err != nil {
    fmt.Println("Erro to get quotation in client: ", err)
    return
  }
  body, err := io.ReadAll(res.Body)
  if err != nil {
    fmt.Println("Erro to read response body: ", err)
    return
  }

  var data Quotation
  err = json.Unmarshal(body, &data)
  if err != nil {
    fmt.Println("Error parsing JSON", err)
    return
  }

  defer res.Body.Close()
  f, err := os.Create("cotacao.txt")
  if err != nil {
    fmt.Println("Erro to create quotation archive: ", err)
    return
  }
  defer f.Close()
  f.WriteString(fmt.Sprintf("Dólar: %s", data.Dolar))
}
