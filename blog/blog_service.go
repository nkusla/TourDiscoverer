package main

type BlogService struct {
	repository *BlogRepository
}

func (s *BlogService) CreateBlog(blog *Blog) error {
	return s.repository.Create(blog)
}
func (s *BlogService) GetAllBlogs() ([]Blog, error) {
	return s.repository.GetAll()
}
