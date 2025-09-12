package main

import "errors"

type BlogService struct {
	repository *BlogRepository
}

func (s *BlogService) CreateBlog(blog *Blog) error {
	return s.repository.Create(blog)
}

func (s *BlogService) GetAllBlogs() ([]Blog, error) {
	return s.repository.GetAll()
}

func (s *BlogService) GetBlogByID(id string) (*Blog, error) {
	return s.repository.GetByID(id)
}

func (s *BlogService) ToggleLike(blogID, username string) (bool, error) {

	isLiked, err := s.repository.IsLikedByUser(blogID, username)
	if err != nil {
		return false, err
	}

	if isLiked {

		err = s.repository.RemoveLike(blogID, username)
		return false, err
	} else {

		err = s.repository.AddLike(blogID, username)
		return true, err
	}
}

func (s *BlogService) GetLikeStatus(blogID, username string) (bool, int, error) {
	blog, err := s.repository.GetByID(blogID)
	if err != nil {
		return false, 0, err
	}

	if blog == nil {
		return false, 0, errors.New("blog not found")
	}

	isLiked, err := s.repository.IsLikedByUser(blogID, username)
	if err != nil {
		return false, 0, err
	}

	return isLiked, blog.LikeCount, nil
}
