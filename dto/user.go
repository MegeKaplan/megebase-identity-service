package dto

type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6"`
	Name     *string `json:"name,omitempty" binding:"omitempty,min=2"`
}
