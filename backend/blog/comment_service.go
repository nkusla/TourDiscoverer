package main

type CommentService struct {
	repository *CommentRepository
}

func (s *CommentService) CreateComment(comment *Comment) error {
	return s.repository.Create(comment)
}

func (s *CommentService) GetComments(blogID string) ([]Comment, error) {
	return s.repository.GetByBlogID(blogID)
}
