package controller

import (
	"strings"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/mapper"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
	"github.com/eztrade/kpi/graph/util"
)

func FlashBulletinList(inputModel model.ListFlashBulletinInput, loggedInUserEntity *entity.LoggedInUser) model.ListFlashBulletinResponse {
	FlashBulletinResponse := model.ListFlashBulletinResponse{}
	entitymap, ListFlashBulletinResponse := mapper.MapListFlashBulletinToEntity(inputModel)
	FlashBulletinResponse = *ListFlashBulletinResponse
	var flashBulletins []entity.FlashBulletinList
	if ListFlashBulletinResponse.Error {
		FlashBulletinResponse.Error = true
		FlashBulletinResponse.Message = ListFlashBulletinResponse.Message
		return FlashBulletinResponse
	} else {
		if loggedInUserEntity.AuthRole == "sfe" || loggedInUserEntity.AuthRole == "cbm" {
			var err error
			flashBulletins, err = postgres.ListFlashBulletinData(entitymap, loggedInUserEntity)
			if err != nil {
				return FlashBulletinResponse
			}
		} else {
			var err error
			flashBulletins, err = postgres.ListFlashBulletinDataForLineManager(entitymap, loggedInUserEntity)
			if err != nil {
				return FlashBulletinResponse
			}
		}

		for _, flashBulletin := range flashBulletins {

			outputModel := model.FlashBulletinData{}
			outputModel.ID = flashBulletin.ID.String
			outputModel.Title = flashBulletin.Title.String
			outputModel.Description = flashBulletin.Description.String
			outputModel.Status = flashBulletin.Status.Bool
			outputModel.StartDate = flashBulletin.StartDate
			outputModel.EndDate = flashBulletin.EndDate
			outputModel.CreatedDate = util.GetTimeUnixTimeStamp(flashBulletin.CreatedDate)
			if flashBulletin.ModifiedDate != nil {
				outputModel.ModifiedDate = util.GetTimeUnixTimeStamp(*flashBulletin.ModifiedDate)
			} else {
				outputModel.ModifiedDate = "0"
			}
			outputModel.Type = postgres.GetCategoryFromCode(flashBulletin.Type.Int64)
			attachmentIDs := flashBulletin.Attachments.String
			attachmentIDs = attachmentIDs[1:]
			attachmentIDs = attachmentIDs[:len(attachmentIDs)-1]
			attachmentIDList := strings.Split(attachmentIDs, ",")
			blobURLList, err := postgres.GetUploadsByIDList(attachmentIDList)
			if err != nil {
				FlashBulletinResponse.Error = true
				FlashBulletinResponse.Message = err.Error()
				return FlashBulletinResponse
			}
			outputModel.Attachments = blobURLList
			FlashBulletinResponse.FlashBulletins = append(FlashBulletinResponse.FlashBulletins, &outputModel)
		}
	}
	if len(FlashBulletinResponse.FlashBulletins) < 1 {
		FlashBulletinResponse.Error = false
		FlashBulletinResponse.Message = "No Flash Bulletins found"
	}
	return FlashBulletinResponse
}
