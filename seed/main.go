package main

import (
	"log"
	"seed/seeder"
)

func main() {
	log.Println("Seeding the database...")

	seederManager := seeder.NewSeederManager()
	seederManager.SeedAll()

	log.Println("Database seeding completed successfully.")
}
