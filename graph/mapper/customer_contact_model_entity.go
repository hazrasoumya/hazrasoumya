package mapper

import (
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	"github.com/eztrade/kpi/graph/util"
	uuid "github.com/gofrs/uuid"
)

func MapCustomerContactModelToEntity(inputModel *model.CustomerContactRequest, loggedInUserEntity *entity.LoggedInUser) (*entity.CustomerContact, *model.ValidationResult) {
	var id *uuid.UUID
	var csId uuid.UUID
	result := &model.ValidationResult{Error: false}

	entity := &entity.CustomerContact{}
	if inputModel.ID != nil && *inputModel.ID != "" {
		uuid, err := uuid.FromString(*inputModel.ID)
		if err == nil {
			id = &uuid
		} else {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Customer Contact ID format is invalid!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}
	entity.ID = id

	if len(inputModel.ContactName) > 0 {
		entity.ContactName = inputModel.ContactName
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Please Provide Contact Name"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	entity.Designation = inputModel.Designation

	if len(inputModel.ContactNumber) > 0 {
		entity.ContactNumber = inputModel.ContactNumber
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Please Provide Contact Number"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	if inputModel.ContactImage != nil && *inputModel.ContactImage != "" {
		if util.IsValidUrl(*inputModel.ContactImage) {
			entity.ContactImage = inputModel.ContactImage
		} else {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Image URL Not Valid"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}

	if inputModel.CustomerID != "" {
		uuid, err := uuid.FromString(inputModel.CustomerID)
		if err == nil {
			csId = uuid
		} else {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Customer ID format is invalid!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}

	country, err := postgres.GetCountryByID(loggedInUserEntity.SalesOrganisaton)
	if err != nil {
		errorMessage := &model.ValidationMessage{Row: 0, Message: err.Error()}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	if inputModel.EmailID == nil || *inputModel.EmailID == "" {
		if country == "VN" {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Email Id is mandatory"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	} else {
		if util.IsValidEmail(*inputModel.EmailID) {
			entity.EmailID = inputModel.EmailID
		} else {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid Email"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}

	}

	errUsr := postgres.CustomersExist([]string{inputModel.CustomerID})
	if errUsr != nil {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid Customer"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	entity.CustomerID = csId

	return entity, result
}

func MapDeleteCustomerContactModelToEntity(inputModel *model.CustomerContactDeleteRequest) (*entity.CustomerContactDelete, *model.ValidationResult) {
	var id uuid.UUID
	result := &model.ValidationResult{Error: false}

	entity := &entity.CustomerContactDelete{}
	uuid, err := uuid.FromString(inputModel.ID)
	if err == nil {
		id = uuid
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Customer Contact ID format is invalid!"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	entity.ID = id

	return entity, result
}
