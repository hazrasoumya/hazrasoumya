package controller

import (
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/mapper"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
)

func UpsertFlashBulletin(inputModel model.FlashBulletinUpsertInput, loggedInUserEntity *entity.LoggedInUser) *model.FlashBulletinUpsertResponse {
	entity, validationResult := mapper.MapFlashBulletinInputModelToEntity(&inputModel, loggedInUserEntity)
	response := &model.FlashBulletinUpsertResponse{}
	if validationResult.Error {
		response.Error = true
		response.Message = "FlashBulletin validation failed"
		response.ValidationErrors = validationResult.ValidationMessage
	} else {
		postgres.UpsertFlashBulletin(entity, response)
	}
	return response
}
