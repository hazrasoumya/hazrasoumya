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

func (r *mutationResolver) UpsertCustomerContact(ctx context.Context, input model.CustomerContactRequest) (*model.KpiResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/upsertCustomerContact : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.KpiResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response := controller.UpsertCustomerContacts(input, userEntity)
	logengine.GetTelemetryClient().TrackRequest("UpsertCustomerContact", "kpi/upsertCustomerContact", time.Since(start), "200")
	return response, nil
}

func (r *mutationResolver) DeleteCustomerContact(ctx context.Context, input model.CustomerContactDeleteRequest) (*model.KpiResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/deleteCustomerContact : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.KpiResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response := controller.DeleteCustomerContacts(input, userEntity)
	logengine.GetTelemetryClient().TrackRequest("DeleteCustomerContact", "kpi/deleteCustomerContact", time.Since(start), "200")
	return response, nil
}
