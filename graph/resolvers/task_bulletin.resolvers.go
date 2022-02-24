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

func (r *mutationResolver) UpsertTaskBulletin(ctx context.Context, input model.TaskBulletinUpsertInput) (*model.TaskBulletinUpsertResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/TaskBulletinUpsertInput : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.TaskBulletinUpsertResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	loggedInUser := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response := controller.UpsertTaskBulletin(input, loggedInUser)
	logengine.GetTelemetryClient().TrackRequest("UpsertTaskBulletin", "kpi/upsertTaskBulletin", time.Since(start), "200")
	return response, nil
}

func (r *mutationResolver) InsertCustomerTaskFeedBack(ctx context.Context, input model.CustomerTaskFeedBackInput) (*model.CustomerTaskFeedBackResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/insertCustomerTaskFeedBack : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.CustomerTaskFeedBackResponse{Error: true, Message: &errMessage}
		return responseModel, nil
	}
	loggedInUser := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response := controller.InsertCustomerTaskFeedBack(input, loggedInUser)
	logengine.GetTelemetryClient().TrackRequest("InsertCustomerTaskFeedBack", "kpi/insertCustomerTaskFeedBack", time.Since(start), "200")
	return response, nil
}

func (r *queryResolver) ListTaskBulletin(ctx context.Context, input *model.ListTaskBulletinInput) (*model.ListTaskBulletinResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/listTaskbulletin : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.ListTaskBulletinResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	response, err := controller.GetCustomerTakBulletin(input, *salesOrgID, *userID)
	if err != nil {
		return nil, err
	}
	logengine.GetTelemetryClient().TrackRequest("ListTaskbulletin", "kpi/listTaskbulletin", time.Since(start), "200")
	response.Error = false
	return &response, nil
}

func (r *queryResolver) FetchCustomerFeedback(ctx context.Context, input model.FetchCustomerFeedbackInput) (*model.FetchCustomerFeedbackResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/fetchCustomerFeedback : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.FetchCustomerFeedbackResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	response, err := controller.GetCustomerTakBulletinFeedback(input, *salesOrgID)
	if err != nil {
		return nil, err
	}
	logengine.GetTelemetryClient().TrackRequest("FetchCustomerFeedback", "kpi/fetchCustomerFeedback", time.Since(start), "200")
	response.Error = false
	return &response, nil
}

func (r *queryResolver) TeamToCustomerDropDown(ctx context.Context, input *model.TaskBulletinInput) (*model.TaskBulletinResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/TeamToCustomerDropDown : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.TaskBulletinResponse{Error: true, Message: errMessage}
		return responseModel, nil
	}
	loggedInUser := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response := controller.TeamToCustomerDropDown(input, loggedInUser)
	logengine.GetTelemetryClient().TrackRequest("TeamToCustomerDropDown", "kpi/TeamToCustomerDropDown", time.Since(start), "200")
	return response, nil
}

func (r *queryResolver) PrincipalDropDown(ctx context.Context, input model.PrincipalDropDownInput) (*model.PrincipalDropDownResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/principalDropDown : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.PrincipalDropDownResponse{Error: true, Message: &errMessage}
		return responseModel, nil
	}
	loggedInUser := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response, err := controller.FetchPrincipalNameDropDown(input, loggedInUser)
	if err != nil {
		return nil, err
	}
	logengine.GetTelemetryClient().TrackRequest("PrincipalDropDown", "kpi/principalDropDown", time.Since(start), "200")
	return response, nil
}

func (r *queryResolver) TaskBulletinTitleDropDown(ctx context.Context, input *model.TaskBulletinTitleInput) (*model.TaskBulletinTitleResponse, error) {
	start := time.Now().UTC()
	inputJson, _ := json.Marshal(input)
	logengine.GetTelemetryClient().TrackEvent("kpi/taskBulletinTypeDropDown : " + string(inputJson))
	userID, salesOrgID, authRole := auth.ForContext(ctx)
	if userID == nil || salesOrgID == nil || authRole == nil {
		errMessage := "You are unauthorized, please login!"
		responseModel := &model.TaskBulletinTitleResponse{Error: true, Message: &errMessage}
		return responseModel, nil
	}
	loggedInUser := &entity.LoggedInUser{ID: *userID, SalesOrganisaton: *salesOrgID, AuthRole: *authRole}
	response, err := controller.FetchTaskBulletinTypeDropDown(input, loggedInUser)
	if err != nil {
		return nil, err
	}
	logengine.GetTelemetryClient().TrackRequest("TaskBulletinTypeDropDown", "kpi/taskBulletinTypeDropDown", time.Since(start), "200")
	return response, nil
}
