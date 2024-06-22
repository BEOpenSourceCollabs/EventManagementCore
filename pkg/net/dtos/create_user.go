package dtos

import (
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/types"
)

type CreateOrUpdateUser struct {
	Register
	Role     types.Role `json:"role"`
	Verified *bool      `json:"verified,omitempty"`
}
