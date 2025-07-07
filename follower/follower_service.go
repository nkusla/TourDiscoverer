package main

type FollowerService struct {
	repository *FollowerRepository
}

func (s *FollowerService) CreateUser(username string) error {
	err := s.repository.CreateUser(username)
	if err != nil {
		return err
	}

	return nil
}

func (s *FollowerService) FollowUser(follower string, followee string) error {
	err := s.repository.CreateFollowRelationship(follower, followee)
	if err != nil {
		return err
	}

	return nil
}

func (s *FollowerService) UnfollowUser(follower string, followee string) error {
	err := s.repository.DeleteFollowRelationship(follower, followee)
	if err != nil {
		return err
	}

	return nil
}
