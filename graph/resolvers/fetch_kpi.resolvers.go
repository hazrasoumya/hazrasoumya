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

func (r *queryResolver) GetKpis(ctx context.Context, input *model.GetKpiInput) (*model.GetKpiResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/getKpis : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.GetKpiResponse{Error: true, Message: &errMessage, ErrorCode: 999}
		return responseModel, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response, err := controller.GetKpiInformation(input, userEntity)
	if err != nil {
		mag := err.Error()
		response.ErrorCode = 999
		response.Message = &mag
		return &response, nil
	}
	response.Error = false
	logengine.GetTelemetryClient().TrackRequest("GetKpis", "kpi/getKpis", time.Since(start), "200")
	return &response, err
}

func (r *queryResolver) GetKpiQuestionAnswerOffline(ctx context.Context, input model.KpiOfflineInput) (*model.KpiOfflineResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/GetKpiQuestionAnswerOffline : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.KpiOfflineResponse{Error: true, Message: errMessage, ErrorCode: 999}
		return responseModel, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response, err := controller.GetEventKpiProductBrand(input, userEntity)
	if err != nil {
		mag := err.Error()
		response.ErrorCode = 999
		response.Message = mag
		return response, nil
	}
	response.Error = false
	logengine.GetTelemetryClient().TrackRequest("GetKpiQuestionAnswerOffline", "kpi/GetKpiQuestionAnswerOffline", time.Since(start), "200")
	return response, err
}
