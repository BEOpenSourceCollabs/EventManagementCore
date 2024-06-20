package dtos

import (
	"reflect"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logger"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

type DTO struct {
	utils.Validatable
}

func (dto *DTO) Validate() (errs []string) {
	// override to implement custom dto validation
	logger.AppLogger.WarnF("DTO", "%s has no validation", reflect.TypeOf(*dto))
	return errs
}
