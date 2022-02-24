package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	excelizev2 "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gofrs/uuid"

	"github.com/eztrade/kpi/graph/azure"
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/mapper"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	"github.com/eztrade/login/auth"
)

func UpsertCustomerTarget(input model.CustomerTargetInput, userEntity *entity.LoggedInUser) *model.KpiResponse {
	entity, validationResult := mapper.CustomerModelToEntity(input, userEntity.SalesOrganisaton, userEntity.ID)
	kpiResponse := &model.KpiResponse{}
	if !validationResult.Error {
		postgres.CustomerTarget(entity, kpiResponse, userEntity.ID)
	} else {
		kpiResponse.Error = true
		kpiResponse.ErrorCode = 999
		kpiResponse.Message = "Upsert customer target validation failed"
		kpiResponse.ValidationErrors = validationResult.ValidationMessage
	}
	return kpiResponse
}

func GetTargetCustomer(input model.GetTargetCustomerRequest, ctx *context.Context) (*model.GetTargetCustomerResponse, error) {
	userAD := auth.GetADName(*ctx)
	salesOrgID := auth.GetSalesOrg(*ctx)
	response := model.GetTargetCustomerResponse{}

	dbInput, err := mapper.ValidateTargetCustomerInput(input, *salesOrgID)
	if err != nil {
		return &response, err
	}

	dbResponse, totalPages, err := postgres.GetTargetCustomer(dbInput)
	if err != nil {
		return &response, err
	}

	if !input.IsExcel {
		for _, dbData := range dbResponse {
			outputModel := model.TargetCustomer{}

			customerTargetId := dbData.CustomerTargetId.String
			outputModel.CustomerTargetID = &customerTargetId
			typeString := dbData.Type.String
			outputModel.Type = &typeString
			category := dbData.Category.String
			outputModel.Category = &category
			productBrandId := dbData.ProductBrandId.String
			outputModel.ProductBrandID = &productBrandId
			productBrandName := dbData.ProductBrandName.String
			outputModel.ProductBrandName = &productBrandName
			outputModel.Year = dbInput.Year

			var target []*model.Target
			if len(dbData.Targets.String) > 0 {
				err := json.Unmarshal([]byte(dbData.Targets.String), &target)
				if err != nil {
					return &response, err
				}
			} else {
				value := 0.0
				for i := 1; i <= 12; i++ {
					var tar model.Target
					month := i
					tar.Month = &month
					tar.Value = &value
					target = append(target, &tar)
				}

			}
			outputModel.Targets = append(outputModel.Targets, target...)

			response.Data = append(response.Data, &outputModel)
		}

		response.TotalPage = totalPages

		return &response, nil
	}
	excelentites := mapper.MapCustomerTargetToEntity(dbResponse, *dbInput.Year)
	if len(excelentites) == 0 {
		log.Println("No data found : ", excelentites)
		return &model.GetTargetCustomerResponse{}, errors.New("No data found")
	}

	sheetName := "customer_target"
	f := excelizev2.NewFile()
	f.SetSheetName("Sheet1", sheetName)

	streamWriter, _ := f.NewStreamWriter(sheetName)
	// Populate excel header columns

	var headerRow []interface{}
	headerRow = append(headerRow, ("Id"))
	headerRow = append(headerRow, ("Type"))
	headerRow = append(headerRow, ("Category"))
	headerRow = append(headerRow, ("Principle Name"))
	headerRow = append(headerRow, ("Material Description"))
	headerRow = append(headerRow, ("Brand Name"))
	headerRow = append(headerRow, ("Year"))
	headerRow = append(headerRow, ("January"))
	headerRow = append(headerRow, ("February"))
	headerRow = append(headerRow, ("March"))
	headerRow = append(headerRow, ("April"))
	headerRow = append(headerRow, ("May"))
	headerRow = append(headerRow, ("June"))
	headerRow = append(headerRow, ("July"))
	headerRow = append(headerRow, ("August"))
	headerRow = append(headerRow, ("September"))
	headerRow = append(headerRow, ("October"))
	headerRow = append(headerRow, ("November"))
	headerRow = append(headerRow, ("December"))

	cell, _ := excelizev2.CoordinatesToCellName(1, 2)
	if err := streamWriter.SetRow(cell, headerRow); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}

	for index, cstarget := range excelentites {
		rowNumber := index + 3
		cell, _ := excelizev2.CoordinatesToCellName(1, rowNumber)
		if err := streamWriter.SetRow(cell, cstarget.Data); err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
		}
	}
	if err := streamWriter.Flush(); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	ExportCustomerTargetnstructionExcel(f)
	filename := *userAD + time.Now().Format("20060102150405") + "customer_targetExcel.xlsx"
	response = model.GetTargetCustomerResponse{}

	// //--------------------------------For saving the file locally----------------------------------------
	// if err := f.SaveAs(filename); err != nil {
	// 	println(err.Error())
	// }
	// response.URL = ""
	//----------------------------------------------------------------------------------------------------

	blobURL, err := azure.UploadBytesToBlob(getBytesFromFileV2(f), filename)
	if err != nil {
		return &model.GetTargetCustomerResponse{}, err
	}
	response.URL = blobURL

	return &response, nil
}

func ExportCustomerTargetnstructionExcel(f *excelizev2.File) {
	instructions := [11]string{"Avoid editing the column header and the sheet name", "Please do not edit id column generated by the system. Leave the id column blank for new records",
		"To delete records, keep the id intact and delete the other column values", "Please avoid data formatting across rows where there is no data",
		"Please avoid to only delete content on the excel, ensure you remove the formatting when data is removed from rows/columns too",
		"Please provide correct type, category", "Please provide correct Principle name, Material description, Brand name",
		"If target is targetbrand then only provide brand name and keep material description, priciple name blank", "If target is targetproduct then provide both principle name and material description and leave brand name blank", ""}

	sheetName := "Instruction"
	f.NewSheet(sheetName)

	streamWriter, err := f.NewStreamWriter(sheetName)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}

	for fieldIndex := 0; fieldIndex < len(instructions); fieldIndex++ {
		var headerRow []interface{}
		headerRow = append(headerRow, instructions[fieldIndex])
		cell, _ := excelizev2.CoordinatesToCellName(1, fieldIndex+1)
		if err := streamWriter.SetRow(cell, headerRow); err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
		}
	}
	// Populate header row

	rowNumber := 11

	var rowData1 []interface{}
	rowData1 = append(rowData1, "")
	rowData1 = append(rowData1, "Mandatory")
	rowData1 = append(rowData1, "Limit")
	cell1, _ := excelizev2.CoordinatesToCellName(1, rowNumber)
	if err := streamWriter.SetRow(cell1, rowData1); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	rowNumber = rowNumber + 1

	var rowData2 []interface{}
	rowData2 = append(rowData2, "id")
	rowData2 = append(rowData2, "N")
	rowData2 = append(rowData2, "-")
	cell2, _ := excelizev2.CoordinatesToCellName(1, rowNumber)
	if err := streamWriter.SetRow(cell2, rowData2); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	rowNumber = rowNumber + 1

	var rowData3 []interface{}
	rowData3 = append(rowData3, "Type")
	rowData3 = append(rowData3, "y")
	rowData3 = append(rowData3, "targetproduct, targetbrand")
	cell3, _ := excelizev2.CoordinatesToCellName(1, rowNumber)
	if err := streamWriter.SetRow(cell3, rowData3); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	rowNumber = rowNumber + 1

	var rowData4 []interface{}
	rowData4 = append(rowData4, "Category")
	rowData4 = append(rowData4, "Y")
	rowData4 = append(rowData4, "posmexecution, promotionexecution,  distribution, brand_share_shelf, must_have_sku_comp")
	cell4, _ := excelizev2.CoordinatesToCellName(1, rowNumber)
	if err := streamWriter.SetRow(cell4, rowData4); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	rowNumber = rowNumber + 1

	var rowData5 []interface{}
	rowData5 = append(rowData5, "Principal Name")
	rowData5 = append(rowData5, "N")
	rowData5 = append(rowData5, "N")
	cell5, _ := excelizev2.CoordinatesToCellName(1, rowNumber)
	if err := streamWriter.SetRow(cell5, rowData5); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	rowNumber = rowNumber + 1

	var rowData6 []interface{}
	rowData6 = append(rowData6, "Material Description")
	rowData6 = append(rowData6, "N")
	rowData6 = append(rowData6, "N")
	cell6, _ := excelizev2.CoordinatesToCellName(1, rowNumber)
	if err := streamWriter.SetRow(cell6, rowData6); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	rowNumber = rowNumber + 1

	var rowData7 []interface{}
	rowData7 = append(rowData7, "Brand Name")
	rowData7 = append(rowData7, "N")
	rowData7 = append(rowData7, "N")
	cell7, _ := excelizev2.CoordinatesToCellName(1, rowNumber)
	if err := streamWriter.SetRow(cell7, rowData7); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	rowNumber = rowNumber + 1

	var rowData8 []interface{}
	rowData8 = append(rowData8, "Year")
	rowData8 = append(rowData8, "Y")
	rowData8 = append(rowData8, "1900-3000")
	cell8, _ := excelizev2.CoordinatesToCellName(1, rowNumber)
	if err := streamWriter.SetRow(cell8, rowData8); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	rowNumber = rowNumber + 1
	var rowData9 []interface{}
	rowData9 = append(rowData9, "Customer")
	rowData9 = append(rowData9, "Y")
	rowData9 = append(rowData9, "Only integer")
	cell9, _ := excelizev2.CoordinatesToCellName(1, rowNumber)
	if err := streamWriter.SetRow(cell9, rowData9); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	rowNumber = rowNumber + 1
	var rowData10 []interface{}
	rowData10 = append(rowData10, "* brand_share_shelf and must_have_sku_comp applicable only for brand")
	cell10, _ := excelizev2.CoordinatesToCellName(1, rowNumber)
	if err := streamWriter.SetRow(cell10, rowData10); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}

	if err := streamWriter.Flush(); err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
}

func getBytesFromFileV2(f *excelizev2.File) []byte {
	buffer, err := f.WriteToBuffer()
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	myresult := buffer.Bytes()
	return myresult
}

func GetCustomerGroup(input *model.CustomerGroupInput, userEntity *entity.LoggedInUser) (*model.CustomerGroupResponse, error) {
	response := model.CustomerGroupResponse{}
	var data []string
	var flag bool
	if input != nil {
		if input.TeamID != nil && *input.TeamID != "" {
			_, err := uuid.FromString(*input.TeamID)
			if err != nil {
				return nil, errors.New("Invalid TeamId")
			}
		}
	}
	if userEntity.AuthRole == "sfe" || userEntity.AuthRole == "cbm" {
		if input != nil {
			if input.TeamID != nil && *input.TeamID != "" {
				data = append(data, *input.TeamID)
			}
		}
	} else {
		teamId, err := postgres.GetLineOneManager(userEntity.ID)
		if err != nil {
			return &response, nil
		}
		data, flag = mapper.CheckLineManager(teamId, input)
		if !flag {
			return nil, errors.New("You are not the line1manager for this input team")
		}
	}
	dbResponse, err := postgres.GetCustomerGroup(userEntity, data, input)
	if err != nil {
		return &response, nil
	}
	result := mapper.GetCustomerGroupByIndustrialCode(dbResponse)
	response.CustoMerData = result
	return &response, nil
}
