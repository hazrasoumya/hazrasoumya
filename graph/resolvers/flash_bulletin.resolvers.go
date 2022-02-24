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

func (r *mutationResolver) UpsertFlashBulletin(ctx context.Context, input model.FlashBulletinUpsertInput) (*model.FlashBulletinUpsertResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/upsertFlashBulletin : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.FlashBulletinUpsertResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	loggedInUser := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response := controller.UpsertFlashBulletin(input, loggedInUser)
	logengine.GetTelemetryClient().TrackRequest("UpsertFlashBulletin", "kpi/upsertFlashBulletin", time.Since(start), "200")
	return response, nil
}

func (r *queryResolver) RetriveFlashBulletinSingle(ctx context.Context, input model.RetriveInfoFlashBulletinInput) (*model.RetriveInfoFlashBulletinleResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/retriveFlashBulletinSingle : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.RetriveInfoFlashBulletinleResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	loggedInUser := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response := controller.RetriveFlashBulletinForEdit(input, loggedInUser)
	logengine.GetTelemetryClient().TrackRequest("RetriveFlashBulletinSingle", "kpi/retriveFlashBulletinSingle", time.Since(start), "200")
	return response, nil
}

func (r *queryResolver) ListFlashbulletin(ctx context.Context, input model.ListFlashBulletinInput) (*model.ListFlashBulletinResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/listFlashbulletin : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.ListFlashBulletinResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	loggedInUser := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response := controller.FlashBulletinList(input, loggedInUser)
	logengine.GetTelemetryClient().TrackRequest("ListFlashbulletin", "kpi/listFlashbulletin", time.Since(start), "200")
	return &response, nil
}
