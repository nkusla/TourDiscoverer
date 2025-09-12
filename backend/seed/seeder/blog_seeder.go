package seeder

type BlogSeeder struct {
	*BaseSeeder
}

func NewBlogSeeder() *BlogSeeder {
	return &BlogSeeder{
		BaseSeeder: NewBaseSeeder("BlogService", "BLOG_SERVICE_URL"),
	}
}

func (s *BlogSeeder) Seed() {
	// Implement the seeding logic for blogs here
	// This could involve creating blog posts, categories, etc.
	// For now, we'll leave it empty as a placeholder
}
