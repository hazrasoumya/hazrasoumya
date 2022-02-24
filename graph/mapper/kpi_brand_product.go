package mapper

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	uuid "github.com/gofrs/uuid"
)

func MapKpiBrandProductToEntity(inputModel model.GetBrandProductRequest) (entity.KpiBrandProductInput, error) {
	outputModel := entity.KpiBrandProductInput{}
	if inputModel.TargetTeam != nil && *inputModel.TargetTeam != "" {
		_, err := uuid.FromString(*inputModel.TargetTeam)
		if err != nil {
			return outputModel, errors.New("Invalid Target  Team Id!")
		} else {
			_, err := postgres.HasTargetTeam(*inputModel.TargetTeam)

			if err != nil {
				return outputModel, errors.New("Target Team Id does not exist!")
			} else {
				outputModel.TargetTeam = inputModel.TargetTeam
			}
		}
	}
	if inputModel.BrandID != nil && *inputModel.BrandID != "" {
		_, err := uuid.FromString(*inputModel.BrandID)
		if err != nil {
			return outputModel, errors.New("Invalid Brand Id!")
		} else {
			_, err := postgres.HasProductBrandId(*inputModel.BrandID, "targetbrand")

			if err != nil {
				return outputModel, errors.New("Target Brand Id does not exist!")
			} else {
				outputModel.BrandID = inputModel.BrandID
			}

		}
	}
	if inputModel.ProductID != nil && *inputModel.ProductID != "" {
		_, err := uuid.FromString(*inputModel.ProductID)
		if err != nil {
			return outputModel, errors.New("Invalid Product Id!")
		} else {
			_, err := postgres.HasProductBrandId(*inputModel.ProductID, "targetproduct")

			if err != nil {
				return outputModel, errors.New("Target Product Id does not exist!")
			} else {
				outputModel.ProductId = inputModel.ProductID
			}

		}
	}
	if inputModel.TeamProductID != nil && *inputModel.TeamProductID != "" {
		_, err := uuid.FromString(*inputModel.TeamProductID)
		if err != nil {
			return outputModel, errors.New("Invalid TeamProduct Id!")
		} else {
			_, err := postgres.IsValidTeamProduct(*inputModel.TeamProductID)

			if err != nil {
				return outputModel, errors.New("Target TeamProduct Id does not exist!")
			} else {
				outputModel.TeamProductId = inputModel.TeamProductID
			}

		}
	}
	if inputModel.IsKpi != nil && *inputModel.IsKpi == true {
		if inputModel.TargetTeam == nil || *inputModel.TargetTeam == "" {
			return outputModel, errors.New("For Kpi details teamid is mandatory!")
		} else {
			outputModel.IsKpi = inputModel.IsKpi
		}

	} else {
		res := false
		outputModel.IsKpi = &res
	}
	if inputModel.IsPriority != nil {
		outputModel.IsPriority = inputModel.IsPriority
	}
	if inputModel.SearchItem != nil && *inputModel.SearchItem != "" {
		outputModel.SearchIteam = inputModel.SearchItem
	}
	outputModel.IsActive = &inputModel.IsActive

	return outputModel, nil

}

func MapKpiBrandProductToModel(getKpiBrandProducts []entity.KpiBrandProductData, isKpi bool) (model.GetKpiBrandProductResponse, error) {
	output := model.GetKpiBrandProductResponse{}
	kpiArrs := []entity.UniqueBrand{}
	for _, value := range getKpiBrandProducts {
		kpiArr := entity.UniqueBrand{}
		kpiArr.BrandID = value.BrandId.String
		kpiArr.BrandName = value.BrandName.String
		kpiArrs = append(kpiArrs, kpiArr)
	}
	UniqueList := UniqueBrand(kpiArrs)

	for _, value := range UniqueList {
		kpiBrandItem := model.KpiBrandItem{}
		kpiBrandItem.BrandID = value.BrandID
		kpiBrandItem.BrandName = value.BrandName
		kpiProducts := []*model.KpiProductItem{}

		if isKpi {
			for _, item := range getKpiBrandProducts {
				kpiProduct := model.KpiProductItem{}

				if value.BrandID == item.BrandId.String && strings.EqualFold(item.Type.String, "product") {
					if strings.EqualFold(item.TeamProductId.String, "00000000-0000-0000-0000-000000000000") {
						kpiProduct.TeamProductID = ""
					} else {
						kpiProduct.TeamProductID = item.TeamProductId.String
					}
					if strings.EqualFold(item.ProductId.String, "00000000-0000-0000-0000-000000000000") {
						kpiProduct.ProductID = ""
					} else {
						kpiProduct.ProductID = item.ProductId.String
					}
					kpiProduct.PrincipalName = item.PrincipalName.String
					kpiProduct.MaterialDescription = item.MaterialDescription.String
					kpiProduct.IsPriority = item.IsPriority.Bool
					kpiProducts = append(kpiProducts, &kpiProduct)
				}
			}

		} else {
			for _, item := range getKpiBrandProducts {
				kpiProduct := model.KpiProductItem{}
				if value.BrandID == item.BrandId.String {

					kpiProduct.TeamProductID = item.TeamProductId.String

					kpiProduct.ProductID = item.ProductId.String

					kpiProduct.PrincipalName = item.PrincipalName.String
					kpiProduct.MaterialDescription = item.MaterialDescription.String
					kpiProduct.IsPriority = item.IsPriority.Bool
					kpiProducts = append(kpiProducts, &kpiProduct)
				}
			}
		}
		kpiBrandItem.Products = kpiProducts

		output.Brands = append(output.Brands, &kpiBrandItem)

	}
	return output, nil
}
func UniqueBrand(brandSlice []entity.UniqueBrand) []entity.UniqueBrand {
	keys := make(map[entity.UniqueBrand]bool)
	uniqueList := []entity.UniqueBrand{}
	for _, entry := range brandSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func UniqueBrandData(brandSlice []entity.UniqueBrandData) []entity.UniqueBrandData {
	keys := make(map[entity.UniqueBrandData]bool)
	uniqueList := []entity.UniqueBrandData{}
	for _, entry := range brandSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func MapKpiVersionEntityToModel(inputEntity entity.KpiData) (model.GetKpiOffline, error) {
	var outputModel model.GetKpiOffline
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

	tempProductDesigns := make([]*model.KPIDesignRes, 0)
	tempBrandDesigns := make([]*model.KPIDesignRes, 0)

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
		tempProductDesigns = append(tempProductDesigns, &tempDesign)
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
		tempBrandDesigns = append(tempBrandDesigns, &tempDesign)
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

func MapKpiProductAnswerDataEntityToModel(kpiProductBrandAnswer []entity.KpiProductBrandAnswer) []*model.KpiProductBrandAnswer {

	var outputModel []*model.KpiProductBrandAnswer
	eventArrs := []entity.UniqueEvent{}
	for _, value := range kpiProductBrandAnswer {
		eventArr := entity.UniqueEvent{}
		eventArr.ScheduleEventId = value.ScheduleEventId.String
		eventArr.TeamMemberCustomer = value.TeamMemberCustomer.String
		eventArrs = append(eventArrs, eventArr)
	}
	uniqueEvent := UniqueEvent(eventArrs)
	for _, itemEvent := range uniqueEvent {
		eventModel := model.KpiProductBrandAnswer{}
		eventModel.EventID = itemEvent.ScheduleEventId
		eventModel.TeamCustomerID = itemEvent.TeamMemberCustomer

		// Brand Operation
		brandArrs := []entity.UniqueBrandData{}
		brandProductArrs := []entity.UniqueBrandProduct{}
		for _, valueBrand := range kpiProductBrandAnswer {
			if itemEvent.ScheduleEventId == valueBrand.ScheduleEventId.String {
				if valueBrand.KpiType.String == "brand" {
					brandArr := entity.UniqueBrandData{}

					brandArr.ItemId = valueBrand.ItemId.String
					brandArr.ItemName = valueBrand.ItemName.String
					brandArr.KpiVersionId = valueBrand.KpiVersionId.String

					brandArrs = append(brandArrs, brandArr)
				} else {
					brandProductArr := entity.UniqueBrandProduct{}

					// assign all data to local variable
					materialDescription := valueBrand.ItemDescription.String
					teamId := valueBrand.TeamID.String
					productId := valueBrand.ProductId.String
					teamProductId := valueBrand.TeamProductID.String
					principleName := valueBrand.ProductPrincipalName.String
					productKpiVersionId := valueBrand.KpiVersionId.String
					month := valueBrand.Month.String
					year := valueBrand.Year.String
					productPrincipleName := valueBrand.ProductPrincipalName.String
					brandId := valueBrand.ItemId.String
					productIsPriority := valueBrand.ProductIsPriority.Bool

					brandProductArr.MaterialDescription = materialDescription
					brandProductArr.TeamId = teamId
					brandProductArr.ProductId = productId
					brandProductArr.TeamProductID = teamProductId
					brandProductArr.PrincipleName = principleName
					brandProductArr.ProductKpiVersionId = productKpiVersionId
					brandProductArr.Month = month
					brandProductArr.Year = year
					brandProductArr.ProductPrincipalName = productPrincipleName
					// brandProductArr.ProductId = valueBrand.ProductId.String
					brandProductArr.BrandId = brandId
					brandProductArr.ProductIsPriority = productIsPriority
					brandProductArr.Category = int(valueBrand.Category.Int32)

					brandProductArrs = append(brandProductArrs, brandProductArr)
				}
			}
		}
		uniqueBrand := UniqueBrandData(brandArrs)
		uniqueBrandProduct := UniqueBrandProduct(brandProductArrs)

		var brandModel []*model.KpiBrandItemOffline
		for _, itemBrand := range uniqueBrand {
			kpiBrand := model.KpiBrandItemOffline{}
			var brandProducts []*model.KpiProductItemOffline
			var answerModels []*model.KpiAnswer

			itemBrandId := itemBrand.ItemId
			brandKpiVersionId := itemBrand.KpiVersionId
			itemBrandName := itemBrand.ItemName

			kpiBrand.BrandID = itemBrandId
			kpiBrand.BrandName = itemBrandName
			kpiBrand.BrandKpiVersionID = &brandKpiVersionId

			for _, valueBrandProduct := range uniqueBrandProduct {
				if valueBrandProduct.BrandId == itemBrand.ItemId {
					teamProductId := valueBrandProduct.TeamProductID
					principleName := valueBrandProduct.ProductPrincipalName
					ProductKpiVersionId := valueBrandProduct.ProductKpiVersionId
					isPriority := valueBrandProduct.ProductIsPriority
					productId := valueBrandProduct.ProductId
					teamId := valueBrandProduct.TeamId
					materialDes := valueBrandProduct.MaterialDescription

					kpiBrandProduct := model.KpiProductItemOffline{}

					kpiBrandProduct.TeamProductID = teamProductId
					kpiBrandProduct.PrincipalName = principleName
					kpiBrandProduct.ProductKpiVersionID = &ProductKpiVersionId
					kpiBrandProduct.IsPriority = isPriority
					kpiBrandProduct.ProductID = productId
					kpiBrandProduct.TeamID = teamId
					kpiBrandProduct.MaterialDescription = materialDes

					answerModels = MapKpiAnswerToProduct(kpiProductBrandAnswer, itemEvent.ScheduleEventId, itemBrand.ItemId, valueBrandProduct.TeamProductID)
					kpiBrandProduct.ProductKpiAnswer = answerModels

					brandProducts = append(brandProducts, &kpiBrandProduct)
				}
			}

			kpiBrand.Products = brandProducts

			answerModels = MapKpiAnswerToBrand(kpiProductBrandAnswer, itemEvent.ScheduleEventId, itemBrand.ItemId)
			kpiBrand.BrandKpiAnswer = answerModels

			brandModel = append(brandModel, &kpiBrand)
		}
		eventModel.Brands = append(eventModel.Brands, brandModel...)

		outputModel = append(outputModel, &eventModel)
	}

	return outputModel
}

func MapKpiAnswerToBrand(kpiProductBrandAnswer []entity.KpiProductBrandAnswer, scheduleEventId string, itemId string) []*model.KpiAnswer {
	var answerModels []*model.KpiAnswer
	for _, valueProductAnswer := range kpiProductBrandAnswer {
		if "brand" == valueProductAnswer.KpiType.String &&
			scheduleEventId == valueProductAnswer.ScheduleEventId.String &&
			itemId == valueProductAnswer.ItemId.String {
			KpiAnswerEntity := []entity.KPIAnswerStruct{}
			err := json.Unmarshal([]byte(valueProductAnswer.Answers.String), &KpiAnswerEntity)
			if err != nil {
				return answerModels
			}

			tempAnswers := make([]model.KPIAnswerRes, 0)
			for _, eachDesign := range KpiAnswerEntity {
				var tempAnswer model.KPIAnswerRes
				tempAnswer.QuestioNnumber = eachDesign.QuestioNnumber
				tempAnswer.Value = eachDesign.Value
				tempAnswers = append(tempAnswers, tempAnswer)
			}

			var answerModel model.KpiAnswer
			answerModel.ID = valueProductAnswer.KpiAnsId.String
			answerModel.Answer = tempAnswers
			answerModel.KpiID = valueProductAnswer.KpiId.String
			answerModel.KpiVersionID = valueProductAnswer.AnsKpiVersionId.String
			answerModel.Category = int(valueProductAnswer.Category.Int32)
			answerModel.TeamMemberCustomerID = valueProductAnswer.AnsTeamMemberCustomer.String
			answerModel.ScheduleEvent = valueProductAnswer.AnsScheduleEventId.String
			answerModel.TargetItem = valueProductAnswer.TargetItem.String
			answerModels = append(answerModels, &answerModel)
		}
	}

	return answerModels
}

func MapKpiAnswerToProduct(kpiProductBrandAnswer []entity.KpiProductBrandAnswer, scheduleEventId string, itemId string, tpId string) []*model.KpiAnswer {
	var answerModels []*model.KpiAnswer
	for _, valueProductAnswer := range kpiProductBrandAnswer {
		if scheduleEventId == valueProductAnswer.ScheduleEventId.String &&
			itemId == valueProductAnswer.ItemId.String &&
			tpId == valueProductAnswer.TeamProductID.String {
			KpiAnswerEntity := []entity.KPIAnswerStruct{}
			err := json.Unmarshal([]byte(valueProductAnswer.Answers.String), &KpiAnswerEntity)
			if err != nil {
				return answerModels
			}

			tempAnswers := make([]model.KPIAnswerRes, 0)
			for _, eachDesign := range KpiAnswerEntity {
				var tempAnswer model.KPIAnswerRes
				tempAnswer.QuestioNnumber = eachDesign.QuestioNnumber
				tempAnswer.Value = eachDesign.Value
				tempAnswers = append(tempAnswers, tempAnswer)
			}

			var answerModel model.KpiAnswer
			answerModel.ID = valueProductAnswer.KpiAnsId.String
			answerModel.Answer = tempAnswers
			answerModel.KpiID = valueProductAnswer.KpiId.String
			answerModel.KpiVersionID = valueProductAnswer.AnsKpiVersionId.String
			answerModel.Category = int(valueProductAnswer.Category.Int32)
			answerModel.TeamMemberCustomerID = valueProductAnswer.AnsTeamMemberCustomer.String
			answerModel.ScheduleEvent = valueProductAnswer.AnsScheduleEventId.String
			answerModel.TargetItem = valueProductAnswer.TargetItem.String
			answerModels = append(answerModels, &answerModel)
		}
	}
	return answerModels
}

// func UniqueKpiVersion(brandSlice []entity.UniqueKpiVersion) []entity.UniqueKpiVersion {
// 	keys := make(map[entity.UniqueKpiVersion]bool)
// 	uniqueList := []entity.UniqueKpiVersion{}
// 	for _, entry := range brandSlice {
// 		if _, value := keys[entry]; !value {
// 			keys[entry] = true
// 			uniqueList = append(uniqueList, entry)
// 		}
// 	}
// 	return uniqueList
// }

func UniqueEvent(brandSlice []entity.UniqueEvent) []entity.UniqueEvent {
	keys := make(map[entity.UniqueEvent]bool)
	uniqueList := []entity.UniqueEvent{}
	for _, entry := range brandSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func UniqueBrandProduct(brandSlice []entity.UniqueBrandProduct) []entity.UniqueBrandProduct {
	keys := make(map[entity.UniqueBrandProduct]bool)
	uniqueList := []entity.UniqueBrandProduct{}
	for _, entry := range brandSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}
