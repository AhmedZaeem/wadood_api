package db

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

var DB *sql.DB

func InitDB() error {
    var err error
    dsn := "root:@tcp(localhost:3306)/wadood"
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    if err = DB.Ping(); err != nil {
        return err
    }
    return createTables()
}

func createTables() error {
    query := `CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT,
        username VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        PRIMARY KEY (id)
    );`
    _, err := DB.Exec(query)
    if err != nil {
        log.Printf("Error creating tables: %v", err)
    }
    return err
}