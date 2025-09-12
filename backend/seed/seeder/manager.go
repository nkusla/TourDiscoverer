package seeder

import "log"

type SeederManager struct {
	seeders []Seeder
}

func NewSeederManager() *SeederManager {
	return &SeederManager{
		seeders: []Seeder{
			NewAuthSeeder(),
			NewFollowerSeeder(),
			NewTourSeeder(),
			NewBlogSeeder(),
		},
	}
}

func (sm *SeederManager) SeedAll() {
	for _, seeder := range sm.seeders {
		if !seeder.IsServiceReady() {
			log.Printf("Skipping %s seeder: service not ready", seeder.Name())
			continue
		}

		seeder.Seed()
		log.Printf("%s seeder completed successfully", seeder.Name())
	}
}
