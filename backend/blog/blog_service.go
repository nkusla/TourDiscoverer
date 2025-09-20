package main

import "errors"

type BlogService struct {
	repository *BlogRepository
	httpClient *HTTPClient
}

func (s *BlogService) CreateBlog(blog *Blog) error {
	return s.repository.Create(blog)
}

func (s *BlogService) GetAllBlogs(username string) ([]Blog, error) {
	// Ako je username prazan (neulogovan korisnik), ne vraćaj nikakve blogove
	if username == "" {
		return []Blog{}, nil
	}

	// Dobij listu korisnika koje trenutni korisnik prati
	followingUsers, err := s.httpClient.GetFollowingUsers(username)
	if err != nil {
		// Ako ne možemo da dobijemo listu praćenih korisnika, vraćamo praznu listu
		return []Blog{}, nil
	}

	// Kreiraj mapu praćenih korisnika za brže pretraživanje
	followingMap := make(map[string]bool)
	for _, user := range followingUsers {
		followingMap[user.Username] = true
	}

	// Dodaj i sebe u listu da može da vidi svoje blogove
	followingMap[username] = true

	// Dobij sve blogove
	allBlogs, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	// Filtriraj blogove - prikaži samo one od korisnika koje prati
	var filteredBlogs []Blog
	for _, blog := range allBlogs {
		if followingMap[blog.Author] {
			// Dodaj is_liked_by_user informaciju
			isLiked, err := s.repository.IsLikedByUser(blog.ID, username)
			if err != nil {
				blog.IsLikedByUser = false
			} else {
				blog.IsLikedByUser = isLiked
			}
			filteredBlogs = append(filteredBlogs, blog)
		}
	}

	return filteredBlogs, nil
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
