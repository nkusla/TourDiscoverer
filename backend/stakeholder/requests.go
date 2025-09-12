package main

type CreateStakeholderRequest struct {
	Username       string `json:"username" validate:"required"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	Biography      string `json:"biography"`
	Motto          string `json:"motto"`
}

type UpdateProfileRequest struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	Biography      string `json:"biography"`
	Motto          string `json:"motto"`
}
