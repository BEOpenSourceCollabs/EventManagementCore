package dtos

import (
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

type Register struct {
	DTO
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Validate implements validatable returns any validation errors
func (reg *Register) Validate() (errs []string) {
	// Validate email.
	// Checks the email is valid and by doing so also ensures that the email was set correcly.
	// No further validation required.
	if !utils.IsEmail(reg.Email) {
		errs = append(errs, "'%s' is not a valid email address", reg.Email)
	}
	// Validate password.
	// Ensure password is alphanumberic.
	if !utils.ContainsAlphabeticCharacters(reg.Password) || !utils.ContainsNumbericCharacters(reg.Password) {
		errs = append(errs, "password must contain both alphanumberic characters")
	}
	// Ensure password length is inbounds
	if !utils.StringLengthInBounds(reg.Password, 6, 50) {
		errs = append(errs, "password must contain between 6 and 50 characters")
	}
	// Validate firstname
	if utils.ContainsNoneAlphabeticCharacters(reg.FirstName) {
		errs = append(errs, "firstname must contain only alphabetic characters without numbers, symbols or spaces")
	}
	if !utils.StringLengthInBounds(reg.FirstName, 1, 50) {
		errs = append(errs, "firstname must contain between 1 and 50 characters")
	}
	// Validate lastname
	if utils.ContainsNoneAlphabeticCharacters(reg.LastName) {
		errs = append(errs, "lastname must contain only alphabetic characters without numbers, symbols or spaces")
	}
	if !utils.StringLengthInBounds(reg.LastName, 1, 50) {
		errs = append(errs, "lastname must contain between 1 and 50 characters")
	}
	// Validate username
	if !utils.IsAlphaNumeric(reg.Username) {
		errs = append(errs, "username must contain only alphanumberic characters")
	}
	if !utils.StringLengthInBounds(reg.Username, 5, 20) {
		errs = append(errs, "username must contain between 5 and 20 characters")
	}

	return errs
}
