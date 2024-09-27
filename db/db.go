package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB // Global DB variable

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1710"
	dbname   = "golangdb_postgresql"
)

// Initialize the database connection
func InitDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Assign the global DB variable
	var err error
	DB, err = sql.Open("postgres", psqlInfo) // Use global DB, not a local variable
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		panic(err)
	}
	// defer DB.Close()
	// Test the database connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
		panic(err)
	}

	// Set connection pool configurations
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createUserTable(DB)
	createTables(DB)

	fmt.Println("Established a successful connection!")
}

// Create tables in the database
func createTables(database *sql.DB) {
	if database == nil {
		log.Fatal("Database connection is nil. Did you forget to initialize it?")
		panic("Database connection is nil. Did you forget to initialize it?")
	}

	createEventTableSQL := `
	CREATE TABLE IF NOT EXISTS events (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime TIMESTAMP NOT NULL,
		user_id UUID,
		CONSTRAINT fk_user
			FOREIGN KEY(user_id)
			REFERENCES users(id)
			ON DELETE CASCADE
	);
	`

	_, err := database.Exec(createEventTableSQL)
	if err != nil {
		log.Fatalf("Error creating table in database: %v", err)
		panic(err)
	}
}

func createUserTable(database *sql.DB) {
	if database == nil {
		log.Fatal("Database connection is nil. Did you forget to initialize it?")
		panic("Database connection is nil. Did you forget to initialize it?")
	}

	createUserTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	`

	_, err := database.Exec(createUserTableSQL)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
		panic(err)
	}

	fmt.Println("Users table created successfully!")
}
