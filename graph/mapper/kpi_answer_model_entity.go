package mapper

import (
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	uuid "github.com/gofrs/uuid"
)

func MapKpiAnswerModelToEntity(inputModel *model.KpiAnswerRequest) (*entity.KpiAnswer, *model.ValidationResult) {
	var kpiVerId *uuid.UUID
	var tmMemCustId *uuid.UUID
	var targetItem *uuid.UUID
	var scheduleEvent *uuid.UUID
	result := &model.ValidationResult{Error: false}
	kpiAns := &entity.KpiAnswer{}

	kvUuid, err := uuid.FromString(inputModel.KpiVersionID)
	if err == nil {
		kpiVerId = &kvUuid
		isTMCId, isActive := postgres.IsKPIVersionId(kpiVerId)
		if !isTMCId {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, ErrorCode: 101, Message: "KPI version doesn't exist!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		} else if isTMCId && !isActive {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, ErrorCode: 102, Message: "Kpi is not active"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, ErrorCode: 103, Message: "Kpi Version format is invalid"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	kpiAns.KpiVersionId = kpiVerId

	if postgres.CheckCodeIDForKpi(int64(inputModel.Category), "KPICategory") {
		kpiAns.Category = inputModel.Category
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, ErrorCode: 104, Message: "Invalid Category ID"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	kpiAnsStructs := []entity.KpiAnswerStruct{}
	for _, eachAnswer := range inputModel.Answers {
		kpiAnsStruct := entity.KpiAnswerStruct{}
		kpiAnsStruct.QuestioNnumber = eachAnswer.QuestioNnumber
		kpiAnsStruct.Value = eachAnswer.Value
		kpiAnsStructs = append(kpiAnsStructs, kpiAnsStruct)
	}

	if len(kpiAnsStructs) > 0 {
		kpiAns.Answers = kpiAnsStructs
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, ErrorCode: 105, Message: "No Answer Found"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	tmcUuid, err := uuid.FromString(inputModel.TeamMemberCustomerID)
	if err == nil {
		tmMemCustId = &tmcUuid
		isTMCId := postgres.CheckTeamMemberCustomerId(tmMemCustId)
		if !isTMCId {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, ErrorCode: 106, Message: "Team Member Customer doesn't exist!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, ErrorCode: 107, Message: "Invalid Team Member Customer id!"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	kpiAns.TeamMemberCustomerID = tmMemCustId

	targetItemUuid, err := uuid.FromString(inputModel.TargetItem)
	if err == nil {
		targetItem = &targetItemUuid
		isTMCId := postgres.IsTargetItemValidForKPIAnswer(kpiVerId, targetItem)
		if !isTMCId {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, ErrorCode: 108, Message: "Target Item is not valid for this KPI Answer"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, ErrorCode: 109, Message: "Invalid Target Item id!"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	kpiAns.TargetItem = targetItem

	scheduleEventUuid, err := uuid.FromString(inputModel.ScheduleEvent)
	if err == nil {
		scheduleEvent = &scheduleEventUuid
		isSeId := postgres.CheckScheduleEventId(scheduleEvent)
		if !isSeId {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, ErrorCode: 108, Message: "Schedule Event Id not found!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, ErrorCode: 109, Message: "Invalid Schedule Event Id!"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	kpiAns.ScheduleEvent = scheduleEvent

	return kpiAns, result
}
