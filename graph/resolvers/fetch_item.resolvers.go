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

func (r *queryResolver) GetKpiBrandProduct(ctx context.Context, input model.GetBrandProductRequest) (*model.GetKpiBrandProductResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/getKpiBrandProductResponse : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.GetKpiBrandProductResponse{Error: true, Message: errMessage, ErrorCode: 999}
		return responseModel, nil
	}
	userEntity := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response, err := controller.GetKpiBrandProductData(input, userEntity)
	if err != nil {
		return nil, err
	}
	response.Error = false
	logengine.GetTelemetryClient().TrackRequest("GetKpiBrandProductResponse", "kpi/getKpiBrandProductResponse", time.Since(start), "200")
	return &response, err
}
