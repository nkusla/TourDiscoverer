package main

type PingResponse struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
}

type FollowUserRequest struct {
	Follower string `json:"follower" validate:"required"`
	Followee string `json:"followee" validate:"required"`
}

type UnfollowUserRequest struct {
	Follower string `json:"follower" validate:"required"`
	Followee string `json:"followee" validate:"required"`
}
