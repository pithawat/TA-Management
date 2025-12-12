package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitTestDB() *sql.DB {

	// if err := godotenv.Load("../../../../../.env"); err != nil {
	// 	log.Println("No .env file found, relying on system environment variables.")
	// }
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting CWD:", err)
	} else {
		log.Println("--- DEBUG CWD: Running tests from:", cwd) // Look at this path!
	}
	_ = godotenv.Load()
	db_host := os.Getenv("DB_TEST_HOST")
	fmt.Print("host :", db_host)
	connStr := os.ExpandEnv("host=$DB_TEST_HOST port=$DB_TEST_PORT user=$DB_TEST_USER password=$DB_TEST_PASSWORD dbname=$DB_TEST_NAME sslmode=disable")
	// connStr := "host=localhost port=5436 user=test_user password=test_password dbname=ta_management_test sslmode=disable"
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
