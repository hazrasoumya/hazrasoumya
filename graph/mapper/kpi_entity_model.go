package mapper

import (
	"encoding/json"
	"strings"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
)

func MapKpiEntityToModel(inputEntity entity.KpiData) (model.GetKpi, error) {
	var outputModel model.GetKpi
	ProductDesignEntity := []model.KPIDesign{}
	BrandDesignEntity := []model.KPIDesign{}

	if len(inputEntity.ProductDesign.String) > 0 {
		err := json.Unmarshal([]byte(inputEntity.ProductDesign.String), &ProductDesignEntity)
		if err != nil {
			return outputModel, err
		}
	}

	if len(inputEntity.BrandDesign.String) > 0 {
		err2 := json.Unmarshal([]byte(inputEntity.BrandDesign.String), &BrandDesignEntity)
		if err2 != nil {
			return outputModel, err2
		}
	}

	tempProductDesigns := make([]model.KPIDesignRes, 0)
	tempBrandDesigns := make([]model.KPIDesignRes, 0)

	for _, eachDesign := range ProductDesignEntity {
		var tempDesign model.KPIDesignRes
		tempDesign.Active = *eachDesign.Active
		tempDesign.CategoryID = int(eachDesign.Category)
		tempDesign.Category = postgres.GetCategoryFromCode(eachDesign.Category)
		tempDesign.Name = *eachDesign.Name
		tempDesign.Type = *eachDesign.Type
		tempDesign.EffectiveStartDate = *eachDesign.EffectiveStartDate
		tempDesign.EffectiveEndDate = *eachDesign.EffectiveEndDate
		tempQuestions := make([]model.KPIQuestionRes, 0)
		for _, eachQuestion := range eachDesign.Questions {
			var tempQuestion model.KPIQuestionRes
			tempQuestion.QuestionNumber = eachQuestion.QuestionNumber
			tempQuestion.Active = *eachQuestion.Active
			tempQuestion.OptionValues = eachQuestion.OptionValues
			tempQuestion.Required = *eachQuestion.Required
			tempQuestion.Title = eachQuestion.Title
			tempQuestion.Type = eachQuestion.Type
			tempQuestions = append(tempQuestions, tempQuestion)
		}
		tempDesign.Questions = tempQuestions
		tempProductDesigns = append(tempProductDesigns, tempDesign)
	}

	for _, eachDesign := range BrandDesignEntity {
		var tempDesign model.KPIDesignRes
		tempDesign.Active = *eachDesign.Active
		tempDesign.CategoryID = int(eachDesign.Category)
		tempDesign.Category = postgres.GetCategoryFromCode(eachDesign.Category)
		tempDesign.Name = *eachDesign.Name
		tempDesign.Type = *eachDesign.Type
		tempDesign.EffectiveStartDate = *eachDesign.EffectiveStartDate
		tempDesign.EffectiveEndDate = *eachDesign.EffectiveEndDate
		tempQuestions := make([]model.KPIQuestionRes, 0)
		for _, eachQuestion := range eachDesign.Questions {
			var tempQuestion model.KPIQuestionRes
			tempQuestion.QuestionNumber = eachQuestion.QuestionNumber
			tempQuestion.Active = *eachQuestion.Active
			tempQuestion.OptionValues = eachQuestion.OptionValues
			tempQuestion.Required = *eachQuestion.Required
			tempQuestion.Title = eachQuestion.Title
			tempQuestion.Type = eachQuestion.Type
			tempQuestions = append(tempQuestions, tempQuestion)
		}
		tempDesign.Questions = tempQuestions
		tempBrandDesigns = append(tempBrandDesigns, tempDesign)
	}

	outputModel.ParentKpiID = inputEntity.ParentKpiID.String
	outputModel.ProductKpiID = inputEntity.ProductKpiID.String
	outputModel.BrandKpiID = inputEntity.BrandKpiID.String
	outputModel.ProductKpiVersionID = inputEntity.ProductKpiVersionID.String
	outputModel.BrandKpiVersionID = inputEntity.BrandKpiVersionID.String
	outputModel.KpiName = inputEntity.KpiName.String
	outputModel.TargetTeamID = inputEntity.TargetTeamID.String
	outputModel.TargetTeamName = inputEntity.TargetTeamName.String
	outputModel.EffectiveMonth = int(inputEntity.EffectiveMonth.Int64)
	outputModel.EffectiveYear = int(inputEntity.EffectiveYear.Int64)
	outputModel.IsPriority = inputEntity.IsPriority.Bool

	for _, item := range inputEntity.TargetProduct {
		outputModel.TargetProduct = append(outputModel.TargetProduct, item.String)
	}

	for _, item := range inputEntity.TargetBrand {
		outputModel.TargetBrand = append(outputModel.TargetBrand, item.String)
	}

	outputModel.ProductDesign = tempProductDesigns
	outputModel.BrandDesign = tempBrandDesigns

	return outputModel, nil
}

func MapCustomerListEntityToModel(inputEntity entity.CustomerList) (model.CustomerList, bool, error) {
	var outputModel model.CustomerList
	KpiAnswerEntity := []model.KPIAnswerStruct{}
	hasImage := false
	if inputEntity.Url.Valid {
		err := json.Unmarshal([]byte(inputEntity.Url.String), &KpiAnswerEntity)
		if err != nil {
			return outputModel, hasImage, err
		}
	}
	for _, eachDesign := range KpiAnswerEntity {
		for _, item := range eachDesign.Value {
			if strings.Contains(strings.ToLower(item), ".png") || strings.Contains(strings.ToLower(item), ".jpg") || strings.Contains(strings.ToLower(item), ".jpeg") {
				outputModel.ID = inputEntity.CustomerId.String
				outputModel.Name = inputEntity.Name.String
				hasImage = true
			}
		}
	}
	return outputModel, hasImage, nil
}

func MapProductListEntityToModel(inputEntity entity.ProductList) (model.ProductList, bool, error) {
	var outputModel model.ProductList
	KpiAnswerEntity := []model.KPIAnswerStruct{}
	hasImage := false
	if inputEntity.Url.Valid {
		err := json.Unmarshal([]byte(inputEntity.Url.String), &KpiAnswerEntity)
		if err != nil {
			return outputModel, hasImage, err
		}
	}
	for _, eachDesign := range KpiAnswerEntity {
		for _, item := range eachDesign.Value {
			if strings.Contains(strings.ToLower(item), ".png") || strings.Contains(strings.ToLower(item), ".jpg") || strings.Contains(strings.ToLower(item), ".jpeg") {
				outputModel.ID = inputEntity.ProductId.String
				outputModel.Name = inputEntity.Name.String
				hasImage = true
			}
		}
	}
	return outputModel, hasImage, nil
}

func MapBrandListEntityToModel(inputEntity entity.BrandList) (model.BrandList, bool, error) {
	var outputModel model.BrandList
	KpiAnswerEntity := []model.KPIAnswerStruct{}
	hasImage := false
	if inputEntity.Url.Valid {
		err := json.Unmarshal([]byte(inputEntity.Url.String), &KpiAnswerEntity)
		if err != nil {
			return outputModel, hasImage, err
		}
	}
	for _, eachDesign := range KpiAnswerEntity {
		for _, item := range eachDesign.Value {
			if strings.Contains(strings.ToLower(item), ".png") || strings.Contains(strings.ToLower(item), ".jpg") || strings.Contains(strings.ToLower(item), ".jpeg") {
				outputModel.ID = inputEntity.BrandId.String
				outputModel.Name = inputEntity.Name.String
				hasImage = true
			}
		}
	}
	return outputModel, hasImage, nil
}
