package mapper

import (
	"encoding/json"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
)

func MapKpiAnswerEntityToModel(inputEntity entity.KpiAnswerData) (model.KpiAnswer, error) {
	var outputModel model.KpiAnswer
	KpiAnswerEntity := []model.KPIAnswerStruct{}
	err := json.Unmarshal([]byte(inputEntity.Answer), &KpiAnswerEntity)
	if err != nil {
		return outputModel, err
	}

	tempAnswers := make([]model.KPIAnswerRes, 0)
	for _, eachDesign := range KpiAnswerEntity {
		var tempAnswer model.KPIAnswerRes
		tempAnswer.QuestioNnumber = eachDesign.QuestioNnumber
		tempAnswer.Value = eachDesign.Value
		tempAnswers = append(tempAnswers, tempAnswer)
	}

	outputModel.ID = inputEntity.ID
	outputModel.Answer = tempAnswers
	outputModel.KpiID = inputEntity.KpiId
	outputModel.KpiVersionID = inputEntity.KpiVersionId
	outputModel.Category = int(inputEntity.Category)
	outputModel.TeamMemberCustomerID = inputEntity.TeamMemberCustomerID
	outputModel.ScheduleEvent = inputEntity.ScheduleEventID
	outputModel.TargetItem = inputEntity.TargetItem

	return outputModel, nil
}
