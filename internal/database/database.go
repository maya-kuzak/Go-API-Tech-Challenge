// create database
// load and initilaize vars that are in docker compose/postgres
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func Init() (*sql.DB, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	dbName := os.Getenv("DATABASE_NAME")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")

	fmt.Printf("DB Env: name = %s, user = %s, password = %s, host = %s, port = %s\n", dbName, dbUser, dbPassword, dbHost, dbPort)

	//dsn - data source name
	//postgres://user:password@host:port/dbname?sslmode=disable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Verify the connection
	err = conn.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
		os.Exit(1)
	}

	log.Println("Database connection established")

	sqlFile, err := os.ReadFile("db_seed.sql")
	if err != nil {
		log.Fatal("Error reading sql file: ", err)
	}

	_, err = conn.Exec(string(sqlFile))
	if err != nil {
		log.Fatal("Error executing sql file: ", err)
	}
	log.Println("Database seeded successfully")

	return conn, nil
}
