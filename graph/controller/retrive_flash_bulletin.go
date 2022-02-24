package controller

import (
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/mapper"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	"github.com/gofrs/uuid"
)

func RetriveFlashBulletinForEdit(inputModel model.RetriveInfoFlashBulletinInput, loggedInUserEntity *entity.LoggedInUser) *model.RetriveInfoFlashBulletinleResponse {
	response := &model.RetriveInfoFlashBulletinleResponse{}
	_, err := uuid.FromString(inputModel.BulletinID)
	if err != nil {
		response.Error = true
		response.Message = "ID is not a valid UUID"
		return response
	}
	if loggedInUserEntity.AuthRole != "sfe" || loggedInUserEntity.AuthRole != "cbm" {
		if !postgres.UserIDChecking(inputModel.BulletinID, loggedInUserEntity) {
			response.Error = true
			response.Message = "You are not allowed to view flash bulletin!"
			return response
		}
	}
	flashBulletinObj, err := postgres.RetiveFlashBulletinDataById(inputModel)
	if err != nil {
		response.Error = true
		response.Message = "No data found"
		return response
	}
	attachments, err := postgres.GetAttachmentsByIDs(flashBulletinObj.Attachments.String)
	if err != nil {
		response.Error = true
		response.Message = "Error getting attachments"
		return response
	}
	recipients := make([]*model.Recipients, 0)

	if flashBulletinObj.Type.String == "Customer" {
		recipients, err = postgres.GetCustomersByIDs(flashBulletinObj.Recipients.String)
		if err != nil {
			response.Error = true
			response.Message = "Error getting Customer recipients"
			return response
		}
	} else if flashBulletinObj.Type.String == "Team" {
		recipients, err = postgres.GetTeamsByIDs(flashBulletinObj.Recipients.String)
		if err != nil {
			response.Error = true
			response.Message = "Error getting Team recipients"
			return response
		}
	} else if flashBulletinObj.Type.String == "Customer Group" {
		recipients, err = postgres.GetCustomerGroups(flashBulletinObj.Recipients.String)
		if err != nil {
			response.Error = true
			response.Message = "Error getting Customer Group recipients"
			return response
		}
	}

	flashBulletinModel, err2 := mapper.MapRetrivePlanEntityToModel(flashBulletinObj, recipients, attachments)
	if err2 != nil {
		response.Error = true
		response.Message = "No data found"
		return response
	} else {
		response.FlashBulletinData = &flashBulletinModel
	}
	return response
}
