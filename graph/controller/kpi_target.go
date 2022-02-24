package controller

import (
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/mapper"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	"github.com/eztrade/kpi/graph/util"
)

func SaveTargetKpi(inputModel model.KPITargetInput, loggedInUserEntity *entity.LoggedInUser) *model.KpiResponse {
	entity, validationResult := mapper.MapKpiTargetToEntity(&inputModel, loggedInUserEntity)
	response := &model.KpiResponse{}
	if !validationResult.Error {
		postgres.UpsertKpiTarget(entity, response, loggedInUserEntity)
	} else {
		response.Error = true
		response.Message = "Failed to save Kpi target data"
		response.ValidationErrors = validationResult.ValidationMessage
	}
	return response
}

func GetKpiTargetData(inputModel *model.GetKpiTargetRequest, loggedInUserEntity *entity.LoggedInUser) (model.GetKpiTargetResponse, error) {
	response := model.GetKpiTargetResponse{}
	if (inputModel.SalesRepID == nil || *inputModel.SalesRepID == "") && (inputModel.TeamID == nil || *inputModel.TeamID == "") {
		return model.GetKpiTargetResponse{Error: true, Message: "Please provide either Team or Salesrep ID!"}, nil
	} else if inputModel.SalesRepID != nil && inputModel.TeamID != nil {
		return model.GetKpiTargetResponse{Error: true, Message: "Team and SalesRep both can't be provided!"}, nil
	}

	year := 0
	if inputModel.Year == nil {
		year = util.GetCurrentTime().Year()
	} else {
		year = *inputModel.Year
	}

	getKpiTargets, err := postgres.GetKpiTargetDetails(inputModel.ID, inputModel.TeamID, inputModel.SalesRepID, inputModel.Status, loggedInUserEntity.SalesOrganisaton, year, "filter")
	if err != nil {
		return model.GetKpiTargetResponse{}, err
	}

	typeData := ""
	rtId := ""
	if inputModel.TeamID != nil {
		typeData = "team"
		rtId = *inputModel.TeamID
	} else if inputModel.SalesRepID != nil {
		typeData = "representative"
		rtId = *inputModel.SalesRepID
	} else {
		typeData = "all"
		rtId = ""
	}

	getActualWorkingDays, err := postgres.ActualWorkingDays(rtId, year, typeData)
	if err != nil {
		return model.GetKpiTargetResponse{}, err
	}

	for _, kpiTrg := range getKpiTargets {
		kpiTrgModel, err := mapper.MapKpiTargetEntityToModel(kpiTrg, getActualWorkingDays)
		if err != nil {
			return response, nil
		}
		response.GetTargets = append(response.GetTargets, &kpiTrgModel)
	}

	return response, nil
}

func ActionTargetKpi(inputModel model.ActionKPITargetInput, loggedInUserEntity *entity.LoggedInUser) *model.KpiResponse {
	entity, validationResult := mapper.MapActionKPITargetModelToEntity(&inputModel)
	response := &model.KpiResponse{}
	if !validationResult.Error {
		postgres.ActionKPITargetData(entity, response, loggedInUserEntity)
	} else {
		response.Error = true
		response.Message = "Kpi Target Details Not Valid"
		response.ValidationErrors = validationResult.ValidationMessage
	}
	return response
}

func GetKpiTargetTitle() (model.KpiTaregetTitleResponse, error) {
	response := model.KpiTaregetTitleResponse{}
	kpiTargetTitles, err := postgres.GetKpiTargetTitleInfo("KPITargetTitle")
	if err != nil {
		return response, err
	}
	for _, kpiTargetTitle := range kpiTargetTitles {
		kpiTargetTitleModel, err := mapper.MapKpiTargetTitleEntityToModel(kpiTargetTitle)
		if err != nil {
			return response, nil
		}
		response.KpiTargetTitles = append(response.KpiTargetTitles, &kpiTargetTitleModel)
	}
	return response, nil
}
