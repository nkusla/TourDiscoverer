package seeder

type TourSeeder struct {
	*BaseSeeder
}

func NewTourSeeder() *TourSeeder {
	return &TourSeeder{
		BaseSeeder: NewBaseSeeder("TourService", "TOUR_SERVICE_URL"),
	}
}

func (s *TourSeeder) Seed() {
	// Implement the seeding logic for blogs here
	// This could involve creating blog posts, categories, etc.
	// For now, we'll leave it empty as a placeholder
}
