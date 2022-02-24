package mapper

import (
	"errors"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	uuid "github.com/gofrs/uuid"
)

func KpiModelToEntity(inputModel model.UpsertKpiRequest) (*entity.UpsertKpi, *model.ValidationResult) {
	var DataEntity entity.UpsertKpi
	var result model.ValidationResult
	var id *uuid.UUID

	kpiBrandDesigns := []entity.UpsertKpiDesign{}
	kpiProductDesigns := []entity.UpsertKpiDesign{}

	for _, eachDesign := range inputModel.BrandDesign {
		kpiQuestions := []entity.KpiQuestion{}
		kpiDesign := entity.UpsertKpiDesign{}

		if eachDesign.Questions != nil {
			for _, eachQuestion := range eachDesign.Questions {
				KpiQuestion := entity.KpiQuestion{}

				if eachQuestion.Active != nil {
					KpiQuestion.Active = *eachQuestion.Active
				}
				if eachQuestion.OptionValues != nil {
					KpiQuestion.OptionValues = eachQuestion.OptionValues
				}
				if eachQuestion.Required != nil {
					KpiQuestion.Required = *eachQuestion.Required
				} else {
					KpiQuestion.Required = true
				}

				KpiQuestion.QuestionNumber = eachQuestion.QuestionNumber
				KpiQuestion.Title = eachQuestion.Title
				KpiQuestion.Type = eachQuestion.Type

				kpiQuestions = append(kpiQuestions, KpiQuestion)
			}

			kpiDesign.Questions = kpiQuestions

			if postgres.CheckCodeIDForKpi(eachDesign.Category, "KPICategory") {
				kpiDesign.Category = eachDesign.Category
			} else {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid category id for design"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
				return &entity.UpsertKpi{}, &result
			}
		}

		if eachDesign.Active {
			kpiDesign.Active = true
		} else {
			kpiDesign.Active = false
		}

		if eachDesign.EffectiveStartDate != "" {
			kpiDesign.EffectiveStartDate = eachDesign.EffectiveStartDate
		} else {
			kpiDesign.EffectiveStartDate = ""
		}

		if eachDesign.EffectiveEndDate != "" {
			kpiDesign.EffectiveEndDate = eachDesign.EffectiveEndDate
		} else {
			kpiDesign.EffectiveEndDate = ""
		}

		kpiDesign.Name = *eachDesign.Name
		kpiDesign.Type = *eachDesign.Type
		kpiBrandDesigns = append(kpiBrandDesigns, kpiDesign)
	}
	DataEntity.BrandDesign = kpiBrandDesigns

	for _, eachDesign := range inputModel.ProductDesign {
		kpiQuestions := []entity.KpiQuestion{}
		kpiDesign := entity.UpsertKpiDesign{}

		if eachDesign.Questions != nil {
			for _, eachQuestion := range eachDesign.Questions {
				KpiQuestion := entity.KpiQuestion{}

				if eachQuestion.Active != nil {
					KpiQuestion.Active = *eachQuestion.Active
				}
				if eachQuestion.OptionValues != nil {
					KpiQuestion.OptionValues = eachQuestion.OptionValues
				}
				if eachQuestion.Required != nil {
					KpiQuestion.Required = *eachQuestion.Required
				} else {
					KpiQuestion.Required = true
				}

				KpiQuestion.QuestionNumber = eachQuestion.QuestionNumber
				KpiQuestion.Title = eachQuestion.Title
				KpiQuestion.Type = eachQuestion.Type

				kpiQuestions = append(kpiQuestions, KpiQuestion)
			}

			kpiDesign.Questions = kpiQuestions

			if postgres.CheckCodeIDForKpi(eachDesign.Category, "KPICategory") {
				kpiDesign.Category = eachDesign.Category
			} else {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid category id for design"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
				return &entity.UpsertKpi{}, &result
			}
		}

		if eachDesign.Active {
			kpiDesign.Active = true
		} else {
			kpiDesign.Active = false
		}

		if eachDesign.EffectiveStartDate != "" {
			kpiDesign.EffectiveStartDate = eachDesign.EffectiveStartDate
		} else {
			kpiDesign.EffectiveStartDate = ""
		}

		if eachDesign.EffectiveEndDate != "" {
			kpiDesign.EffectiveEndDate = eachDesign.EffectiveEndDate
		} else {
			kpiDesign.EffectiveEndDate = ""
		}

		kpiDesign.Name = *eachDesign.Name
		kpiDesign.Type = *eachDesign.Type
		kpiProductDesigns = append(kpiProductDesigns, kpiDesign)
	}
	DataEntity.ProductDesign = kpiProductDesigns

	if inputModel.ID != nil && *inputModel.ID != "" {
		uuid, err := uuid.FromString(*inputModel.ID)
		if err == nil {
			id = &uuid
		} else {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Kpi ID format is invalid!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}

		DataEntity.ID = id
		productId := postgres.GetIDProductBrandIdFromKpiParentId(*inputModel.ID, "product")
		DataEntity.ProductId = productId
		brandId := postgres.GetIDProductBrandIdFromKpiParentId(*inputModel.ID, "brand")
		DataEntity.BrandId = brandId

		if !postgres.IsMonthyearValid(inputModel.IsPriority, inputModel.EffectiveMonth, inputModel.EffectiveYear, *inputModel.ID) {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Please provide correct month, year, isPriority"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}

	DataEntity.TargetProduct = inputModel.TargetProducts
	if inputModel.ID == nil || *inputModel.ID == "" {
		kpiExist, err := postgres.KpiCombinationExist(inputModel.TargetTeam, inputModel.EffectiveMonth, inputModel.EffectiveYear)
		if err == nil {
			if kpiExist {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "KPI for these team for provided month already exist"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
		}

		for _, BrandId := range inputModel.TargetBrand {
			if !postgres.ProductBelongInAnyBrand(BrandId, DataEntity.TargetProduct) {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Select at least one product from each selected brands!"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
		}
	}

	DataEntity.TargetBrand = inputModel.TargetBrand

	if inputModel.IsDeleted && DataEntity.ID != nil {
		DataEntity.IsDeleted = true
	} else if inputModel.IsDeleted && DataEntity.ID == nil {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Provide KPI parent Id for deleting KPI"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	} else {
		DataEntity.IsDeleted = false
	}

	if DataEntity.IsDeleted {
		existingID, err := postgres.GetIDForSameTeamAndMonth(inputModel.TargetTeam, inputModel.EffectiveMonth, inputModel.EffectiveYear)
		if err != nil {
			DataEntity.ExistingID = ""
		} else {
			DataEntity.ExistingID = existingID
		}
	}
	DataEntity.EffectiveMonth = inputModel.EffectiveMonth
	DataEntity.EffectiveYear = &inputModel.EffectiveYear
	if inputModel.IsPriority {
		DataEntity.IsPriority = true
	} else {
		DataEntity.IsPriority = false
	}
	if inputModel.Name != "" {
		DataEntity.Name = inputModel.Name
	}
	DataEntity.TypeName = inputModel.Type
	if inputModel.TargetTeam != "" {
		DataEntity.TargetTeam = inputModel.TargetTeam
	}
	return &DataEntity, &result
}

func ValidateGetKpisInput(input *model.GetKpiInput, salesOrgId string) (entity.GetKpisInput, error) {
	inputEntity := entity.GetKpisInput{
		SalesOrgId: &salesOrgId,
	}

	if input.ParentKpiID != nil && *input.ParentKpiID != "" {
		_, err := uuid.FromString(*input.ParentKpiID)
		if err != nil {
			return inputEntity, errors.New("Invalid Parent KPI Id!")
		} else {
			if !postgres.HasKpi(*input.ParentKpiID) {
				return inputEntity, errors.New("Parent KPI does not exist!")
			} else {
				inputEntity.ParentKpiID = input.ParentKpiID
			}
		}
	}

	if input.KpiID != nil && *input.KpiID != "" {
		_, err := uuid.FromString(*input.KpiID)
		if err != nil {
			return inputEntity, errors.New("Invalid KPI Id!")
		} else {
			if !postgres.HasKpi(*input.KpiID) {
				return inputEntity, errors.New("KPI does not exist!")
			} else {
				inputEntity.KpiID = input.KpiID
			}
		}
	}

	if input.KpiVersionID != nil && *input.KpiVersionID != "" {
		_, err := uuid.FromString(*input.KpiVersionID)
		if err != nil {
			return inputEntity, errors.New("Invalid KPI Version Id!")
		} else {
			if !postgres.HasKPIVersion(*input.KpiVersionID) {
				return inputEntity, errors.New("KPI Version does not exist!")
			} else {
				inputEntity.KpiVersionID = input.KpiVersionID
			}
		}
	}

	if input.TeamID != nil && *input.TeamID != "" {
		_, err := uuid.FromString(*input.TeamID)
		if err != nil {
			return inputEntity, errors.New("Invalid Team Id!")
		} else {
			if !postgres.HasTeam(*input.TeamID) {
				return inputEntity, errors.New("Team does not exist!")
			} else {
				inputEntity.TeamID = input.TeamID
			}
		}
	}

	if input.BrandID != nil && *input.BrandID != "" {
		_, err := uuid.FromString(*input.BrandID)
		if err != nil {
			return inputEntity, errors.New("Invalid Brand Id!")
		} else {
			if !postgres.HasBrand(*input.BrandID) {
				return inputEntity, errors.New("Brand does not exist!")
			} else {
				inputEntity.BrandID = input.BrandID
			}
		}
	}

	if input.TeamProductID != nil && *input.TeamProductID != "" {
		_, err := uuid.FromString(*input.TeamProductID)
		if err != nil {
			return inputEntity, errors.New("Invalid Team Product Id!")
		} else {
			if !postgres.HasTeamProduct(*input.TeamProductID) {
				return inputEntity, errors.New("Team Product does not exist!")
			} else {
				inputEntity.TeamProductID = input.TeamProductID
			}
		}
	}

	if input.Year != nil && *input.Year != 0 {
		if *input.Year < 1900 || *input.Year > 3000 {
			return inputEntity, errors.New("Invalid Year!")
		} else {
			inputEntity.Year = input.Year
		}
	}

	if input.Month != nil && *input.Month != 0 {
		if *input.Month < 1 || *input.Month > 12 {
			return inputEntity, errors.New("Invalid Month!")
		} else {
			inputEntity.Month = input.Month
		}
	}

	if input.SearchItem != nil && *input.SearchItem != "" {
		inputEntity.SearchItem = input.SearchItem
	}

	if input.Limit != nil && *input.Limit > 0 {
		inputEntity.Limit = input.Limit
		var initalOffset = 0
		inputEntity.Offset = &initalOffset
		if input.PageNo != nil && *input.PageNo > 0 {
			inputEntity.PageNo = input.PageNo
			*inputEntity.Offset = *input.Limit * (*input.PageNo - 1)
		}
	}

	return inputEntity, nil
}
