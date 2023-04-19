package users

type UserCreate struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdate struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type UpdatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
