package controller

import (
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/mapper"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
)

func InsertKpiAnswers(inputModel model.KpiAnswerRequest, loggedInUserEntity *entity.LoggedInUser) *model.KpiResponse {
	entity, validationResult := mapper.MapKpiAnswerModelToEntity(&inputModel)
	kpiResponse := &model.KpiResponse{}
	if !validationResult.Error {
		postgres.InsertKpiAnswer(entity, kpiResponse, loggedInUserEntity)
	} else {
		kpiResponse.Error = true
		kpiResponse.ErrorCode = 999
		kpiResponse.Message = "Upsert kpi answer validation failed"
		kpiResponse.ValidationErrors = validationResult.ValidationMessage
	}
	return kpiResponse
}

func GetKpiAnswerData(inputModel *model.GetKpiAnswerRequest, loggedInUserEntity *entity.LoggedInUser) (model.GetKpiAnswerResponse, error) {
	response := model.GetKpiAnswerResponse{}

	getKpiAnswers, isProposedStock, err := postgres.GetKpiAnswerDetails(inputModel, loggedInUserEntity, true)
	if err != nil {
		return model.GetKpiAnswerResponse{}, err
	}

	if len(getKpiAnswers) < 1 {
		getKpiAnswers, isProposedStock, err = postgres.GetKpiAnswerDetails(inputModel, loggedInUserEntity, false)
		if err != nil {
			return model.GetKpiAnswerResponse{}, err
		}

		response.IsOldAnswer = true
	} else {
		response.IsOldAnswer = false
	}

	for _, kpiAns := range getKpiAnswers {
		kpiAnsModel, err := mapper.MapKpiAnswerEntityToModel(kpiAns)
		if err != nil {
			return response, nil
		}
		response.GetAnswers = append(response.GetAnswers, &kpiAnsModel)
	}

	if len(getKpiAnswers) < 1 {
		response.Message = "No data Found"
		response.ErrorCode = 999
		response.IsOldAnswer = false
	}

	response.IsProposedStock = isProposedStock

	return response, nil
}
