package controller

import (
	"archive/zip"
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/eztrade/kpi/graph/azure"
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/mapper"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	"github.com/eztrade/kpi/graph/util"
)

func RetrievePictures(input model.PicturesInput, loggedInUserEntity *entity.LoggedInUser) (model.RetrievePicturesResponse, error) {
	response := model.RetrievePicturesResponse{}

	if input.Type != "" {
		if input.Type != "Brand" && input.Type != "Product" && input.Type != "Customer" {
			response.Error = true
			response.Message = "Provied correct type!"
			return response, nil
		}
	}

	var StartDate string
	var EndDate string
	monthFilter := true
	yearFilter := true
	var sYear, sMonth, sDay int
	var eYear, eMonth, eDay int

	timeObjsMap := make(map[string]time.Time)
	if input.StartDate != "" {
		startDateObj, startDateString, startDateValidateErr := util.IsValidDateWithUSDateString(input.StartDate)
		if startDateValidateErr != nil {
			response.Error = true
			response.Message = "Wrong Start Date Format!"
			return response, nil
		}
		timeObjsMap["start"] = startDateObj
		StartDate = startDateString
		startYear := strings.Trim(startDateString[0:4], " \t\n\r")
		startMonth := strings.Trim(startDateString[5:7], " \t\n\r")
		startDay := strings.Trim(startDateString[8:10], " \t\n\r")

		sYear, _ = strconv.Atoi(startYear)
		sMonth, _ = strconv.Atoi(startMonth)
		sDay, _ = strconv.Atoi(startDay)
	}

	if input.EndDate != "" {
		startDateObj, EndDateString, startDateValidateErr := util.IsValidDateWithUSDateString(input.EndDate)
		if startDateValidateErr != nil {
			response.Error = true
			response.Message = "Wrong End Date Format!"
			return response, nil
		}
		timeObjsMap["End"] = startDateObj
		EndDate = EndDateString
		endYear := strings.Trim(EndDateString[0:4], " \t\n\r")
		endMonth := strings.Trim(EndDateString[5:7], " \t\n\r")
		endDay := strings.Trim(EndDateString[8:10], " \t\n\r")

		eYear, _ = strconv.Atoi(endYear)
		eMonth, _ = strconv.Atoi(endMonth)
		eDay, _ = strconv.Atoi(endDay)
	}

	var productBrandCustomerId []string

	if input.Type == "Product" && input.ProductID != nil {
		for _, value := range input.ProductID {
			productBrandCustomerId = append(productBrandCustomerId, *value)
		}
		if input.Type == "Product" && (input.BrandID != nil || input.CustomerID != nil) {
			response.Error = true
			response.Message = "You selected product, please provide only product data!"
			return response, nil
		}
	} else {
		if input.Type == "Product" && (input.BrandID != nil || input.CustomerID != nil) {
			response.Error = true
			response.Message = "You selected product, please provide only product data!"
			return response, nil
		}
	}

	if input.Type == "Brand" && input.BrandID != nil {
		for _, value := range input.BrandID {
			productBrandCustomerId = append(productBrandCustomerId, *value)
		}
		if input.Type == "Brand" && (input.ProductID != nil || input.CustomerID != nil) {
			response.Error = true
			response.Message = "You selected brand, please provide only brand data!"
			return response, nil
		}
	} else {
		if input.Type == "Brand" && (input.ProductID != nil || input.CustomerID != nil) {
			response.Error = true
			response.Message = "You selected brand, please provide only brand data!"
			return response, nil
		}
	}

	if input.Type == "Customer" && input.CustomerID != nil {
		for _, value := range input.CustomerID {
			productBrandCustomerId = append(productBrandCustomerId, *value)
		}
		if input.Type == "Customer" && (input.BrandID != nil || input.ProductID != nil) {
			response.Error = true
			response.Message = "You selected customer, please provide only customer data!"
			return response, nil
		}
	} else {
		if input.Type == "Customer" && (input.BrandID != nil || input.ProductID != nil) {
			response.Error = true
			response.Message = "You selected customer, please provide only customer data!"
			return response, nil
		}
	}

	t1 := util.DateDifference(sYear, sMonth, sDay)
	t2 := util.DateDifference(eYear, eMonth, eDay)
	days := t2.Sub(t1).Hours() / 24

	if days > 92 {
		response.Error = true
		response.Message = "Date range can not be greater than 3 months!"
		return response, nil
	}

	pictureZip := model.PictureZipInput{}
	getTeamIDS, _ := postgres.GetTeamIDsBySalesOrg(loggedInUserEntity)
	getPictures, err := postgres.RetrievePicturesInfo(input, pictureZip, false, loggedInUserEntity, getTeamIDS, monthFilter, StartDate, yearFilter, EndDate, productBrandCustomerId)
	if err != nil {
		return model.RetrievePicturesResponse{}, err
	}

	teamArrs := []entity.UniqueTeam{}
	customerArrs := []entity.UniqueCustomer{}
	for _, item := range getPictures {
		teamArr := entity.UniqueTeam{}
		teamArr.TeamID = item.Team.String
		teamArrs = append(teamArrs, teamArr)
	}
	for _, item := range getPictures {
		customerArr := entity.UniqueCustomer{}
		customerArr.CustomerID = item.Customer.String
		customerArr.TeamID = item.Team.String
		customerArrs = append(customerArrs, customerArr)
	}
	UniqueCustomerList := UniqueCustomer(customerArrs)
	UniqueTeamList := UniqueTeam(teamArrs)
	outPuts := []model.Teams{}
	for _, teamValue := range UniqueTeamList {
		outPut := model.Teams{}
		outPut.TeamID = teamValue.TeamID
		outPut.TeamName = postgres.GetTeamName(teamValue.TeamID)
		customerOutputs := []model.Customer{}
		for _, customerValue := range UniqueCustomerList {
			customerOutput := model.Customer{}
			productValue := model.Images{}
			brandValue := model.Images{}
			surveyValue := model.Images{}
			competitorValue := model.Images{}
			flashBulletin := make([]model.Images, 0)
			promotionValue := model.Images{}
			if teamValue.TeamID == customerValue.TeamID {
				for _, item := range getPictures {
					if item.Team.String == teamValue.TeamID && item.Customer.String == customerValue.CustomerID {
						pictureModel, competitorName, surveyItem, promotionItem, err := mapper.MapPicturesEntityToModel(item)
						if err != nil {
							return response, nil
						}
						if strings.ToLower(item.Type.String) == "product" && len(pictureModel) > 0 {
							productValue.Name = item.Name.String
							productValue.URL = append(productValue.URL, pictureModel...)
						}
						if strings.ToLower(item.Type.String) == "brand" && len(pictureModel) > 0 {
							brandValue.Name = item.Name.String
							brandValue.URL = append(brandValue.URL, pictureModel...)
						}
						if strings.ToLower(item.Type.String) == "competitor" && len(pictureModel) > 0 {
							competitorValue.Name = competitorName
							competitorValue.URL = append(competitorValue.URL, pictureModel...)
						}
						if strings.ToLower(item.Type.String) == "survey" && len(pictureModel) > 0 {
							surveyValue.Name = surveyItem
							surveyValue.URL = append(surveyValue.URL, pictureModel...)
						}
						if strings.ToLower(item.Type.String) == "promotion" && len(pictureModel) > 0 {
							promotionValue.Name = promotionItem
							promotionValue.URL = append(promotionValue.URL, pictureModel...)
						}
					}
				}
				customerOutput.CustomerID = customerValue.CustomerID
				customerOutput.CustomerName = postgres.GetCustomerName(customerValue.CustomerID)
				customerOutput.Product = productValue
				customerOutput.Brand = brandValue
				customerOutput.Competitor = competitorValue
				customerOutput.Promotion = promotionValue
				customerOutput.Survey = surveyValue
				customerOutput.FlashBulletin = flashBulletin
				if len(productValue.URL) > 0 || len(brandValue.URL) > 0 || len(competitorValue.URL) > 0 || len(flashBulletin) > 0 || len(surveyValue.URL) > 0 || len(promotionValue.URL) > 0 {
					customerOutputs = append(customerOutputs, customerOutput)
				}
			}
		}
		outPut.Customers = customerOutputs

		if len(outPut.Customers) > 0 {
			outPuts = append(outPuts, outPut)
		}
	}
	if len(outPuts) < 1 {
		response.Message = "No data found"
	}
	response.Error = false
	response.Data = outPuts
	return response, nil
}

func UniqueTeam(brandSlice []entity.UniqueTeam) []entity.UniqueTeam {
	keys := make(map[entity.UniqueTeam]bool)
	uniqueList := []entity.UniqueTeam{}
	for _, entry := range brandSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func UniqueCustomer(brandSlice []entity.UniqueCustomer) []entity.UniqueCustomer {
	keys := make(map[entity.UniqueCustomer]bool)
	uniqueList := []entity.UniqueCustomer{}
	for _, entry := range brandSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func RetrievePictureZip(ctx *context.Context, input model.PictureZipInput, loggedInUserEntity *entity.LoggedInUser) (model.RetrievePictureZip, error) {
	response := model.RetrievePictureZip{}

	userName := postgres.GetUserName(loggedInUserEntity.ID)
	dirName := userName + "_" + input.Type + time.Now().Format("20060102150405") + "PictureRetrieve"
	err := os.Mkdir(dirName, 0755)
	for _, value := range input.Selections {
		valid, fileName := util.ValidateURL(value)
		if valid {
			file, downldErr := azure.DownloadFileFromBlobURL(value)
			if downldErr == nil {
				fileBytes, _ := ioutil.ReadAll(file)
				err = ioutil.WriteFile(dirName+"/"+fileName, fileBytes, 0644)
			}
		} else {
			continue
		}
	}

	dir, _ := os.Getwd()
	err = Zipit(dir+"/"+dirName, dir+"/"+dirName+".zip")
	zipByte, _ := ioutil.ReadFile(dir + "/" + dirName + ".zip")
	blobURL, err := azure.UploadBytesToBlob(zipByte, dirName+".zip")
	errd := os.RemoveAll(dir + "/" + dirName)
	errf := os.RemoveAll(dir + "/" + dirName + ".zip")
	if err != nil || errd != nil || errf != nil {
		response.Error = true
		response.Message = err.Error()
		return response, err
	}

	response.Error = false
	response.URL = blobURL
	return response, nil
}

// func CreateFolder(files []model.Teams, userAD *string) (string, error) {
// 	dirName := *userAD + time.Now().Format("20060102150405") + "PictureRetrieve"
// 	err := os.Mkdir(dirName, 0755)

// 	for _, item := range files {
// 		err = os.MkdirAll(dirName+"/"+item.TeamName, 0755)
// 		for _, each := range item.Customers {
// 			err = os.MkdirAll(dirName+"/"+item.TeamName+"/"+each.CustomerName, 0755)
// 			if len(each.Product.URL) > 0 {
// 				ProductFolderPath := dirName + "/" + item.TeamName + "/" + each.CustomerName + "/product_" + each.Product.Name
// 				err = os.MkdirAll(ProductFolderPath, 0755)
// 				for _, eachURL := range each.Product.URL {
// 					valid, fileName := util.ValidateURL(eachURL)
// 					if valid {
// 						file, downldErr := azure.DownloadFileFromBlobURL(eachURL)
// 						if downldErr == nil {
// 							fileBytes, _ := ioutil.ReadAll(file)
// 							err = ioutil.WriteFile(ProductFolderPath+"/"+fileName, fileBytes, 0644)
// 						} else {
// 							continue
// 						}
// 					}
// 				}
// 			}
// 			if len(each.Brand.URL) > 0 {
// 				brandFolderPath := dirName + "/" + item.TeamName + "/" + each.CustomerName + "/brand_" + each.Brand.Name
// 				err = os.MkdirAll(brandFolderPath, 0755)
// 				for _, eachURL := range each.Brand.URL {
// 					valid, fileName := util.ValidateURL(eachURL)
// 					if valid {
// 						file, downldErr := azure.DownloadFileFromBlobURL(eachURL)
// 						if downldErr == nil {
// 							fileBytes, _ := ioutil.ReadAll(file)
// 							err = ioutil.WriteFile(brandFolderPath+"/"+fileName, fileBytes, 0644)
// 						} else {
// 							continue
// 						}
// 					}
// 				}
// 			}
// 			if len(each.Competitor.URL) > 0 {
// 				competitorFolderPath := dirName + "/" + item.TeamName + "/" + each.CustomerName + "/competitor_" + each.Competitor.Name
// 				err = os.MkdirAll(competitorFolderPath, 0755)
// 				for _, eachURL := range each.Competitor.URL {
// 					valid, fileName := util.ValidateURL(eachURL)
// 					if valid {
// 						file, downldErr := azure.DownloadFileFromBlobURL(eachURL)
// 						if downldErr == nil {
// 							fileBytes, _ := ioutil.ReadAll(file)
// 							err = ioutil.WriteFile(competitorFolderPath+"/"+fileName, fileBytes, 0644)
// 						} else {
// 							continue
// 						}
// 					}
// 				}
// 			}
// 			if len(each.Survey.URL) > 0 {
// 				surveyFolderPath := dirName + "/" + item.TeamName + "/" + each.CustomerName + "/survey_" + each.Survey.Name
// 				err = os.MkdirAll(surveyFolderPath, 0755)
// 				for _, eachURL := range each.Survey.URL {
// 					valid, fileName := util.ValidateURL(eachURL)
// 					if valid {
// 						file, downldErr := azure.DownloadFileFromBlobURL(eachURL)
// 						if downldErr == nil {
// 							fileBytes, _ := ioutil.ReadAll(file)
// 							err = ioutil.WriteFile(surveyFolderPath+"/"+fileName, fileBytes, 0644)
// 						} else {
// 							continue
// 						}
// 					}
// 				}
// 			}
// 			if len(each.Promotion.URL) > 0 {
// 				promotionFolderPath := dirName + "/" + item.TeamName + "/" + each.CustomerName + "/promotion_" + each.Promotion.Name
// 				err = os.MkdirAll(promotionFolderPath, 0755)
// 				for _, eachURL := range each.Promotion.URL {
// 					valid, fileName := util.ValidateURL(eachURL)
// 					if valid {
// 						file, downldErr := azure.DownloadFileFromBlobURL(eachURL)
// 						if downldErr == nil {
// 							fileBytes, _ := ioutil.ReadAll(file)
// 							err = ioutil.WriteFile(promotionFolderPath+"/"+fileName, fileBytes, 0644)
// 						} else {
// 							continue
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	dir, _ := os.Getwd()
// 	err = Zipit(dir+"/"+dirName, dir+"/"+dirName+".zip")
// 	zipByte, _ := ioutil.ReadFile(dir + "/" + dirName + ".zip")
// 	blobURL, err := azure.UploadBytesToBlob(zipByte, dirName+".zip")
//	errd := os.RemoveAll(dir + "/" + dirName)
//	errf := os.RemoveAll(dir + "/" + dirName + ".zip")
// 	if err != nil {
// 		return "", err
// 	}
// 	return blobURL, nil
// }

func Zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}

func CustomerListData(input *model.ListInput) (model.CustomerListResponse, error) {
	getCustomerList, err := postgres.GetPictureRetrieveCustomerList(input)
	if err != nil {
		return model.CustomerListResponse{}, err
	}
	response := model.CustomerListResponse{}
	uniqueCustomer := make(map[string]int)
	for _, eachItem := range getCustomerList {
		customerModel, hasImage, err := mapper.MapCustomerListEntityToModel(eachItem)
		if err != nil {
			return response, nil
		}
		if hasImage {
			if _, exist := uniqueCustomer[customerModel.ID]; !exist {
				response.Data = append(response.Data, &customerModel)
				uniqueCustomer[customerModel.ID] = 1
			}
		}
	}
	if len(response.Data) < 1 {
		response.Message = "No data Found"
	}
	return response, nil
}

func ProductListData(input *model.ListInput) (model.ProductListResponse, error) {
	getProductList, err := postgres.GetPictureRetrieveProductList(input)
	if err != nil {
		return model.ProductListResponse{}, err
	}
	response := model.ProductListResponse{}
	uniqueProduct := make(map[string]int)
	for _, eachItem := range getProductList {
		productModel, hasImage, err := mapper.MapProductListEntityToModel(eachItem)
		if err != nil {
			return response, nil
		}
		if hasImage {
			if _, exist := uniqueProduct[productModel.ID]; !exist {
				response.Data = append(response.Data, &productModel)
				uniqueProduct[productModel.ID] = 1
			}
		}
	}
	if len(response.Data) < 1 {
		response.Message = "No data Found"
	}
	return response, nil
}

func BrandListData(input *model.ListInput) (model.BrandListResponse, error) {
	getBrandList, err := postgres.GetPictureRetrieveBrandList(input)
	if err != nil {
		return model.BrandListResponse{}, err
	}
	response := model.BrandListResponse{}
	uniqueBrand := make(map[string]int)
	for _, eachItem := range getBrandList {
		brandModel, hasImage, err := mapper.MapBrandListEntityToModel(eachItem)
		if err != nil {
			return response, nil
		}
		if hasImage {
			if _, exist := uniqueBrand[brandModel.ID]; !exist {
				response.Data = append(response.Data, &brandModel)
				uniqueBrand[brandModel.ID] = 1
			}
		}
	}
	if len(response.Data) < 1 {
		response.Message = "No data Found"
	}
	return response, nil
}
