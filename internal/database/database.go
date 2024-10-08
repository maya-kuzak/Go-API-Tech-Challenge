// create database
// load and initilaize vars that are in docker compose/postgres
package database

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() (*gorm.DB, error) {

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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	log.Println("Database connection established")

	err = db.AutoMigrate(&dbName, &models.Person{}, &models.Course{}, &models.PersonCourse{})
	if err != nil {
		log.Fatal("Could not auto migrate", err)
	}
	return db, nil
}
