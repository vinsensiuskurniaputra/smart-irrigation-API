package dto

type RegisterAuthDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UpdatePasswordDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
