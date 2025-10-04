package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() *sql.DB {
	connStr := "host=localhost port=5434 user=admin password=admin1234 dbname=mydatabase sslmode=disable"
	// dsn := "postgres://admin:admin123@localhost:5432/mydatabase?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error validating database arguments : %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Connected to Database")
	return db
}
