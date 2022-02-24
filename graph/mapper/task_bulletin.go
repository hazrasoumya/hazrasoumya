package mapper

import (
	"errors"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	"github.com/eztrade/kpi/graph/util"
	uuid "github.com/gofrs/uuid"
)

func MapTaskBulletinInputModelToEntity(inputModel *model.TaskBulletinUpsertInput, loggedInUser *entity.LoggedInUser) (*entity.TaskBulletin, *model.ValidationResult) {
	var id *uuid.UUID
	result := &model.ValidationResult{Error: false}
	entity := &entity.TaskBulletin{}

	if loggedInUser.AuthRole != "sfe" {
		if !postgres.IsLineOneManager(loggedInUser.ID) {
			result.Error = true
			errorMessage := &model.ValidationMessage{Row: 0, Message: "You are not allowed to create Task Bulletin!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			return entity, result
		}
	}
	if inputModel.ID == nil {
		if inputModel.IsDeleted != nil && *inputModel.IsDeleted {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "To Delete Id is mandatory!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}

	if inputModel.ID != nil {
		uuid, err := uuid.FromString(*inputModel.ID)
		if err == nil {
			id = &uuid
		} else {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Task Bulletin ID format is invalid!"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
	}
	if inputModel.PrincipalName == "" {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Principal name can not be blank"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	entity.PrincipalName = inputModel.PrincipalName
	entity.ID = id
	if inputModel.Type != "" {
		typeId, err := postgres.ValidateTaskBulletinType(inputModel.Type)
		if err != nil {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: err.Error()}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		}
		entity.TypeID = typeId
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Type is Empty!"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	if inputModel.Description != "" {
		entity.Description = inputModel.Description
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Description can not be empty!"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	if len(inputModel.TeamMemberCustomer) > 0 {
		for _, value := range inputModel.TeamMemberCustomer {
			_, err := uuid.FromString(*value)
			if err != nil {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Wrong TeamMemberCustomer id"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			} else {
				_, err := postgres.HasTeamMemberCustomerId(*value, loggedInUser.SalesOrganisaton)
				if err != nil {
					if !result.Error {
						result.Error = true
						result.ValidationMessage = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "TeamMemberCustomer Doesn't Exist"}
					result.ValidationMessage = append(result.ValidationMessage, errorMessage)
				} else {
					entity.TeamMemberCustomer = append(entity.TeamMemberCustomer, *value)
				}
			}
		}
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "TeamMemberCustomer can not be blank"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	if inputModel.Title != "" {
		entity.Title = inputModel.Title
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Title is Empty"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	creationDateObj, creationDateValidateErr := util.IsValidDateWithDateObect(inputModel.CreationDate)
	if creationDateValidateErr != nil {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Incorrect creation date"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	targetDateObj, targetDateValidateErr := util.IsValidDateWithDateObect(inputModel.TargetDate)
	if targetDateValidateErr != nil {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Incorrect Target date"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	if creationDateValidateErr == nil && targetDateValidateErr == nil {
		if creationDateObj.After(targetDateObj) {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "creation date is greater than Target date"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		} else {
			entity.CreationDate = inputModel.CreationDate
			entity.TargetDate = inputModel.TargetDate
		}
	}

	userID, err := uuid.FromString(loggedInUser.ID)
	if err != nil {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid uuid"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	if inputModel.ID != nil {
		entity.ModifiedBy = userID
	} else {
		entity.CreatedBy = userID
	}
	flag1 := true
	flag2 := false
	if inputModel.IsDeleted == nil || !*inputModel.IsDeleted {
		entity.IsActive = &flag1
		entity.IsDeleted = &flag2
	} else {
		entity.IsActive = &flag2
		entity.IsDeleted = &flag1
	}

	if inputModel.Attachments != nil {
		entity.Attachments = MapAttachmentTaskInputModelsToEntities(inputModel, userID, result)
	}
	return entity, result
}

func MapAttachmentTaskInputModelsToEntities(inputModel *model.TaskBulletinUpsertInput, userId uuid.UUID, result *model.ValidationResult) []entity.Attachment {
	var entities []entity.Attachment
	if len(inputModel.Attachments) < 1 {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Attachments are required"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	updateFlag := inputModel.ID != nil
	for _, attachment := range inputModel.Attachments {
		var entity entity.Attachment
		if attachment.ID != nil && updateFlag {
			var id *uuid.UUID
			uuid, err := uuid.FromString(*attachment.ID)
			if err == nil {
				id = &uuid
				err = postgres.AttachmentBelongsToTaskBulletin(inputModel.ID, uuid)
				if err != nil {
					if !result.Error {
						result.Error = true
						result.ValidationMessage = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Attachment does not belong to Task Bulletin!"}
					result.ValidationMessage = append(result.ValidationMessage, errorMessage)
				}
			} else {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Attachment ID format is invalid!"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
			entity.ID = id
		}
		entity.BlobName = attachment.Filename
		entity.BlobUrl = attachment.URL
		if updateFlag {
			entity.ModifiedBy = userId
		} else {
			entity.CreatedBy = userId
		}
		entity.ValidateData(0, result)
		entities = append(entities, entity)
	}
	return entities
}

func ValidateCustomerFeedbackInput(input model.FetchCustomerFeedbackInput, salesorg string) (entity.TaskBulletinfeedbackInput, error) {
	entity := &entity.TaskBulletinfeedbackInput{}
	if input.TeamMemberCustomerID == "" || input.TaskBulletinID == "" {
		return *entity, errors.New("TeamMemberCustomerId and TaskBulletinID Both Should Be Provided")
	}
	if input.TeamMemberCustomerID != "" {
		uuid, err := uuid.FromString(input.TeamMemberCustomerID)
		if err != nil {
			return *entity, errors.New("TeamMemberCustomerId format is invalid")
		} else {
			_, err := postgres.HasTeamMemberCustomerId(input.TeamMemberCustomerID, salesorg)
			if err != nil {
				return *entity, errors.New("TeamMemberCustomerId not found")
			} else {
				entity.TeamMemberCustomerId = uuid
			}
		}
	}
	if input.TaskBulletinID != "" {
		uuid, err := uuid.FromString(input.TaskBulletinID)
		if err != nil {
			return *entity, errors.New("TaskBulletinID format is invalid")
		} else {
			_, err := postgres.HasTaskBulletinID(input.TaskBulletinID)
			if err != nil {
				return *entity, errors.New("TaskBulletin ID not found")
			} else {
				entity.TaskBulletinId = uuid
			}
		}
	}
	return *entity, nil

}

func MapCustomerTaskFeedbackToModel(input []entity.TaskBulletinFeedbackOutput) []*model.CustomerFeedback {
	response := []*model.CustomerFeedback{}
	for _, data := range input {
		res := model.CustomerFeedback{}
		res.StatusTitle = data.StatusTitle.String
		res.StatusValue = data.StatusValue.String
		res.Remarks = data.Remarks.String
		res.DateCreated = data.DateCreated.String
		attachments, _ := postgres.GetAttachmentsByIDs(data.Attachment.String)
		res.Attachments = attachments
		response = append(response, &res)

	}
	return response

}

func MapCustomerTaskFeedBackInputModelToEntity(inputModel *model.CustomerTaskFeedBackInput, loggedInUser *entity.LoggedInUser) (*entity.CustomerTaskFeedBack, *model.ValidationResult) {
	result := &model.ValidationResult{Error: false}
	entitys := &entity.CustomerTaskFeedBack{}
	if inputModel.TaskBulletinID != "" {
		uuid, err := uuid.FromString(inputModel.TaskBulletinID)
		if err != nil {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "TaskBulletin ID format invalid"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		} else {
			_, err := postgres.HasTaskBulletinID(inputModel.TaskBulletinID)
			if err != nil {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "TaskBulletin ID Not Found"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			} else {
				entitys.TaskBulletinId = uuid
			}
		}
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "TaskBulletin ID Can Not Be Blank"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	if inputModel.TeamMemberCustomerID != "" {
		uuid, err := uuid.FromString(inputModel.TeamMemberCustomerID)
		if err != nil {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "TeamMemberCustomer ID format invalid"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		} else {
			_, err := postgres.HasTeamMemberCustomerId(inputModel.TeamMemberCustomerID, loggedInUser.SalesOrganisaton)
			if err != nil {

				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Team member Customer Id Does Not Exists"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			} else {
				entitys.TeamMemberCustomerId = uuid
			}
		}
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: " Team member Customer Id Cant Not Be Blank"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}

	if inputModel.Status != "" {
		statusCode, err := postgres.GetCodeIdForStatus(inputModel.Status, "TaskFeedback")
		if err != nil {
			if !result.Error {
				result.Error = true
				result.ValidationMessage = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Invalid customer feedback Status"}
			result.ValidationMessage = append(result.ValidationMessage, errorMessage)
		} else {
			entitys.Status = *statusCode
		}
	} else {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Customer Feedback Status Can Not Be Blank"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)
	}
	if inputModel.Remarks != nil && *inputModel.Remarks != "" {
		entitys.Remarks = *inputModel.Remarks
	}
	entitys.CreatedBy = loggedInUser.ID

	if !postgres.CombiNationExistsForFeedback(inputModel.TaskBulletinID, inputModel.TeamMemberCustomerID) {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Customer Feedback Combination Does Not Exists"}
		result.ValidationMessage = append(result.ValidationMessage, errorMessage)

	}
	if inputModel.Attachments != nil {
		for _, data := range inputModel.Attachments {
			var attachmentsData entity.Attachmentdata
			if data.ID != nil && *data.ID != "" {
				attachmentsData.Id = *data.ID
			}

			if data.Filename == "" {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Attachments File Name can not be blank"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			} else {
				attachmentsData.BlobName = data.Filename
			}

			if data.URL == "" {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Attachments URL can not be blank"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			} else if util.IsValidUrl(data.URL) {
				attachmentsData.BlobUrl = data.URL
			} else {
				if !result.Error {
					result.Error = true
					result.ValidationMessage = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Attachments URL Format Is Invalid"}
				result.ValidationMessage = append(result.ValidationMessage, errorMessage)
			}
			entitys.Attachments = append(entitys.Attachments, attachmentsData)
		}
	}

	return entitys, result
}

func ValidateCustomerCustomerTakBulletinInput(input *model.ListTaskBulletinInput, salesorg string, userId string) (entity.TaskBulletinInput, error) {
	entity := entity.TaskBulletinInput{}
	if input != nil {
		if input.TeamMemberCustomerID != nil {
			_, err := uuid.FromString(*input.TeamMemberCustomerID)
			if err != nil {
				return entity, errors.New("TeamMemberCustomerId format is invalid")
			} else {
				_, err := postgres.HasTeamMemberCustomerId(*input.TeamMemberCustomerID, salesorg)
				if err != nil {
					return entity, errors.New("TeamMemberCustomerId not found")
				} else {
					entity.TeamMemberCustomerId = input.TeamMemberCustomerID
				}
			}
		}

		if input.ID != nil {
			_, err := uuid.FromString(*input.ID)
			if err != nil {
				return entity, errors.New("TaskBulletinID format is invalid")
			} else {
				_, err := postgres.HasTaskBulletinID(*input.ID)
				if err != nil {
					return entity, errors.New("TaskBulletin ID not found")
				} else {
					entity.Id = input.ID
				}
			}
		}

		if input.Type != nil {
			_, err := postgres.IsValidTaskBulletinType(*input.Type)
			if err != nil {
				return entity, errors.New("task Bulletin type is invalid")
			}
			entity.Type = input.Type
		}

		if input.CreationDate != nil || input.TargetDate != nil {
			if input.CreationDate != nil {
				_, creationDateValidateErr := util.IsValidDateWithDateObect(*input.CreationDate)
				if creationDateValidateErr != nil {
					return entity, errors.New("creation is invalid")
				}
				entity.CreationDate = input.CreationDate
			}
			if input.TargetDate != nil {
				_, targetDateValidateErr := util.IsValidDateWithDateObect(*input.TargetDate)
				if targetDateValidateErr != nil {
					return entity, errors.New("target is invalid")
				}
				entity.TargetDate = input.TargetDate
			}
			if input.CreationDate != nil && input.TargetDate != nil {
				creationDateObj, _ := util.IsValidDateWithDateObect(*input.CreationDate)
				targetDateObj, _ := util.IsValidDateWithDateObect(*input.TargetDate)
				if creationDateObj.After(targetDateObj) {
					return entity, errors.New("creation day is after the target date")

				}
			}
		}
		if input.TeamMemberID != nil && *input.TeamMemberID != "" {
			_, err := uuid.FromString(*input.TeamMemberID)
			if err != nil {
				return entity, errors.New("TeamMemberId format is invalid")
			} else {
				_, err := postgres.HasTeamMemberId(*input.TeamMemberID, salesorg)
				if err != nil {
					return entity, errors.New("TeamMemberId not found")
				} else {
					entity.TeamMemberId = input.TeamMemberID
				}
			}

		}

		if input.SearchItem != nil && *input.SearchItem != "" {
			entity.SearchItem = input.SearchItem
		}

		if input.Limit != nil && *input.Limit > 0 {
			entity.Limit = input.Limit
			var initalOffset = 0
			entity.Offset = &initalOffset
			if input.PageNo != nil && *input.PageNo > 0 {
				entity.PageNo = input.PageNo
				*entity.Offset = *input.Limit * (*input.PageNo - 1)
			}
		}

		entity.UserId = &userId
		if input.IsActive != nil && *input.IsActive {
			entity.IsActive = true
		} else {
			entity.IsActive = false
		}
	}
	entity.Salesorg = &salesorg

	return entity, nil

}

func MapTaskBulletinEntityToModel(input []entity.TaskBulletinOutput) []*model.TaskBulletinData {
	var outputModel []*model.TaskBulletinData
	teamMembers := []entity.UniqueTeammember{}
	bulletins := []entity.UniqueTaskBulletin{}
	attachments := []entity.UniqueAttachment{}
	customers := []entity.UniqueCustomers{}
	for _, data := range input {
		teamMember := entity.UniqueTeammember{}
		bulletin := entity.UniqueTaskBulletin{}
		attach := entity.UniqueAttachment{}
		customer := entity.UniqueCustomers{}
		teamMember.TaskBulletinId = data.Id.String
		teamMember.TeamId = data.TeamId.String
		teamMember.TeamMemberId = data.TeamMemberId.String
		teamMember.UserId = data.UserId.String
		teamMember.FirstName = data.FirstName.String
		teamMember.LastName = data.LastName.String
		teamMember.ActiveDirectory = data.ActiveDirectory.String
		teamMember.Email = data.Email.String
		teamMember.ApprovalRoleTitle = data.ApprovalRoleTitle.String
		teamMember.ApprovalRoleValues = data.ApprovalRoleValues.String
		bulletin.Id = data.Id.String
		bulletin.TypeTitle = data.TypeTitle.String
		bulletin.TypeValue = data.TypeValue.String
		bulletin.Title = data.Title.String
		bulletin.Description = data.Description.String
		bulletin.TeamId = data.TeamId.String
		bulletin.TeamName = data.TeamName.String
		bulletin.CreationDate = data.CreationDate.String
		bulletin.TargetDate = data.TargetDate.String
		bulletin.PrincipalName = data.PrincipalName.String
		attach.Id = data.Id.String
		attach.BlobId = data.BlobId.String
		attach.BlobName = data.BlobName.String
		attach.Url = data.BlobURL.String
		customer.TeamMemberId = data.TeamMemberId.String
		customer.TeamMemberCustomerID = data.TeamMemberCustomerId.String
		customer.CusomerId = data.CustomerID.String
		customer.CustomerName = data.CustomerName.String
		customer.SoldTo = data.SoldTo
		customer.ShipTo = data.ShipTo
		customer.TaskBulletinId = data.Id.String
		teamMembers = append(teamMembers, teamMember)
		bulletins = append(bulletins, bulletin)
		attachments = append(attachments, attach)
		customers = append(customers, customer)

	}
	uniquetaskBuletin := UniqueTaskBulletins(bulletins)
	uniqueteammember := UniqueTeamMembers(teamMembers)
	uniqueAttachments := UniqueAttachments(attachments)
	uniqueCustomer := UniqueCustoers(customers)

	for _, bulletin := range uniquetaskBuletin {
		var innerModel model.TaskBulletinData
		var teamMembers []*model.SalesRepData
		innerModel.ID = bulletin.Id
		innerModel.TypeTitle = bulletin.TypeTitle
		innerModel.TypeValue = bulletin.TypeValue
		innerModel.Title = bulletin.Title
		innerModel.Description = bulletin.Description
		innerModel.PrincipalName = bulletin.PrincipalName
		team := bulletin.TeamId
		innerModel.TeamID = &team
		innerModel.TeamName = bulletin.TeamName
		innerModel.TargetDate = bulletin.TargetDate
		innerModel.CreationDate = bulletin.CreationDate
		var attachments []*model.Attachment
		for _, attch := range uniqueAttachments {
			if attch.Id == bulletin.Id {
				var attachment model.Attachment
				attachment.ID = attch.BlobId
				attachment.URL = attch.Url
				attachment.Filename = attch.BlobName
				attachments = append(attachments, &attachment)
			}
		}
		innerModel.Attachments = attachments
		for _, teateamMembers := range uniqueteammember {
			if bulletin.TeamId == teateamMembers.TeamId {
				if teateamMembers.TaskBulletinId == bulletin.Id {
					var teamMember model.SalesRepData
					TeamMemberID := teateamMembers.TeamMemberId
					teamMember.TeamMemberID = &TeamMemberID
					user := teateamMembers.UserId
					teamMember.UserID = &user
					firstName := teateamMembers.FirstName
					teamMember.FirstName = &firstName
					lastName := teateamMembers.LastName
					teamMember.LastName = &lastName
					activeDir := teateamMembers.ActiveDirectory
					teamMember.ActiveDirectory = &activeDir
					email := teateamMembers.Email
					teamMember.Email = &email
					approveTitle := teateamMembers.ApprovalRoleTitle
					teamMember.ApprovalRoleTitle = &approveTitle
					approveValue := teateamMembers.ApprovalRoleValues
					teamMember.ApprovalRoleValues = &approveValue

					var customers []*model.CustomerData
					for _, value := range uniqueCustomer {
						if teateamMembers.TeamMemberId == value.TeamMemberId {
							if value.TaskBulletinId == bulletin.Id {
								var customer model.CustomerData
								customerId := value.CusomerId
								customer.CustomerID = &customerId
								tmcId := value.TeamMemberCustomerID
								customer.TeamMemberCustomerID = &tmcId
								customerName := value.CustomerName
								customer.CustomerName = &customerName
								soldTo := value.SoldTo
								customer.SoldTo = &soldTo
								shipTo := value.ShipTo
								customer.ShipTo = &shipTo
								cm := &customer
								customers = append(customers, cm)
							}
						}
					}
					teamMember.Customers = customers
					tm := &teamMember
					teamMembers = append(teamMembers, tm)
				}
			}
		}
		innerModel.SalesRep = teamMembers
		im := &innerModel
		outputModel = append(outputModel, im)
	}
	return outputModel
}

func UniqueTeamMembers(teammember []entity.UniqueTeammember) []entity.UniqueTeammember {
	keys := make(map[entity.UniqueTeammember]bool)
	uniqueList := []entity.UniqueTeammember{}
	for _, entry := range teammember {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func UniqueTaskBulletins(taskBulletin []entity.UniqueTaskBulletin) []entity.UniqueTaskBulletin {
	keys := make(map[entity.UniqueTaskBulletin]bool)
	uniqueList := []entity.UniqueTaskBulletin{}
	for _, entry := range taskBulletin {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func UniqueAttachments(attachment []entity.UniqueAttachment) []entity.UniqueAttachment {
	keys := make(map[entity.UniqueAttachment]bool)
	uniqueList := []entity.UniqueAttachment{}
	for _, entry := range attachment {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func UniqueCustoers(cusomers []entity.UniqueCustomers) []entity.UniqueCustomers {
	keys := make(map[entity.UniqueCustomers]bool)
	uniqueList := []entity.UniqueCustomers{}
	for _, entry := range cusomers {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func TeamToCustomerEntityToModel(input []entity.TeamToCustomer) []*model.TeamMemberDropdown {
	var UniqueTeams []entity.UniqueTeamEntity
	var UniqueEmployees []entity.UniqueEmployeeEntity
	var UniqueCustomers []entity.UniqueCustomerEntity

	for _, inputData := range input {
		var UniqueTeam entity.UniqueTeamEntity
		UniqueTeam.TeamId = inputData.TeamId.String
		UniqueTeam.TeamName = inputData.TeamName.String
		UniqueTeams = append(UniqueTeams, UniqueTeam)

		var UniqueEmployee entity.UniqueEmployeeEntity
		UniqueEmployee.TeamId = inputData.TeamId.String
		UniqueEmployee.ActiveDirectory = inputData.ActiveDirectory.String
		UniqueEmployee.ApprovalRoleTitle = inputData.ApproverTitle.String
		UniqueEmployee.ApprovalRoleValue = inputData.ApprovalValue.String
		UniqueEmployee.Email = inputData.Email.String
		UniqueEmployee.FirstName = inputData.FirstName.String
		UniqueEmployee.LastName = inputData.LastName.String
		UniqueEmployee.TeamMemberID = inputData.TeamMemberId.String
		UniqueEmployee.UserId = inputData.UserId.String
		UniqueEmployees = append(UniqueEmployees, UniqueEmployee)

		var UniqueCustomer entity.UniqueCustomerEntity
		UniqueCustomer.TeamMemberId = inputData.TeamMemberId.String
		UniqueCustomer.CustomerId = inputData.CustomerId.String
		UniqueCustomer.CustomerName = inputData.CustomerName.String
		UniqueCustomer.ShipTo = int(inputData.CustomerShipTo.Int64)
		UniqueCustomer.SoldTo = int(inputData.CustomerSoldTo.Int64)
		UniqueCustomer.TeamMemberCustomerId = inputData.TeamMemberCustomerId.String
		UniqueCustomers = append(UniqueCustomers, UniqueCustomer)
	}

	UniqueTeamData := UniqueTeamData(UniqueTeams)
	UniqueEmployeeData := UniqueEmployeeData(UniqueEmployees)
	UniqueCustomerData := UniqueCustomerData(UniqueCustomers)
	var teamResponse []*model.TeamMemberDropdown
	for _, teamValue := range UniqueTeamData {
		var finalTeamOutput model.TeamMemberDropdown

		newTeam := teamValue.TeamId
		newTeamName := teamValue.TeamName
		finalTeamOutput.TeamID = &newTeam
		finalTeamOutput.TeamName = &newTeamName

		var employee []*model.SalesRepDataDropDown
		for _, employeeValue := range UniqueEmployeeData {
			var employeeData model.SalesRepDataDropDown
			if employeeValue.TeamId == teamValue.TeamId {

				activeDirectory := employeeValue.ActiveDirectory
				approvalValue := employeeValue.ApprovalRoleTitle
				approvalTitle := employeeValue.ApprovalRoleValue
				email := employeeValue.Email
				firstName := employeeValue.FirstName
				lastName := employeeValue.LastName
				teamMember := employeeValue.TeamMemberID
				userId := employeeValue.UserId

				employeeData.ActiveDirectory = &activeDirectory
				employeeData.ApprovalRoleTitle = &approvalTitle
				employeeData.ApprovalRoleValues = &approvalValue
				employeeData.Email = &email
				employeeData.FirstName = &firstName
				employeeData.LastName = &lastName
				employeeData.TeamMemberID = &teamMember
				employeeData.UserID = &userId

				var customer []*model.CustomerDataDropDown
				for _, customerValues := range UniqueCustomerData {
					var customerData model.CustomerDataDropDown

					if customerValues.TeamMemberId == employeeValue.TeamMemberID {

						customerId := customerValues.CustomerId
						customerName := customerValues.CustomerName
						customerShipTo := customerValues.ShipTo
						customerSoldTo := customerValues.SoldTo
						customerTeamMemberCustomer := customerValues.TeamMemberCustomerId

						customerData.CustomerID = &customerId
						customerData.CustomerName = &customerName
						customerData.ShipTo = &customerShipTo
						customerData.SoldTo = &customerSoldTo
						customerData.TeamMemberCustomerID = &customerTeamMemberCustomer
						customer = append(customer, &customerData)
					}
				}
				employeeData.Customers = customer
				employee = append(employee, &employeeData)
			}
		}
		finalTeamOutput.Employee = employee
		teamResponse = append(teamResponse, &finalTeamOutput)
	}
	return teamResponse
}

func UniqueTeamData(TeamSlice []entity.UniqueTeamEntity) []entity.UniqueTeamEntity {
	keys := make(map[entity.UniqueTeamEntity]bool)
	uniqueList := []entity.UniqueTeamEntity{}
	for _, entry := range TeamSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func UniqueEmployeeData(EmployeeSlice []entity.UniqueEmployeeEntity) []entity.UniqueEmployeeEntity {
	keys := make(map[entity.UniqueEmployeeEntity]bool)
	uniqueList := []entity.UniqueEmployeeEntity{}
	for _, entry := range EmployeeSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func UniqueCustomerData(CustomerSlice []entity.UniqueCustomerEntity) []entity.UniqueCustomerEntity {
	keys := make(map[entity.UniqueCustomerEntity]bool)
	uniqueList := []entity.UniqueCustomerEntity{}
	for _, entry := range CustomerSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func MapInputForPrinciPalName(input model.PrincipalDropDownInput, salesorg string) (string, error) {
	if input.TeamID == "" {
		return "", errors.New("Team Id Is Blank")
	}
	_, err := uuid.FromString(input.TeamID)
	if err != nil {
		return "", errors.New("Team Id Format Is Invalid")
	}
	_, erro := postgres.HasTeamID(input.TeamID, salesorg)
	if erro != nil {
		return "", errors.New("Team Id Not Found")
	}
	TeamId := input.TeamID
	return TeamId, nil
}

func MapPrincipalNameToEntity(input []entity.PrincipalDropDownOutput) []*model.PrincipalDropDownData {
	outPutModel := []*model.PrincipalDropDownData{}
	for _, value := range input {
		innerModel := model.PrincipalDropDownData{}
		innerModel.PrincipalName = value.PrincipalName
		outPutModel = append(outPutModel, &innerModel)
	}
	return outPutModel
}

func MapTaskBulletinTypeToEntity(input []entity.Titlelues) []*model.TitleValue {
	outPut := []*model.TitleValue{}
	for _, value := range input {
		innerModel := model.TitleValue{}
		innerModel.Type = value.Type
		outPut = append(outPut, &innerModel)
	}
	return outPut
}
