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

func (r *mutationResolver) UpsertKpiTarget(ctx context.Context, input model.KPITargetInput) (*model.KpiResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/upsertKpiTarget : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.KpiResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response := controller.SaveTargetKpi(input, userEntity)
	logengine.GetTelemetryClient().TrackRequest("UpsertKpiTarget", "kpi/upsertKpiTarget", time.Since(start), "200")
	return response, nil
}

func (r *mutationResolver) ActionKpiTarget(ctx context.Context, input model.ActionKPITargetInput) (*model.KpiResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/actionKpiTarget : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.KpiResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	if *authRole == "countrymanager" {
		userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
		response := controller.ActionTargetKpi(input, userEntity)
		logengine.GetTelemetryClient().TrackRequest("ActionKpiTarget", "kpi/actionKpiTarget", time.Since(start), "200")
		return response, nil
	} else {
		errMessage := "You don't have any access for this service!"
		responseModel := &model.KpiResponse{Error: true, Message: errMessage}
		logengine.GetTelemetryClient().TrackRequest("ActionKpiTarget", "kpi/actionKpiTarget", time.Since(start), "200")
		return responseModel, nil
	}
}
