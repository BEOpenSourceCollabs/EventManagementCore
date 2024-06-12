package dtos

type GoogleSignUpRequest struct {
	IdToken   string `json:"id_token" validate:"required"`
	FirstName string `json:"first_name" validate:"omitempty,alpha,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,alpha,max=50"`
	Username  string `json:"username" validate:"required,max=20"`
}

type GoogleSignInRequest struct {
	IdToken string `json:"id_token" validate:"required"`
}
