package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	host := os.Getenv("AUTH_DB_HOST")
	port := os.Getenv("AUTH_DB_PORT")
	user := os.Getenv("AUTH_DB_USER")
	password := os.Getenv("AUTH_DB_PASSWORD")
	dbname := os.Getenv("AUTH_DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&User{})

	return db
}

func SeedDB(db *gorm.DB) {
	// Check if data already exists
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		log.Println("Test data already exists, skipping seed")
		return
	}

	// Create test users
	testUsers := []User{
		{
			Username: "admin123",
			Password: hashPassword("admin123"),
			Email:    "admin123@gmail.com",
			Role:     RoleAdmin,
		},
		{
			Username: "tourist123",
			Password: hashPassword("tourist123"),
			Email:    "tourist123@gmail.com",
			Role:     RoleTourist,
		},
		{
			Username: "guide123",
			Password: hashPassword("guide123"),
			Email:    "guide123@gmail.com",
			Role:     RoleGuide,
		},
	}

	for _, user := range testUsers {
		result := db.Create(&user)
		if result.Error != nil {
			log.Printf("Error creating test user %s: %v", user.Username, result.Error)
		} else {
			log.Printf("Created test user: %s", user.Username)
		}
	}
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	return string(hash)
}
