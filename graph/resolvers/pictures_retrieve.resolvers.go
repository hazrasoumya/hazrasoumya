package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eztrade/kpi/graph/controller"
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/login/auth"
)

func (r *queryResolver) RetrievePictures(ctx context.Context, input model.PicturesInput) (*model.RetrievePicturesResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/retrievePictures : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.RetrievePicturesResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response, err := controller.RetrievePictures(input, userEntity)
	if err != nil {
		return nil, err
	}
	response.Error = false
	logengine.GetTelemetryClient().TrackRequest("RetrievePictures", "kpi/retrievePictures", time.Since(start), "200")
	return &response, err
}

func (r *queryResolver) RetrievePictureZip(ctx context.Context, input model.PictureZipInput) (*model.RetrievePictureZip, error) {
	// panic(fmt.Errorf("not implemented"))
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/retrievePictureZip : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.RetrievePictureZip{Error: true, Message: errMessage}
		return responseModel, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response, err := controller.RetrievePictureZip(&ctx, input, userEntity)
	if err != nil {
		return nil, err
	}
	response.Error = false
	logengine.GetTelemetryClient().TrackRequest("RetrievePictureZip", "kpi/retrievePictureZip", time.Since(start), "200")
	return &response, err
}

func (r *queryResolver) RetrievePictureCustomerList(ctx context.Context, input *model.ListInput) (*model.CustomerListResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/retrievePictureCustomerList : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.CustomerListResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	response, err := controller.CustomerListData(input)
	if err != nil {
		return nil, err
	}
	response.Error = false
	logengine.GetTelemetryClient().TrackRequest("RetrievePictureCustomerList", "kpi/retrievePictureCustomerList", time.Since(start), "200")
	return &response, err
}

func (r *queryResolver) RetrievePictureProductList(ctx context.Context, input *model.ListInput) (*model.ProductListResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/retrievePictureProductList : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.ProductListResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	response, err := controller.ProductListData(input)
	if err != nil {
		return nil, err
	}
	response.Error = false
	logengine.GetTelemetryClient().TrackRequest("RetrievePictureProductList", "kpi/retrievePictureProductList", time.Since(start), "200")
	return &response, err
}

func (r *queryResolver) RetrievePictureBrandList(ctx context.Context, input *model.ListInput) (*model.BrandListResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/retrievePictureBrandList : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.BrandListResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	response, err := controller.BrandListData(input)
	if err != nil {
		return nil, err
	}
	response.Error = false
	logengine.GetTelemetryClient().TrackRequest("RetrievePictureBrandList", "kpi/retrievePictureBrandList", time.Since(start), "200")
	return &response, err
}
