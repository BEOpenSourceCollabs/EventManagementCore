package dtos

type Register struct {
	Email     string `json:"email" validate:"required,email,min=1,max=256"`
	Password  string `json:"password" validate:"required,is-strong,min=6,max=50"`
	FirstName string `json:"first_name" validate:"omitempty,alpha,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,alpha,max=50"`
	Username  string `json:"username" validate:"required,max=20"`
}
