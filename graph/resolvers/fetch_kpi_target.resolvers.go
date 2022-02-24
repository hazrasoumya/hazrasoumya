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

func (r *queryResolver) GetKpiTargets(ctx context.Context, input *model.GetKpiTargetRequest) (*model.GetKpiTargetResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/getKpiTargets : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		response := model.GetKpiTargetResponse{}
		return &response, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response, err := controller.GetKpiTargetData(input, userEntity)
	if err != nil {
		return nil, err
	}
	logengine.GetTelemetryClient().TrackRequest("GetKpiTargets", "kpi/getKpiTargets", time.Since(start), "200")
	return &response, err
}

func (r *queryResolver) GetKpiTargetTitle(ctx context.Context) (*model.KpiTaregetTitleResponse, error) {
	start := time.Now().UTC()
	logengine.GetTelemetryClient().TrackEvent("kpi/getKpiTargetTitle : NO INPUT")
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.KpiTaregetTitleResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	response, err := controller.GetKpiTargetTitle()
	if err != nil {
		return nil, err
	}
	response.Error = false
	logengine.GetTelemetryClient().TrackRequest("GetKpiTargetTitle", "kpi/getKpiTargetTitle", time.Since(start), "200")
	return &response, err
}
