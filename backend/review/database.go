package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase() *gorm.DB {
	host := os.Getenv("REVIEW_DB_HOST")
	port := "5432"
	user := os.Getenv("REVIEW_DB_USER")
	password := os.Getenv("REVIEW_DB_PASSWORD")
	dbname := os.Getenv("REVIEW_DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate review-related models
	err = db.AutoMigrate(&Review{}, &ReviewImage{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Review database connection established and migrations completed")
	return db
}
