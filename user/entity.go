package user

import "time"

type User struct {
	ID             int
	Name           string
	Email          string
	Password       string
	AvatarFileName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Occupation     string
	Roles          string
}
