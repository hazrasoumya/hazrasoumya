package controller

import (
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/mapper"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
)

func GetTaskBulletinReportData(input model.TaskBulletinReportInput, user *entity.LoggedInUser) (*model.TaskReportOutput, error) {
	output := model.TaskReportOutput{}
	if len(input.Tittle) < 1 {
		output.Error = true
		output.Message = "Please provide task bulletin tittle"
		return &output, nil
	}

	response, err := postgres.GetTaskBulletinReportValues(user.SalesOrganisaton, input.Tittle)
	if err != nil {
		return &output, err
	}
	if len(response) < 1 {
		output.Error = false
		output.Message = "No data found"
	}
	result := []entity.ReportData{}
	validData := mapper.MapValidTaskBulletinData(response, input.DateRange)
	data := mapper.MapLatestFeedBack(validData)
	if !input.IsExcel {
		result = mapper.MapValueOfMapToEntity(data, input.DateRange)
		output.Values = mapper.MapFeedBackReportToModel(result)
		return &output, nil
	}
	result = mapper.MapValueOfMapToEntityForReport(data, input.DateRange)
	output.Values = mapper.MapFeedBackReportToModel(result)
	retuenValue, err := mapper.GrnaretReortForTaskBulletin(output.Values, input.DateRange)
	if err != nil {
		return &retuenValue, err
	}

	return &retuenValue, nil

}
