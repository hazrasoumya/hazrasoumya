package controller

import (
	"errors"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/mapper"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	"github.com/gofrs/uuid"
)

func UpsertTaskBulletin(inputModel model.TaskBulletinUpsertInput, loggedInUserEntity *entity.LoggedInUser) *model.TaskBulletinUpsertResponse {
	entity, validationResult := mapper.MapTaskBulletinInputModelToEntity(&inputModel, loggedInUserEntity)
	response := &model.TaskBulletinUpsertResponse{}
	if validationResult.Error {
		response.Error = true
		response.Message = "TaskBulletin validation failed"
		response.ValidationErrors = validationResult.ValidationMessage
	} else {
		postgres.UpsertTaskBulletin(entity, response, loggedInUserEntity)
	}
	return response
}

func GetCustomerTakBulletinFeedback(input model.FetchCustomerFeedbackInput, salesorg string) (model.FetchCustomerFeedbackResponse, error) {
	response := model.FetchCustomerFeedbackResponse{}
	dbInput, err := mapper.ValidateCustomerFeedbackInput(input, salesorg)
	if err != nil {
		response.Error = true
		return response, err
	}
	dbOutput, err := postgres.TaskBulletinFeedbackDetails(dbInput)
	if err != nil {
		return response, err
	}
	if len(dbOutput) < 1 {
		response.Error = false
		msg := "No Data Found"
		response.Message = msg
	}
	response.CustomerFeedback = mapper.MapCustomerTaskFeedbackToModel(dbOutput)
	return response, nil
}

func InsertCustomerTaskFeedBack(inputModel model.CustomerTaskFeedBackInput, loggedInUserEntity *entity.LoggedInUser) *model.CustomerTaskFeedBackResponse {
	entity, validationResult := mapper.MapCustomerTaskFeedBackInputModelToEntity(&inputModel, loggedInUserEntity)
	response := &model.CustomerTaskFeedBackResponse{}
	if validationResult.Error {
		response.Error = true
		msg := "CustomerTaskFeedBack validation failed"
		response.Message = &msg
		response.ValidationErrors = validationResult.ValidationMessage
	} else {
		postgres.InsertCustomerTaskFeedBack(entity, response, loggedInUserEntity)
	}
	return response
}

func GetCustomerTakBulletin(input *model.ListTaskBulletinInput, salesorg string, userId string) (model.ListTaskBulletinResponse, error) {
	response := model.ListTaskBulletinResponse{}
	dbInput, err := mapper.ValidateCustomerCustomerTakBulletinInput(input, salesorg, userId)
	if err != nil {
		response.Error = true
		return response, err
	}
	dbOutput, totalPages, err := postgres.TaskBulletinDetails(dbInput)
	if err != nil {
		return response, err
	}
	if len(dbOutput) < 1 {
		response.Error = false
		msg := "No Data Found"
		response.Message = msg
	}
	response.TotalPages = totalPages
	response.TaskBulletins = mapper.MapTaskBulletinEntityToModel(dbOutput)
	return response, nil
}

func TeamToCustomerDropDown(input *model.TaskBulletinInput, loggedIn *entity.LoggedInUser) *model.TaskBulletinResponse {
	var response model.TaskBulletinResponse
	Entity, err := postgres.TeamToCustomerDropdownFetch(input, loggedIn)
	if err != nil {
		response.Error = true
		response.Message = err.Error()
		return &response
	}
	data := mapper.TeamToCustomerEntityToModel(Entity)
	response.DropDown = data
	return &response
}

func FetchPrincipalNameDropDown(input model.PrincipalDropDownInput, loggedIn *entity.LoggedInUser) (*model.PrincipalDropDownResponse, error) {
	var response model.PrincipalDropDownResponse
	team, err := mapper.MapInputForPrinciPalName(input, loggedIn.SalesOrganisaton)
	if err != nil {
		response.Error = true
		return &response, err
	}
	dbOutput, err := postgres.GetPriciPalNames(team, loggedIn.SalesOrganisaton)
	if err != nil {
		response.Error = true
		return &response, nil
	}
	if len(dbOutput) < 1 {
		response.Error = false
		msg := "No Data Found"
		response.Message = &msg
	}
	response.Data = mapper.MapPrincipalNameToEntity(dbOutput)
	return &response, nil
}

func FetchTaskBulletinTypeDropDown(input *model.TaskBulletinTitleInput, loggedIn *entity.LoggedInUser) (*model.TaskBulletinTitleResponse, error) {
	var response model.TaskBulletinTitleResponse
	var teamId string
	if input != nil {
		if input.TeamID != nil && *input.TeamID != "" {
			_, err := uuid.FromString(*input.TeamID)
			if err != nil {
				return &response, errors.New("TeamId Format Is Invalid")
			} else {
				_, erro := postgres.HasTeamID(*input.TeamID, loggedIn.SalesOrganisaton)
				if erro != nil {
					return &response, errors.New("TeamId  Not Found")
				}
			}
			teamId = *input.TeamID
		}
	}
	output, err := postgres.TypeOfTaskBulletin(teamId, loggedIn.SalesOrganisaton)
	if err != nil {
		response.Error = true
		return &response, nil
	}
	response.TypeDetails = mapper.MapTaskBulletinTypeToEntity(output)
	return &response, nil
}
