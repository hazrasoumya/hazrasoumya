package mapper

import (
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	"github.com/eztrade/kpi/graph/util"
	uuid "github.com/gofrs/uuid"
)

func MapFlashBulletinInputModelToEntity(inputModel *model.FlashBulletinUpsertInput, loggedInUser *entity.LoggedInUser) (*entity.FlashBulletin, *model.ValidationResult) {
	var id *uuid.UUID
	result := &model.ValidationResult{Error: false}
	entity := &entity.FlashBulletin{}
	var customerGroups []string

	if loggedInUser.AuthRole != "sfe" {
		if loggedInUser.AuthRole != "cbm" {
			if !postgres.IsLineOneManager(loggedInUser.ID) {
				result.Error = true
				errorMessage := &model.ValidationMessage{Row: 0, Message: "You are not allowed to create Flash Bulletin!"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
				return entity, result
			}
		}
	}

	if inputModel.ID != nil {
		uuid, err := uuid.FromString(*inputModel.ID)
		if err == nil {
			id = &uuid
		} else {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Flash Bulletin ID format is invalid!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}
	entity.ID = id
	entity.TypeID = inputModel.Type
	entity.Title = inputModel.Title
	entity.Description = inputModel.Description
	entity.IsDeleted = inputModel.IsDeleted
	entity.IsActive = inputModel.IsActive
	startDateObj, startDateValidateErr := util.IsValidDateWithDateObect(inputModel.ValidityDateStart)
	if startDateValidateErr != nil {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Incorrect start date"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	endDateObj, endDateValidateErr := util.IsValidDateWithDateObect(inputModel.ValidityDateEnd)
	if endDateValidateErr != nil {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Incorrect end date"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	if startDateValidateErr == nil && endDateValidateErr == nil {
		if startDateObj.After(endDateObj) {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Start date is greater than end date"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		} else {
			if startDateObj.Equal(endDateObj) {
				entity.ValidityDate = "[" + inputModel.ValidityDateStart + ",)"
			} else {
				entity.ValidityDate = "[" + inputModel.ValidityDateStart + "," + inputModel.ValidityDateEnd + ")"
			}
		}
	}
	typeValue, err := postgres.ValidateFlashBulletinType(entity.TypeID)
	if err != nil {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: err.Error()}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	if typeValue != "" {
		if typeValue == "salesorganisation" {
			if loggedInUser.AuthRole == "cbm" {
				if len(inputModel.Recipients) > 0 {
					entity.Recipients = nil
				}
			} else {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Client business manager can only create flash bulletin for entire sales organisation"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
		} else {
			if len(inputModel.Recipients) < 1 {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Recipients are required"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
			inputModel.Recipients = util.RemoveDuplicatesFromSlice(inputModel.Recipients)
			entity.TypeValue = typeValue
			var err error
			switch typeValue {
			case "channel":
				err = postgres.ChannelsExist(inputModel.Recipients)
			case "customer":
				err = postgres.CustomersExistBySalesOrg(inputModel.Recipients, loggedInUser.SalesOrganisaton)
			case "team":
				err = postgres.TeamsExistBySalesOrg(inputModel.Recipients, loggedInUser.SalesOrganisaton)
			}
			if err != nil {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid Recipients/Recipients not belong to this sales organization"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
			entity.Recipients = inputModel.Recipients
		}
		if typeValue == "customergroup" {
			if len(inputModel.CustomerGroup) < 1 {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Please provide atleast one customer group"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			} else {
				for _, value := range inputModel.CustomerGroup {
					customerGroups = append(customerGroups, *value)

				}
			}
		} else {
			if len(inputModel.CustomerGroup) > 0 {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Flash bulletin type is not customer group"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
		}
	}
	userID, err := uuid.FromString(loggedInUser.ID)
	if err != nil {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid uuid"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	entity.CustomerGroup = customerGroups
	if inputModel.ID != nil {
		entity.ModifiedBy = userID
	} else {
		entity.CreatedBy = userID
	}

	sUuid, err := uuid.FromString(loggedInUser.SalesOrganisaton)
	if err != nil {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Sales Organization ID format is invalid!"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	entity.SalesOrgId = sUuid

	entity.ValidateData(0, result)
	entity.Attachments = MapAttachmentInputModelsToEntities(inputModel, userID, result)
	return entity, result
}

func MapAttachmentInputModelsToEntities(inputModel *model.FlashBulletinUpsertInput, userId uuid.UUID, result *model.ValidationResult) []entity.Attachment {
	var entities []entity.Attachment
	if len(inputModel.Attachments) < 1 {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Attachments are required"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	updateFlag := inputModel.ID != nil
	for _, attachment := range inputModel.Attachments {
		var entity entity.Attachment
		if updateFlag && attachment.ID != nil {
			var id *uuid.UUID
			uuid, err := uuid.FromString(*attachment.ID)
			if err == nil {
				id = &uuid
				err = postgres.AttachementBelongsToFlashBulletin(inputModel.ID, uuid)
				if err != nil {
					if !result.Error {
						result.Error = true
						result.ValidationMessage = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Attachment does not belong to Flash Bulletin!"}
					result.ValidationMessage = append(result.ValidationMessage, errorMessage)
				}
			} else {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Attachment ID format is invalid!"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
			entity.ID = id
		}
		entity.BlobName = attachment.Filename
		entity.BlobUrl = attachment.URL
		if updateFlag {
			entity.ModifiedBy = userId
		} else {
			entity.CreatedBy = userId
		}
		entity.ValidateData(0, result)
		entities = append(entities, entity)
	}
	return entities
}
