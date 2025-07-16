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

func InitDatabase() *gorm.DB {
	host := os.Getenv("TOUR_DB_HOST") // Use container name for internal Docker networking
	port := "5432"       // Use internal Docker port
	user := os.Getenv("TOUR_DB_USER")
	password := os.Getenv("TOUR_DB_PASSWORD")
	dbname := os.Getenv("TOUR_DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate tour-related models
	db.AutoMigrate(&Tour{}, &KeyPoint{})

	return db
}

func SeedTour(db *gorm.DB) {
	var count int64
	db.Model(&Tour{}).Count(&count)
	if count > 0 {
		log.Println("Tour data already exists, skipping seed")
		return
	}

	sampleTours := []Tour{
		{
			Name:           "Belgrade City Tour",
			Description:    "Explore the beautiful capital of Serbia",
			Difficulty:     DifficultyEasy,
			Tags:           "city,culture,history",
			Status:         TourStatusDraft,
			Price:          0,
			AuthorUsername: "admin123",
			KeyPoints: []KeyPoint{
				{
					Name:        "Kalemegdan Fortress",
					Description: "Historic fortress overlooking the confluence of rivers",
					Latitude:    44.8225,
					Longitude:   20.4514,
					ImageURL:    "https://example.com/kalemegdan.jpg",
					Order:       1,
				},
				{
					Name:        "Republic Square",
					Description: "Central square with National Theatre",
					Latitude:    44.8176,
					Longitude:   20.4633,
					ImageURL:    "https://example.com/republic-square.jpg",
					Order:       2,
				},
			},
		},
		{
			Name:           "Novi Sad Adventure",
			Description:    "Discover the cultural capital of Vojvodina",
			Difficulty:     DifficultyMedium,
			Tags:           "culture,fortress,danube",
			Status:         TourStatusDraft,
			Price:          0,
			AuthorUsername: "admin123",
			KeyPoints: []KeyPoint{
				{
					Name:        "Petrovaradin Fortress",
					Description: "Fortress with amazing views over Danube",
					Latitude:    45.2517,
					Longitude:   19.8369,
					ImageURL:    "https://example.com/petrovaradin.jpg",
					Order:       1,
				},
			},
		},
	}

	for _, tour := range sampleTours {
		result := db.Create(&tour)
		if result.Error != nil {
			log.Printf("Error creating tour %s: %v", tour.Name, result.Error)
		} else {
			log.Printf("Created tour: %s with %d key points", tour.Name, len(tour.KeyPoints))
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
