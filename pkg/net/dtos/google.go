package dtos

import (
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

type GoogleSignUpRequest struct {
	DTO
	IdToken   string `json:"id_token" validate:"required"`
	FirstName string `json:"first_name" validate:"omitempty,alpha,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,alpha,max=50"`
	Username  string `json:"username" validate:"required,max=20"`
}

// Validate implements validatable returns any validation errors
func (dto *GoogleSignUpRequest) Validate() (errs []string) {
	if len(dto.IdToken) < 1 {
		errs = append(errs, "id_token is required")
	}
	// Validate firstname
	if utils.ContainsNoneAlphabeticCharacters(dto.FirstName) {
		errs = append(errs, "firstname must contain only alphabetic characters without numbers, symbols or spaces")
	}
	if !utils.StringLengthInBounds(dto.FirstName, 1, 50) {
		errs = append(errs, "firstname must contain between 1 and 50 characters")
	}
	// Validate lastname
	if utils.ContainsNoneAlphabeticCharacters(dto.LastName) {
		errs = append(errs, "lastname must contain only alphabetic characters without numbers, symbols or spaces")
	}
	if !utils.StringLengthInBounds(dto.LastName, 1, 50) {
		errs = append(errs, "lastname must contain between 1 and 50 characters")
	}

	// Validate username
	if !utils.IsAlphaNumeric(dto.Username) {
		errs = append(errs, "username must contain only alphanumberic characters")
	}
	if !utils.StringLengthInBounds(dto.Username, 5, 20) {
		errs = append(errs, "username must contain between 5 and 20 characters")
	}
	return errs
}

type GoogleSignInRequest struct {
	DTO
	IdToken string `json:"id_token" validate:"required"`
}

// Validate implements validatable returns any validation errors
func (dto *GoogleSignInRequest) Validate() (errs []string) {
	if len(dto.IdToken) < 1 {
		errs = append(errs, "id_token is required")
	}
	return errs
}
