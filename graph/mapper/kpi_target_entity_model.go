package mapper

import (
	"encoding/json"
	"math"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
)

func MapKpiTargetEntityToModel(inputEntity entity.KpiTargetData, getActualWorkingDays []entity.CallPlanData) (model.KpiTarget, error) {
	var outputModel model.KpiTarget
	KpiTargetEntity := []entity.KpiTarget{}
	err := json.Unmarshal([]byte(inputEntity.Target.String), &KpiTargetEntity)
	if err != nil {
		return outputModel, err
	}

	otherWorkDayMap := make(map[int64]float64)
	workingDayMap := make(map[int64]float64)
	actualWorkingDayMap := make(map[int64]float64)
	sellingTimeMap := make(map[int64]float64)
	workedDayMap := make(map[int64]float64)

	//Workable Days, Vacation Days, Other Work, Working Days
	tempTargets := make([]model.KpiTargetRes, 0)
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

		if *value == "otherwork" {
			for _, eachValue := range eachDesign.Values {
				otherWorkDayMap[eachValue.Month] = otherWorkDayMap[eachValue.Month] + eachValue.Value
			}
		}
		if *value == "workingdays" {
			for _, eachValue := range eachDesign.Values {
				workingDayMap[eachValue.Month] = workingDayMap[eachValue.Month] + eachValue.Value
			}
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
			tempTargets = append(tempTargets, tempTarget)
		}
	}

	//Number of FTE
	var tempTarget1 model.KpiTargetRes
	title1, err := postgres.GetCodeTitleFromValue("numberoffte", "KPITargetTitle")
	if err != nil {
		tempTarget1.KpiTitle = ""
		tempTarget1.KpiValue = ""
	} else {
		tempTarget1.KpiTitle = *title1
		tempTarget1.KpiValue = "numberoffte"
	}
	tempValues := make([]model.TargetValueRes, 0)
	for i := 1; i <= 12; i++ {
		var tempValue model.TargetValueRes
		tempValue.Month = int64(i)
		tempValue.Value = 0
		tempValues = append(tempValues, tempValue)
	}
	tempTarget1.Values = tempValues
	tempTargets = append(tempTargets, tempTarget1)

	//Selling Days (getActualWorkingDays - Unplanned Leaves, where Unplanned Leaves currently not included)
	var tempTarget2 model.KpiTargetRes
	title2, err := postgres.GetCodeTitleFromValue("selldays", "KPITargetTitle")
	if err != nil {
		tempTarget2.KpiTitle = ""
		tempTarget2.KpiValue = ""
	} else {
		tempTarget2.KpiTitle = *title2
		tempTarget2.KpiValue = "selldays"
	}
	if len(getActualWorkingDays) > 0 {
		for _, eachDesign := range getActualWorkingDays {
			actualWorkingDayMap[eachDesign.Month.Int64] = float64(eachDesign.Value.Int64)
		}

		tempValues := make([]model.TargetValueRes, 0)
		for i := 1; i <= 12; i++ {
			var tempValue model.TargetValueRes
			tempValue.Month = int64(i)
			tempValue.Value = actualWorkingDayMap[int64(i)]
			tempValues = append(tempValues, tempValue)
		}
		tempTarget2.Values = tempValues
		tempTargets = append(tempTargets, tempTarget2)
	} else {
		tempValues := make([]model.TargetValueRes, 0)
		for i := 1; i <= 12; i++ {
			var tempValue model.TargetValueRes
			tempValue.Month = int64(i)
			tempValue.Value = 0
			tempValues = append(tempValues, tempValue)
		}
		tempTarget2.Values = tempValues
		tempTargets = append(tempTargets, tempTarget2)
	}

	//Selling Time % Note: actualWorkingDayMap = Selling Days here
	var tempTarget3 model.KpiTargetRes
	title3, err := postgres.GetCodeTitleFromValue("selltimepercentage", "KPITargetTitle")
	if err != nil {
		tempTarget3.KpiTitle = ""
		tempTarget3.KpiValue = ""
	} else {
		tempTarget3.KpiTitle = *title3
		tempTarget3.KpiValue = "selltimepercentage"
	}
	if len(actualWorkingDayMap) > 0 && len(workingDayMap) > 0 {
		for i := 1; i <= 12; i++ {
			if workingDayMap[int64(i)] > 0 {
				sellingTimeMap[int64(i)] = math.Round(((actualWorkingDayMap[int64(i)]/workingDayMap[int64(i)])*100)*100) / 100
			} else {
				sellingTimeMap[int64(i)] = 0
			}
		}

		tempValues := make([]model.TargetValueRes, 0)
		for i := 1; i <= 12; i++ {
			var tempValue model.TargetValueRes
			tempValue.Month = int64(i)
			tempValue.Value = sellingTimeMap[int64(i)]
			tempValues = append(tempValues, tempValue)
		}
		tempTarget3.Values = tempValues
		tempTargets = append(tempTargets, tempTarget3)
	} else {
		tempValues := make([]model.TargetValueRes, 0)
		for i := 1; i <= 12; i++ {
			var tempValue model.TargetValueRes
			tempValue.Month = int64(i)
			tempValue.Value = 0
			tempValues = append(tempValues, tempValue)
		}
		tempTarget3.Values = tempValues
		tempTargets = append(tempTargets, tempTarget3)
	}

	//Worked Days
	var tempTarget4 model.KpiTargetRes
	title4, err := postgres.GetCodeTitleFromValue("workeddays", "KPITargetTitle")
	if err != nil {
		tempTarget4.KpiTitle = ""
		tempTarget4.KpiValue = ""
	} else {
		tempTarget4.KpiTitle = *title4
		tempTarget4.KpiValue = "workeddays"
	}
	if len(otherWorkDayMap) > 0 {
		for i := 1; i <= 12; i++ {
			if actualWorkingDayMap[int64(i)] >= otherWorkDayMap[int64(i)] {
				workedDayMap[int64(i)] = actualWorkingDayMap[int64(i)] - otherWorkDayMap[int64(i)]
			}
		}
		tempValues := make([]model.TargetValueRes, 0)
		for i := 1; i <= 12; i++ {
			var tempValue model.TargetValueRes
			tempValue.Month = int64(i)
			tempValue.Value = workedDayMap[int64(i)]
			tempValues = append(tempValues, tempValue)
		}
		tempTarget4.Values = tempValues
		tempTargets = append(tempTargets, tempTarget4)
	} else {
		tempValues := make([]model.TargetValueRes, 0)
		for i := 1; i <= 12; i++ {
			var tempValue model.TargetValueRes
			tempValue.Month = int64(i)
			tempValue.Value = 0
			tempValues = append(tempValues, tempValue)
		}
		tempTarget4.Values = tempValues
		tempTargets = append(tempTargets, tempTarget4)
	}

	statusName := postgres.GetCategoryFromCode(inputEntity.Status.Int64)

	outputModel.ID = inputEntity.ID.String
	outputModel.Year = inputEntity.Year.Int64
	if inputEntity.SalesRep.Valid {
		outputModel.SalesRep = inputEntity.SalesRep.String
	}
	if inputEntity.TeamName.Valid {
		outputModel.TeamName = inputEntity.TeamName.String
	}
	outputModel.Region = inputEntity.Region.String
	outputModel.Country = inputEntity.Country.String
	outputModel.Currency = inputEntity.Currency.String
	outputModel.Plants = inputEntity.Plants.Int64
	outputModel.Bergu = inputEntity.Bergu.String
	outputModel.Status = statusName
	outputModel.Target = tempTargets

	return outputModel, nil
}

func MapKpiTargetTitleEntityToModel(inputEntity entity.KpiTargetTitleList) (outputModel model.KpiTargetTitle, err error) {
	outputModel.Description = inputEntity.Description.String
	outputModel.Title = inputEntity.Title.String
	outputModel.Value = inputEntity.Value.String
	return outputModel, nil
}
