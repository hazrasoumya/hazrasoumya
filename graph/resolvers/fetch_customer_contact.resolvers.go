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

func (r *queryResolver) GetGetCustomerContacts(ctx context.Context, input *model.GetCustomerContactRequest) (*model.GetCustomerContactResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/getGetCustomerContacts : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		response := model.GetCustomerContactResponse{}
		return &response, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response, err := controller.GetCustomerContact(input, userEntity)
	if err != nil {
		return nil, err
	}
	logengine.GetTelemetryClient().TrackRequest("GetGetCustomerContacts", "kpi/getGetCustomerContacts", time.Since(start), "200")
	return &response, err
}
