package dtos

import (
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

type DTO struct {
	utils.Validatable
}

func (dto *DTO) Validate() (errs []string) {
	// override to implement custom dto validation
	return errs
}
