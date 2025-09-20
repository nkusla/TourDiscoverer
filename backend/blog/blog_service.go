package main

import "errors"

type BlogService struct {
	repository *BlogRepository
}

func (s *BlogService) CreateBlog(blog *Blog) error {
	return s.repository.Create(blog)
}

func (s *BlogService) GetAllBlogs(username string) ([]Blog, error) {
	blogs, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	
	// Ako je username prazan (neulogovan korisnik), vrati blog-ove bez like status-a
	if username == "" {
		return blogs, nil
	}
	
	// Za ulogovane korisnike, dodaj is_liked_by_user informaciju
	for i := range blogs {
		isLiked, err := s.repository.IsLikedByUser(blogs[i].ID, username)
		if err != nil {
			// Ako ima greška, postavi na false
			blogs[i].IsLikedByUser = false
		} else {
			blogs[i].IsLikedByUser = isLiked
		}
	}
	
	return blogs, nil
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
