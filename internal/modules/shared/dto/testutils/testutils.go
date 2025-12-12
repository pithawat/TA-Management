package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitTestDB() *sql.DB {

	// if err := godotenv.Load("../../../../../.env"); err != nil {
	// 	log.Println("No .env file found, relying on system environment variables.")
	// }
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	log.Println("Error getting CWD:", err)
	// } else {
	// 	log.Println("--- DEBUG CWD: Running tests from:", cwd) // Look at this path!
	// }
	// loadEnvFile()
	// // db_host := os.Getenv("DB_TEST_HOST")
	// // fmt.Print("host :", db_host)

	// connStr := os.ExpandEnv("host=$DB_TEST_HOST port=$DB_TEST_PORT user=$DB_TEST_USER password=$DB_TEST_PASSWORD dbname=$DB_TEST_NAME sslmode=disable")
	// Debug: Print what we're using
	db_host := os.Getenv("DB_TEST_HOST")
	db_port := os.Getenv("DB_TEST_PORT")
	db_user := os.Getenv("DB_TEST_USER")
	db_name := os.Getenv("DB_TEST_NAME")

	log.Printf("=== Database Connection Config ===")
	log.Printf("DB_TEST_HOST: %s", db_host)
	log.Printf("DB_TEST_PORT: %s", db_port)
	log.Printf("DB_TEST_USER: %s", db_user)
	log.Printf("DB_TEST_NAME: %s", db_name)
	log.Printf("=================================")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db_host, db_port, db_user, os.Getenv("DB_TEST_PASSWORD"), db_name)

	// var err error
	testDB, err := sql.Open("postgres", connStr)
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	if err = testDB.Ping(); err != nil {
		panic("Failed to ping to database: " + err.Error())
	}
	return testDB
}

func loadEnvFile() {
	// Start from current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Could not get working directory:", err)
		return
	}

	// Try to find .env file by going up directories
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			// Found .env file
			if err := godotenv.Load(envPath); err != nil {
				log.Printf("Error loading .env from %s: %v", envPath, err)
			} else {
				log.Printf("Loaded .env from: %s", envPath)
			}
			return
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root, stop searching
			log.Println("No .env file found, using environment variables")
			return
		}
		dir = parent
	}
}
