package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eztrade/kpi/graph/controller"
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/generated"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/login/auth"
)

func (r *mutationResolver) SaveCustomerTarget(ctx context.Context, input model.CustomerTargetInput) (*model.KpiResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/saveCustomerTarget : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.KpiResponse{Error: true, Message: errMessage, ErrorCode: 999}
		return responseModel, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response := controller.UpsertCustomerTarget(input, userEntity)
	logengine.GetTelemetryClient().TrackRequest("SaveKpiAnswers", "kpi/SaveCustomerTarget", time.Since(start), "200")
	return response, nil
}

func (r *queryResolver) GetTargetCustomer(ctx context.Context, input model.GetTargetCustomerRequest) (*model.GetTargetCustomerResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/getTargetCustomer : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.GetTargetCustomerResponse{Error: true, Message: &errMessage}
		return responseModel, nil
	}
	logengine.GetTelemetryClient().TrackRequest("GetTargetCustomer", "kpi/getTargetCustomer", time.Since(start), "200")
	response, err := controller.GetTargetCustomer(input, &ctx)
	if err != nil {
		return nil, err
	}
	response.Error = false
	return response, nil
}

func (r *queryResolver) GetCustomerGroup(ctx context.Context, input *model.CustomerGroupInput) (*model.CustomerGroupResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/getTargetCustomer : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.CustomerGroupResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	logengine.GetTelemetryClient().TrackRequest("getCustomerGroup", "kpi/getCustomerGroup", time.Since(start), "200")
	response, err := controller.GetCustomerGroup(input, userEntity)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
