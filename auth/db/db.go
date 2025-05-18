package db

import (
    "database/sql"
)

var DB *sql.DB

func InitDB() error {
    var err error
    dsn := "root:@tcp(127.0.0.1:3306)/wadood?parseTime=true"
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    return DB.Ping()
}