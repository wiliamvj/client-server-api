package server

import (
  "context"
  "database/sql"
  "encoding/json"
  "fmt"
  "io"
  "log"
  "net/http"
  "time"

  _ "github.com/mattn/go-sqlite3"
)

const (
  dbFile       = "database/currencies.db"
  table        = "currencies"
  quotationURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
)

type Quotation struct {
  USDBRL struct {
    Code       string `json:"code"`
    Codein     string `json:"codein"`
    Name       string `json:"name"`
    High       string `json:"high"`
    Low        string `json:"low"`
    VarBid     string `json:"varBid"`
    PctChange  string `json:"pctChange"`
    Bid        string `json:"bid"`
    Ask        string `json:"ask"`
    Timestamp  string `json:"timestamp"`
    CreateDate string `json:"create_date"`
  } `json:"USDBRL"`
}

func connectDB() (*sql.DB, error) {
  db, err := sql.Open("sqlite3", dbFile)
  if err != nil {
    return nil, fmt.Errorf("error opening database: %v", err)
  }

  _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, currency TEXT, value REAL, timestamp INTEGER)", table))
  if err != nil {
    return nil, fmt.Errorf("error creating table: %v", err)
  }
  return db, nil
}

func saveDB(db *sql.DB, bid string) error {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
  defer cancel()

  select {
  case <-time.After(10 * time.Millisecond):
    fmt.Println("Quotation saved in DB with successfully!")
  case <-ctx.Done():
    fmt.Println("Context DB is canceled!")
  default:
  }

  _, err := db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s (currency, value, timestamp) VALUES (?, ?, ?)", table), "USD-BRL", bid, time.Now())
  if err != nil {
    return fmt.Errorf("error saving to database: %v", err)
  }

  return nil
}

func HTTPServer() {
  db, err := connectDB()
  if err != nil {
    log.Fatalf("Error connecting to database: %v", err)
  }

  http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
    defer cancel()

    select {
    case <-time.After(200 * time.Millisecond):
      fmt.Println("Quotation processed with successfully!")
    case <-ctx.Done():
      http.Error(w, "Context HTTP is canceled!", http.StatusInternalServerError)
    default:
    }

    client := http.Client{}
    req, err := http.NewRequestWithContext(ctx, "GET", quotationURL, nil)
    if err != nil {
      http.Error(w, "Error creating request", http.StatusInternalServerError)
      return
    }

    resp, err := client.Do(req)
    if err != nil {
      http.Error(w, "Error making request", http.StatusInternalServerError)
      return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
      http.Error(w, "Server returned non-200 status code", http.StatusInternalServerError)
      return
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
      http.Error(w, "Error reading response", http.StatusInternalServerError)
      return
    }

    var data Quotation
    err = json.Unmarshal(body, &data)
    if err != nil {
      http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
      return
    }

    // save to db
    bid := data.USDBRL.Bid
    if err := saveDB(db, bid); err != nil {
      http.Error(w, "Error saving to database", http.StatusInternalServerError)
      return
    }

    result := map[string]string{"dolar": data.USDBRL.Bid}
    response, err := json.Marshal(result)
    if err != nil {
      http.Error(w, "Error creating response JSON", http.StatusInternalServerError)
      return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(response)

  })

  fmt.Println("Server listening on http://localhost:8080/cotacao")
  http.ListenAndServe(":8080", nil)
}
