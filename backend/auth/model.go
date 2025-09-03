package main

type User struct {
	Username  string `json:"username" gorm:"primaryKey"`
	Password  string `json:"-" gorm:"not null"` // "-" excludes from JSON
	Email     string `json:"email" gorm:"not null;uniqueIndex:idx_user_email"`
	Role      string `json:"role" gorm:"not null;default:'tourist'"`
	IsBlocked bool   `json:"is_blocked" gorm:"default:false"`
}

const (
	RoleTourist = "tourist"
	RoleGuide   = "guide"
	RoleAdmin   = "admin"
)

var UserRoles = []string{
	RoleTourist,
	RoleGuide,
	RoleAdmin,
}
