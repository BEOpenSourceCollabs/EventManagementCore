package dtos

import (
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
)

type CreateOrUpdateUser struct {
	Register
	Role     constants.Role `json:"role"`
	Verified *bool          `json:"verified,omitempty"`
}
