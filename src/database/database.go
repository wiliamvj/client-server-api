package database

import (
  "database/sql"
  "fmt"
)

const (
  dbFile = "currencies.db"
  table  = "currencies"
)

func ConnectDB() (*sql.DB, error) {
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
