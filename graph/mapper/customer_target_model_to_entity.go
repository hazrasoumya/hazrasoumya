package mapper

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	uuid "github.com/gofrs/uuid"
)

func CustomerModelToEntity(inputModel model.CustomerTargetInput, salesOrg, loggedInId string) (*entity.CustomerTarget, *model.ValidationResult) {
	result := &model.ValidationResult{Error: false}
	entityData := entity.CustomerTarget{}
	entityData.SalesOrgId = salesOrg
	var id *uuid.UUID

	if inputModel.ID != nil && *inputModel.ID != "" {
		uuid, err := uuid.FromString(*inputModel.ID)
		if err == nil {
			id = &uuid
			if !postgres.IdPresent(*inputModel.ID, salesOrg) {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Please Provide Valid ID"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
		} else {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "ID format is invalid"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}

	}

	CustomerMonth := []int64{}

	if inputModel.Answers != nil {
		for _, value := range inputModel.Answers {
			CustomerMonth = append(CustomerMonth, value.Month)
		}
	}

	uniqueMonth := uniqueKpiMonth(CustomerMonth)
	monthLen := len(uniqueMonth)

	if monthLen != 12 {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid Customer Target Month(s)!"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	} else {
		for _, eachUniqueMonth := range uniqueMonth {
			if eachUniqueMonth > 12 || eachUniqueMonth < 1 {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Inappropriate Customer Target Month(s) tag!"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
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
		entityData.Year = inputModel.Year
	}

	if !postgres.ValidBrandId(inputModel.ProductBrandID, salesOrg) && !postgres.ValidProductId(inputModel.ProductBrandID, salesOrg) {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Please provide valid product-brand-id"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	} else {
		entityData.ProductBraindId = inputModel.ProductBrandID
	}

	if !postgres.ValidBrandProductValue(inputModel.Type, "TargetCustomerType") {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Please Provide Valid Type"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)

	} else {
		TypeId := postgres.GetValueAndCategory(inputModel.Type, "TargetCustomerType")
		entityData.TypeID = TypeId
	}

	if !postgres.ValidBrandProductValue(inputModel.Category, "MerchandisingType") {

		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Please Provide Valid Category"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)

	} else {
		if (inputModel.Category == "brand_share_shelf" || inputModel.Category == "must_have_sku_comp") && !(inputModel.Type == "targetbrand") {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Brand Share shelf and Must Have SKU % Compliance only applicable for Brand"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		} else {
			entityData.Category = postgres.GetValueAndCategory(inputModel.Category, "MerchandisingType")
		}
	}
	if inputModel.ID == nil || *inputModel.ID == "" {
		if postgres.AlreadyExists(&entityData, entityData.SalesOrgId) {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "This Data Set already Exists"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}

	if inputModel.ID != nil && *inputModel.ID != "" {
		if !postgres.IsCombinationExist(&entityData, entityData.SalesOrgId, inputModel.ID) {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Wrong customer target data combination!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}

	MonthlyData := []entity.TargetValue{}

	for _, targetData := range inputModel.Answers {
		var value entity.TargetValue
		value.Month = targetData.Month
		value.Value = targetData.Value
		MonthlyData = append(MonthlyData, value)
	}

	if inputModel.IsDeleted != nil {
		entityData.IsDeleted = *inputModel.IsDeleted
	}

	entityData.ID = id
	entityData.Targets = MonthlyData
	entityData.SalesOrgId = salesOrg

	return &entityData, result
}

func GetBrandShareID() {
	panic("unimplemented")
}

func ValidateTargetCustomerInput(input model.GetTargetCustomerRequest, salesOrgId string) (entity.TargetCustomerInput, error) {
	inputEntity := entity.TargetCustomerInput{
		SalesOrgId: &salesOrgId,
	}

	if input.ID != nil && *input.ID != "" {
		_, err := uuid.FromString(*input.ID)
		if err != nil {
			return inputEntity, errors.New("Invalid Customer Target Id!")
		} else {
			_, err := postgres.HasTargetCustomerId(*input.ID)

			if err != nil {
				return inputEntity, errors.New("Customer Target Id does not exist!")
			} else {
				inputEntity.TargetCustomerId = input.ID
			}
		}
	}

	if input.Type != nil && *input.Type != "" {
		if !strings.EqualFold(*input.Type, "targetbrand") && !strings.EqualFold(*input.Type, "targetproduct") {
			return inputEntity, errors.New("Invalid Customer Target Type!")
		} else {
			inputEntity.Type = input.Type

			if input.ProductBrandID != nil && *input.ProductBrandID != "" {
				_, err := uuid.FromString(*input.ProductBrandID)
				if err != nil {
					return inputEntity, errors.New("Invalid Product/Brand Id!")
				} else {
					isActive, err := postgres.HasProductBrandId(*input.ProductBrandID, *input.Type)

					if err != nil {
						return inputEntity, errors.New("Product/Brand does not exist!")
					} else {
						if isActive {
							inputEntity.ProductBrandId = input.ProductBrandID
						} else {
							if strings.EqualFold(*inputEntity.Type, "targetproduct") {
								return inputEntity, errors.New("Product does not exist!")
							} else {
								return inputEntity, errors.New("Brand does not exist!")
							}
						}
					}
				}
			}
		}
	} else {
		if input.ProductBrandID != nil && *input.ProductBrandID != "" {
			_, err := uuid.FromString(*input.ProductBrandID)
			if err != nil {
				return inputEntity, errors.New("Invalid Product/Brand Id!")
			} else {
				isActive, err := postgres.HasProductBrandId(*input.ProductBrandID, "targetproduct")

				if err != nil {
					isActive, err := postgres.HasProductBrandId(*input.ProductBrandID, "targetbrand")

					if err != nil {
						return inputEntity, errors.New("Product/Brand does not exist!")
					} else {
						if isActive {
							typeString := "targetbrand"
							inputEntity.Type = &typeString
							inputEntity.ProductBrandId = input.ProductBrandID
						} else {
							return inputEntity, errors.New("The Brand does not exist!")
						}
					}
				} else {
					if isActive {
						typeString := "targetproduct"
						inputEntity.Type = &typeString
						inputEntity.ProductBrandId = input.ProductBrandID
					} else {
						return inputEntity, errors.New("The Product does not exist!")
					}
				}
			}
		} else {
			typeString := "targetproduct"
			inputEntity.Type = &typeString
		}
	}

	if input.Category != nil && !strings.EqualFold(*input.Category, "") {
		if !postgres.IsValidCategory(*input.Category) {
			return inputEntity, errors.New("Invalid Customer Target Category!")
		} else {
			inputEntity.Category = input.Category
		}
	}

	if input.ProductBrandName != nil && *input.ProductBrandName != "" {
		inputEntity.ProductBrandName = input.ProductBrandName
	}

	if input.Year != nil && *input.Year != 0 {
		if *input.Year < 1900 || *input.Year > 3000 {
			return inputEntity, errors.New("Invalid Year!")
		} else {
			inputEntity.Year = input.Year
		}
	} else {
		yearInt := time.Now().Year()
		inputEntity.Year = &yearInt
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

func MapCustomerTargetToEntity(input []entity.TargetCustomerResponse, year int) []entity.CustomerTargetExcelInterface {
	outputEntity := []entity.CustomerTargetExcelInterface{}
	for _, data := range input {
		innerEntity := entity.CustomerTargetExcelInterface{}
		var rowData []interface{}
		rowData = append(rowData, data.CustomerTargetId.String)
		rowData = append(rowData, data.Type.String)
		rowData = append(rowData, data.Category.String)
		if strings.EqualFold(data.Type.String, "targetproduct") {
			full := strings.Split(data.ProductBrandName.String, "|")
			principal := full[0]
			material := full[1]
			rowData = append(rowData, principal)
			rowData = append(rowData, material)
			rowData = append(rowData, "")
		} else {
			rowData = append(rowData, "")
			rowData = append(rowData, "")
			rowData = append(rowData, data.ProductBrandName.String)
		}
		rowData = append(rowData, year)
		var target []*model.Target
		if len(data.Targets.String) > 0 {
			err := json.Unmarshal([]byte(data.Targets.String), &target)
			if err == nil {
				for _, value := range target {
					rowData = append(rowData, *value.Value)
				}
			}
		} else {
			for i := 0; i < 12; i++ {
				rowData = append(rowData, 0)
			}
		}

		innerEntity.Data = rowData
		outputEntity = append(outputEntity, innerEntity)
	}
	return outputEntity

}

func GetCustomerGroupByIndustrialCode(dbResponse []entity.CustomerGroupResponse) []*model.CustomerGroup {
	industryCode := []entity.UniqueIndustryCode{}
	for _, value := range dbResponse {
		industry := entity.UniqueIndustryCode{}
		industry.IndustryCode = value.IndustrialCode
		industryCode = append(industryCode, industry)
	}
	uniqueCodeList := UniqueIndustryCode(industryCode)
	finalData := []*model.CustomerGroup{}
	for _, value := range uniqueCodeList {
		result := &model.CustomerGroup{}
		result.InDusTrialCode = value.IndustryCode.String
		customerData := []*model.CustomerResponse{}
		for _, data := range dbResponse {
			if value.IndustryCode.String == data.IndustrialCode.String {
				customer := model.CustomerResponse{}
				customer.CustoMerID = data.CustomerId.String
				customer.CustoMerName = data.CustomerName.String
				customer.SoldTo = data.SoldTo.String
				customer.ShipTo = data.ShipTo.String
				customerData = append(customerData, &customer)
			}
		}
		result.CustomeDetails = customerData
		finalData = append(finalData, result)
	}
	return finalData
}

func UniqueIndustryCode(data []entity.UniqueIndustryCode) []entity.UniqueIndustryCode {
	keys := make(map[entity.UniqueIndustryCode]bool)
	uniqueCodeList := []entity.UniqueIndustryCode{}
	for _, entry := range data {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueCodeList = append(uniqueCodeList, entry)
		}
	}
	return uniqueCodeList
}

func CheckLineManager(teamId []entity.CustomerLineManager, input *model.CustomerGroupInput) ([]string, bool) {
	var teams []string
	var flag bool
	if input != nil && input.TeamID != nil && *input.TeamID != "" {
		for _, value := range teamId {
			if value.TeamId.String == *input.TeamID {
				teams = append(teams, value.TeamId.String)
				flag = true
			}
		}
	} else {
		for _, value := range teamId {
			teams = append(teams, value.TeamId.String)
		}
		flag = true
	}
	return teams, flag
}
