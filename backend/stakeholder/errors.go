package main

import "errors"

var (
	ErrStakeholderNotFound      = errors.New("stakeholder profile not found")
	ErrStakeholderAlreadyExists = errors.New("stakeholder profile already exists")
	ErrUsernameRequired         = errors.New("username is required")
	ErrInvalidProfileData       = errors.New("invalid profile data")
	ErrPositionNotFound         = errors.New("position not found")
)

const (
	StakeholderUsernamePrimaryKey  = "stakeholders_pkey"
	StakeholderUsernameUniqueIndex = "stakeholders_username_key"
)
