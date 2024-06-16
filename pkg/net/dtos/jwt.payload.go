package dtos

import "github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"

type JwtPayload struct {
	Id   string
	Role constants.Role
}
