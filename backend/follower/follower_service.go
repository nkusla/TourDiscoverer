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

func (s *FollowerService) GetFollowers(username string) ([]User, error) {
	followers, err := s.repository.GetFollowers(username)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (s *FollowerService) GetFollowing(username string) ([]User, error) {
	following, err := s.repository.GetFollowing(username)
	if err != nil {
		return nil, err
	}

	return following, nil
}

func (s *FollowerService) IsFollowing(follower string, followee string) (bool, error) {
	isFollowing, err := s.repository.IsFollowing(follower, followee)
	if err != nil {
		return false, err
	}

	return isFollowing, nil
}
