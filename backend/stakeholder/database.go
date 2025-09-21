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
	host := os.Getenv("STAKEHOLDER_DB_HOST")
	port := os.Getenv("STAKEHOLDER_DB_PORT")
	user := os.Getenv("STAKEHOLDER_DB_USER")
	password := os.Getenv("STAKEHOLDER_DB_PASSWORD")
	dbname := os.Getenv("STAKEHOLDER_DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to stakeholder database")

	// Auto migrate the schema
	err = db.AutoMigrate(&Stakeholder{}, &TouristPosition{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed")
	return db
}
