package models

import "github.com/goravel/framework/database/orm"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	orm.Model
	Name     string
	Email    string
	Alias    string
	Password string
	Role     string
	Active   bool
}

func (u User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
