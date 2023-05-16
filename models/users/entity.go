package users

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Role      string
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "tbl_user"
}
