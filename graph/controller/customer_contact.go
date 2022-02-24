package controller

import (
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/mapper"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
)

func UpsertCustomerContacts(inputModel model.CustomerContactRequest, loggedInUserEntity *entity.LoggedInUser) *model.KpiResponse {
	entity, validationResult := mapper.MapCustomerContactModelToEntity(&inputModel, loggedInUserEntity)
	kpiResponse := &model.KpiResponse{}
	if !validationResult.Error {
		postgres.UpsertCustomerContactData(entity, kpiResponse, loggedInUserEntity)
	} else {
		kpiResponse.Error = true
		kpiResponse.Message = "Customer Contact Not Valid"
		kpiResponse.ValidationErrors = validationResult.ValidationMessage
	}
	return kpiResponse
}

func GetCustomerContact(inputModel *model.GetCustomerContactRequest, loggedInUserEntity *entity.LoggedInUser) (model.GetCustomerContactResponse, error) {
	getContacts, err := postgres.GetCustomerContactData(inputModel, loggedInUserEntity)
	if err != nil {
		return model.GetCustomerContactResponse{}, err
	}
	hasConsentForm := postgres.HasConsentFormSalesOrg("consentform", loggedInUserEntity.SalesOrganisaton)
	response := model.GetCustomerContactResponse{}
	for _, contact := range getContacts {
		contactModel := model.CustomerContact{}
		contactModel.ID = contact.ID.String
		contactModel.ContactName = contact.ContactName.String
		contactModel.Designation = contact.Designation.String
		contactModel.ContactNumber = contact.ContactNumber.String
		contactModel.ContactImage = contact.ContactImage.String
		contactModel.CustomerID = contact.CustomerID.String
		contactModel.CustomerName = contact.CustomerName.String
		contactModel.EmailID = contact.EmailID.String
		if hasConsentForm {
			contactModel.HasConsent = contact.HasConsent.Bool
		} else {
			contactModel.HasConsent = true
		}
		response.GetCustomerContact = append(response.GetCustomerContact, &contactModel)
	}

	return response, nil
}

func DeleteCustomerContacts(inputModel model.CustomerContactDeleteRequest, loggedInUserEntity *entity.LoggedInUser) *model.KpiResponse {
	entity, validationResult := mapper.MapDeleteCustomerContactModelToEntity(&inputModel)
	kpiResponse := &model.KpiResponse{}
	if !validationResult.Error {
		postgres.DeleteCustomerContactData(entity, kpiResponse, loggedInUserEntity)
	} else {
		kpiResponse.Error = true
		kpiResponse.Message = "Customer Contact Details Not Valid"
		kpiResponse.ValidationErrors = validationResult.ValidationMessage
	}
	return kpiResponse
}
