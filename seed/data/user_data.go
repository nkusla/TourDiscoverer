package data

const (
	RoleTourist = "tourist"
	RoleGuide   = "guide"
	RoleAdmin   = "admin"
)

var Users = []map[string]interface{}{
	{
		"username": "tourist1",
		"password": "tourist1",
		"email":    "tourist1@gmail.com",
		"role":     RoleTourist,
	},
	{
		"username": "tourist2",
		"password": "tourist2",
		"email":    "tourist2@gmail.com",
		"role":     RoleTourist,
	},
	{
		"username": "tourist3",
		"password": "tourist3",
		"email":    "tourist3@gmail.com",
		"role":     RoleTourist,
	},
	{
		"username": "guide1",
		"password": "guide1",
		"email":    "guide1@gmail.com",
		"role":     RoleGuide,
	},
	{
		"username": "guide2",
		"password": "guide2",
		"email":    "guide2@gmail.com",
		"role":     RoleGuide,
	},
	{
		"username": "guide3",
		"password": "guide3",
		"email":    "guide3@gmail.com",
		"role":     RoleGuide,
	},
}
