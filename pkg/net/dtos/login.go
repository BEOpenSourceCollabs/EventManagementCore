package dtos

type Login struct {
	DTO
	Email    string `json:"email"`
	Password string `json:"password"`
}
