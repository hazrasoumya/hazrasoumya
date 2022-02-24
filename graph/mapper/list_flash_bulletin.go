package mapper

import (
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/kpi/graph/postgres"
	"github.com/eztrade/kpi/graph/util"
	"github.com/gofrs/uuid"
)

func MapListFlashBulletinToEntity(inputModel model.ListFlashBulletinInput) (*entity.FlashBulletinListInput, *model.ListFlashBulletinResponse) {
	var id *uuid.UUID
	response := &model.ListFlashBulletinResponse{}
	entity := &entity.FlashBulletinListInput{}
	if inputModel.ReceipientID != nil {
		uuid, err := uuid.FromString(*inputModel.ReceipientID)
		if err != nil {
			response.Error = true
			response.Message = "Receipient ID format is invalid!"
			return nil, response
		}
		id = &uuid
	}
	entity.ReceipientID = id
	if inputModel.TeamMemberCustomerID != nil {
		uuid, err := uuid.FromString(*inputModel.TeamMemberCustomerID)
		if err != nil {
			response.Error = true
			response.Message = "TeamMemberCustomer ID format is invalid!"
			return nil, response
		}
		id = &uuid
		entity.TeamMemberCustomerID = id
		err = postgres.CheckTeamMemberCustomerID(entity)
		if err != nil {
			response.Error = true
			response.Message = "TeamMemberCustomer does not exists!"
			return nil, response
		}
	}
	var typeValue string
	var err error
	if inputModel.Type != nil {
		entity.Type = *inputModel.Type
		typeValue, err = postgres.ValidateFlashBulletinType(entity.Type)
		if err != nil {
			response.Error = true
			response.Message = "Incorrect type"
			return nil, response
		}
	}

	if typeValue != "" {
		entity.TypeValue = typeValue
	}
	entity.Status = inputModel.IsActive
	dateFlag := false
	if inputModel.StartDate != nil {
		entity.StartDate = *inputModel.StartDate
		dateFlag = true
	}
	if inputModel.EndDate != nil {
		entity.EndDate = *inputModel.EndDate
		dateFlag = true
	}
	if dateFlag {
		if inputModel.StartDate == nil || inputModel.EndDate == nil {
			response.Error = true
			response.Message = "Both start date and end date are required for date filtering."
			return nil, response
		}
		// if entity.Type == 0 {
		// 	response.Error = true
		// 	response.Message = "Type cannot be Blank"
		// 	return nil, response
		// }

		if entity.StartDate == "" {
			response.Error = true
			response.Message = "Start Date cannot be Blank"
			return nil, response
		}

		if entity.EndDate == "" {
			response.Error = true
			response.Message = "End Date cannot be Blank"
			return nil, response
		}

		startDateObj, startDateValidateErr := util.IsValidDateWithDateObect(entity.StartDate)
		if startDateValidateErr != nil {
			response.Error = true
			response.Message = "Incorrect Start Date"
			return nil, response
		}
		endDateObj, endDateValidateErr := util.IsValidDateWithDateObect(entity.EndDate)
		if endDateValidateErr != nil {
			response.Error = true
			response.Message = "Incorrect End Date"
			return nil, response
		}

		if startDateValidateErr == nil && endDateValidateErr == nil {
			if startDateObj.After(endDateObj) {
				if !response.Error {
					response.Error = true
					response.Message = "Start date is later than end date"
					return nil, response
				}
			}
		}
	}
	return entity, response
}
