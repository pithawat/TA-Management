package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() *sql.DB {

	_ = godotenv.Load()
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")

	log.Printf("=== Database Connection Config ===")
	log.Printf("DB_TEST_HOST: %s", db_host)
	log.Printf("DB_TEST_PORT: %s", db_port)
	log.Printf("DB_TEST_USER: %s", db_user)
	log.Printf("DB_TEST_NAME: %s", db_name)
	log.Printf("=================================")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db_host, db_port, db_user, os.Getenv("DB_TEST_PASSWORD"), db_name)
	// connStr := "host=localhost port=5434 user=admin password=admin1234 dbname=mydatabase sslmode=disable"
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
docker run --rm -it ghcr.io/pithawat/ta-management:61