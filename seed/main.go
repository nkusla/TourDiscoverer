package main

import "log"

func main() {
	log.Println("Seeding the database...")

	seeder := NewSeederService()
	seeder.SeedAll()

	log.Println("Database seeding completed successfully.")
}
