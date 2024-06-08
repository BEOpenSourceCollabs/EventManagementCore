package dtos

type LoginUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

type LoginSuccess struct {
	User        LoginUser `json:"user"`
	AccessToken string    `json:"access_token"`
}
