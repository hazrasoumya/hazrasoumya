package controller

import (
	"strconv"
	"time"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/mapper"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
)

func SaveKpi(inputModel model.UpsertKpiRequest, loggedInUserEntity *entity.LoggedInUser) *model.KpiResponse {
	entity, validationResult := mapper.KpiModelToEntity(inputModel)
	kpiResponse := &model.KpiResponse{}
	if !validationResult.Error {
		check := postgres.SaveKpiInformation(*entity, loggedInUserEntity)
		kpiResponse = &check
	} else {
		kpiResponse.Error = true
		kpiResponse.Message = "Kpi validation failed"
		kpiResponse.ValidationErrors = validationResult.ValidationMessage
	}
	return kpiResponse
}

func GetKpiInformation(input *model.GetKpiInput, loggedInUserEntity *entity.LoggedInUser) (model.GetKpiResponse, error) {
	inputEntity, err := mapper.ValidateGetKpisInput(input, loggedInUserEntity.SalesOrganisaton)
	if err != nil {
		return model.GetKpiResponse{}, err
	}
	teams := make([]string, 0)
	response := model.GetKpiResponse{}

	getKpis, totalPages, err := postgres.GetKpiInfo(inputEntity, loggedInUserEntity, teams)
	if err != nil {
		return model.GetKpiResponse{}, err
	}

	for _, kpi := range getKpis {
		kpiModel, err := mapper.MapKpiEntityToModel(kpi)
		if err != nil {
			return response, nil
		}
		response.Data = append(response.Data, &kpiModel)
	}

	if len(getKpis) < 1 {
		msg := "No data Found"
		response.Message = &msg
		response.ErrorCode = 999
	}

	response.TotalPage = totalPages

	return response, nil
}

func GetKpiBrandProductData(inputModel model.GetBrandProductRequest, loggedInUserEntity *entity.LoggedInUser) (model.GetKpiBrandProductResponse, error) {
	response := model.GetKpiBrandProductResponse{}

	inputentity, err := mapper.MapKpiBrandProductToEntity(inputModel)
	if err != nil {
		return response, err
	}
	getKpiBrandProducts, err := postgres.GetKpiBrandProduct(&inputentity, loggedInUserEntity)
	if err != nil {
		return response, err
	}
	if len(getKpiBrandProducts) < 1 {
		response.Message = "No data Found"
		response.ErrorCode = 999
		return response, nil
	}
	var res bool
	if *inputentity.IsKpi == true {
		res = true
	} else {
		res = false
	}
	response, _ = mapper.MapKpiBrandProductToModel(getKpiBrandProducts, res)
	return response, nil

}

func GetEventKpiProductBrand(inputModel model.KpiOfflineInput, loggedInUserEntity *entity.LoggedInUser) (*model.KpiOfflineResponse, error) {

	responseModel := model.KpiOfflineResponse{}
	var timestampLength int

	//Validation
	var startDate string
	var endDate string
	var endMonth string
	var year string

	if inputModel.StartDate != nil && inputModel.EndDate != nil {
		if *inputModel.StartDate > *inputModel.EndDate {
			responseModel.Error = true
			responseModel.Message = "Start date is  greater than end date!"
			return &responseModel, nil
		}
	}
	if inputModel.StartDate != nil {
		if *inputModel.StartDate > 0 {
			timestampLength = len(strconv.Itoa(*inputModel.StartDate))
		}
		if timestampLength != 10 && timestampLength != 13 {
			responseModel.Error = true
			responseModel.Message = "Invalid timestamp !"
			return &responseModel, nil
		} else {
			var i int64
			var err error
			if timestampLength == 10 {
				i, err = strconv.ParseInt(strconv.Itoa(*inputModel.StartDate), 10, 64)
			} else {
				i, err = strconv.ParseInt(strconv.Itoa(*inputModel.StartDate/1000), 10, 64)
			}
			if err != nil {
				panic(err)
			}
			tm := time.Unix(i, 0)
			startDate = tm.Format("01/02/2006")
		}
	}

	if inputModel.EndDate != nil {
		if *inputModel.EndDate > 0 {
			timestampLength = len(strconv.Itoa(*inputModel.EndDate))
		}
		if timestampLength != 10 && timestampLength != 13 {
			responseModel.Error = true
			responseModel.Message = "Invalid timestamp !"
			return &responseModel, nil
		} else {
			var i int64
			var err error
			if timestampLength == 10 {
				i, err = strconv.ParseInt(strconv.Itoa(*inputModel.EndDate), 10, 64)
			} else {
				i, err = strconv.ParseInt(strconv.Itoa(*inputModel.EndDate/1000), 10, 64)
			}
			if err != nil {
				panic(err)
			}
			tm := time.Unix(i, 0)
			endDate = tm.Format("01/02/2006")

			endMonth = endDate[0:2]
			year = endDate[6:10]
		}
	}

	entity, err := postgres.GetKpiEvents(startDate, endDate, loggedInUserEntity.ID, loggedInUserEntity.SalesOrganisaton)
	if err != nil {
		return &responseModel, nil
	}

	for _, kpi := range entity {
		eventKpiVersions, err := mapper.MapKpiVersionEntityToModel(kpi)
		if err != nil {
			return &responseModel, nil
		}
		responseModel.GetKpiOffline = append(responseModel.GetKpiOffline, &eventKpiVersions)
	}
	productBrandData, err := postgres.GetKpiProductBrandAnswer(startDate, endDate, loggedInUserEntity.ID, endMonth, year)
	if err != nil {
		return &responseModel, nil
	}
	productBrandDataModel := mapper.MapKpiProductAnswerDataEntityToModel(productBrandData)
	if productBrandDataModel == nil {
		return &responseModel, nil
	}
	responseModel.KpiProductBrandAnswer = productBrandDataModel
	return &responseModel, nil

}
