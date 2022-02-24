package mapper

import (
	"encoding/json"
	"strconv"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	uuid "github.com/gofrs/uuid"
)

func MapKpiTargetToEntity(inputModel *model.KPITargetInput, loggedInUserEntity *entity.LoggedInUser) (*entity.KpiTargetInput, *model.ValidationResult) {
	var id *uuid.UUID
	var salesRepId *uuid.UUID
	var teamId *uuid.UUID
	result := &model.ValidationResult{Error: false}
	kpiTargets := []entity.KpiTarget{}

	kpiTargetData := make([]model.KpiTargetRes, 0)

	if inputModel.ID != nil && *inputModel.ID != "" {
		uuid, err := uuid.FromString(*inputModel.ID)
		if err == nil {
			id = &uuid
			getKpiTargets, err := postgres.GetKpiTargetDetails(inputModel.ID, nil, nil, nil, "", 0, "id")
			if len(getKpiTargets) == 0 || err != nil {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Kpi Target ID does not exists!"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			} else {
				KpiTargetEntity := []entity.KpiTarget{}
				err := json.Unmarshal([]byte(getKpiTargets[0].Target.String), &KpiTargetEntity)
				if err != nil {
					if !result.Error {
						result.Error = true
						result.ValidationMessage = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: err.Error()}
					result.ValidationMessage = append(result.ValidationMessage, errorMessage)
				} else {
					for _, eachDesign := range KpiTargetEntity {
						var tempTarget model.KpiTargetRes
						title, value, err := postgres.GetCodeValueForKpi(*eachDesign.KpiTitle, "KPITargetTitle")
						if err != nil {
							tempTarget.KpiTitle = ""
							tempTarget.KpiValue = ""
						} else {
							tempTarget.KpiTitle = *title
							tempTarget.KpiValue = *value
						}

						if *value == "workabledays" || *value == "vacationdays" || *value == "otherwork" || *value == "workingdays" {
							tempValues := make([]model.TargetValueRes, 0)
							for _, eachValue := range eachDesign.Values {
								var tempValue model.TargetValueRes
								tempValue.Month = eachValue.Month
								tempValue.Value = eachValue.Value
								tempValues = append(tempValues, tempValue)
							}
							tempTarget.Values = tempValues
							kpiTargetData = append(kpiTargetData, tempTarget)
						}
					}
				}
			}
		} else {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Kpi Target ID format is invalid!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}

	if (inputModel.SalesRepID != nil && *inputModel.SalesRepID != "") && (inputModel.TeamID != nil && *inputModel.TeamID != "") {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Don't provide both Sales Representative ID and Team ID for KPI Target"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	} else {
		if inputModel.SalesRepID != nil && *inputModel.SalesRepID != "" {
			uuid, err := uuid.FromString(*inputModel.SalesRepID)
			if err == nil {
				if postgres.IsSalesRep(*inputModel.SalesRepID) {
					salesRepId = &uuid
				} else {
					if !result.Error {
						result.Error = true
						result.ValidationMessage = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Sales Representative ID does not exist!"}
					result.ValidationMessage = append(result.ValidationMessage, errorMessage)
				}
			} else {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Sales Representative ID format is invalid!"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
		} else if inputModel.TeamID != nil && *inputModel.TeamID != "" {
			uuid, err := uuid.FromString(*inputModel.TeamID)
			if err == nil {
				if !postgres.HasTeamId(*inputModel.TeamID) {
					teamId = &uuid
				} else {
					if !result.Error {
						result.Error = true
						result.ValidationMessage = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Team ID does not exist!"}
					result.ValidationMessage = append(result.ValidationMessage, errorMessage)
				}
			} else {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Team ID format is invalid!"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
		} else {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Please provide either Sales Representative ID or Team ID for KPI Target"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}

	if inputModel.Year < 1900 || inputModel.Year > 3000 {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Please provide valid year"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	} else {
		isPresent, value := postgres.ValidateYear(inputModel.Year, id, salesRepId, teamId)
		if isPresent {
			tempYear := strconv.Itoa(inputModel.Year)
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: tempYear + " data already enlisted for this " + value}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}

	kpiVal := []string{}
	for _, eachTarget := range inputModel.Target {
		kpiVal = append(kpiVal, eachTarget.KpiTitle)

		kpiMonth := []int64{}
		for _, eachValues := range eachTarget.Values {
			kpiMonth = append(kpiMonth, eachValues.Month)
		}
		uniqueMonth := uniqueKpiMonth(kpiMonth)
		monthLen := len(uniqueMonth)

		if monthLen != 12 {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Kpi Target Month(s) are missing!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		} else {
			for _, eachUniqueMonth := range uniqueMonth {
				if eachUniqueMonth > 12 || eachUniqueMonth < 1 {
					if !result.Error {
						result.Error = true
						result.ValidationMessage = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Inappropriate Kpi Target Month(s) tag!"}
					result.ValidationMessage = append(result.ValidationMessage, errorMessage)
				}
			}
		}
	}
	uniqueVal := uniqueKpiTitle(kpiVal)

	countCommon := 0
	kpiTargetTitle := []*int8{}
	for _, eachUniqueVal := range uniqueVal {
		if eachUniqueVal == "workabledays" || eachUniqueVal == "vacationdays" || eachUniqueVal == "otherwork" || eachUniqueVal == "workingdays" {
			title, err := postgres.GetCodeIdForKpi(eachUniqueVal, "KPITargetTitle")
			if err != nil {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Inappropriate Kpi Target(s) tag!"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			} else {
				countCommon = countCommon + 1
				kpiTargetTitle = append(kpiTargetTitle, title)
			}
		}
	}

	if countCommon != 4 {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Kpi Target(s) missing!"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	for i, eachTarget := range inputModel.Target {
		kpiTarget := entity.KpiTarget{}
		values := []entity.TargetValue{}
		if eachTarget.KpiTitle == "workabledays" || eachTarget.KpiTitle == "vacationdays" ||
			eachTarget.KpiTitle == "otherwork" || eachTarget.KpiTitle == "workingdays" {
			kpiTarget.KpiTitle = kpiTargetTitle[i]
			if id != nil {
				if len(kpiTargetData) > 0 {
					for _, eachPreviousTarget := range kpiTargetData {
						if eachTarget.KpiTitle == eachPreviousTarget.KpiValue {
							for j, eachValues := range eachTarget.Values {
								value := entity.TargetValue{}
								if eachValues.Month == eachPreviousTarget.Values[j].Month {
									// if util.IsCurrentPastMonthYear(strconv.Itoa(int(eachValues.Month)), strconv.Itoa(inputModel.Year)) && eachValues.Value != eachPreviousTarget.Values[j].Value   {
									// 	if !result.Error {
									// 		result.Error = true
									// 		result.ValidationMessage = []*model.ValidationMessage{}
									// 	}
									// 	errorMessage := &model.ValidationMessage{Row: 0, Message: "Can't edit past month's KPI Target value!"}
									// 	result.ValidationMessage = append(result.ValidationMessage, errorMessage)
									// } else {
									value.Month = eachValues.Month
									value.Value = eachValues.Value
									values = append(values, value)
									// }
								}
							}
						}
					}
				} else {
					if !result.Error {
						result.Error = true
						result.ValidationMessage = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Can't find previous KPI Target value!"}
					result.ValidationMessage = append(result.ValidationMessage, errorMessage)
				}
			} else {
				for _, eachValues := range eachTarget.Values {
					value := entity.TargetValue{}
					value.Month = eachValues.Month
					value.Value = eachValues.Value
					values = append(values, value)
				}
			}
			kpiTarget.Values = values
			kpiTargets = append(kpiTargets, kpiTarget)
		}
	}

	statusCode, err := postgres.GetCodeIdForKpi("pending", "KPITargetStatus")
	if err != nil {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid KPI Target Status"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	entity := &entity.KpiTargetInput{}

	entity.ID = id
	entity.SalesRepID = salesRepId
	entity.TeamID = teamId
	entity.Year = inputModel.Year
	entity.Target = kpiTargets
	entity.Status = statusCode

	return entity, result
}

func MapActionKPITargetModelToEntity(inputModel *model.ActionKPITargetInput) (*entity.ActionKPITarget, *model.ValidationResult) {
	var id uuid.UUID
	result := &model.ValidationResult{Error: false}

	entity := &entity.ActionKPITarget{}
	uuid, err := uuid.FromString(inputModel.ID)
	if err == nil {
		id = uuid
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "KPI Target ID format is invalid!"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	entity.ID = id

	if inputModel.Action {
		statusCode, err := postgres.GetCodeIdForKpi("approved", "KPITargetStatus")
		if err != nil {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid KPI Target Status"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
		entity.Status = statusCode
	} else {
		statusCode, err := postgres.GetCodeIdForKpi("rejected", "KPITargetStatus")
		if err != nil {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid KPI Target Status"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
		entity.Status = statusCode
	}

	return entity, result
}

func uniqueKpiMonth(valSlice []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}
	for _, entry := range valSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func uniqueKpiTitle(valSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range valSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
