package mapper

import (
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/kpi/graph/util"
	uuid "github.com/gofrs/uuid"
)

func MapRetrivePlanEntityToModel(inputEntity entity.RetriveFlashBulletin, recipients []*model.Recipients, attachments []*model.Attachment) (model.FlashBulletin, error) {
	var err error
	var outputModel model.FlashBulletin
	bulletinID, err := uuid.FromString(inputEntity.ID.String)
	if err != nil {
		return model.FlashBulletin{}, err
	}
	outputModel.ID = util.UUIDV4ToString(bulletinID)
	outputModel.Type = inputEntity.Type.String
	outputModel.Title = inputEntity.Title.String
	outputModel.Description = inputEntity.Description.String
	outputModel.Attachments = attachments
	outputModel.Recipients = recipients
	outputModel.ValidityDate = inputEntity.ValidityDate.String

	return outputModel, err
}
