package postgres

import (
	"context"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/kpi/graph/util"
	"github.com/gofrs/uuid"
)

func UpsertFlashBulletin(entity *entity.FlashBulletin, response *model.FlashBulletinUpsertResponse) {
	if pool == nil {
		pool = GetPool()
	}
	tx, err := pool.Begin(context.Background())
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		response.Message = "Failed to begin transaction"
		response.Error = true
	}
	defer tx.Rollback(context.Background())
	timenow := util.GetCurrentTime()
	var attachmentIds []string
	for _, attachment := range entity.Attachments {
		var attachmentId string
		if attachment.ID == nil {
			querystring := "INSERT INTO uploads (blob_url, blob_name, created_by, status) VALUES($1, $2, $3, $4) RETURNING(id)"
			err := tx.QueryRow(context.Background(), querystring, attachment.BlobUrl, attachment.BlobName, attachment.CreatedBy, "").Scan(&attachmentId)
			if err != nil {
				logengine.GetTelemetryClient().TrackException(err.Error())
				response.Message = "Failed to insert attachment data"
				response.Error = true
			}
		} else {
			attachmentId = attachment.ID.String()
			updateStmt := "UPDATE uploads SET blob_url = $2, blob_name = $3, modified_by = $4, last_modified = now() WHERE id = $1"
			commandTag, err := tx.Exec(context.Background(), updateStmt, attachmentId, attachment.BlobUrl, attachment.BlobName, attachment.ModifiedBy)
			if err != nil {
				logengine.GetTelemetryClient().TrackException(err.Error())
				response.Message = "Failed to update attachment data"
				response.Error = true
			}
			if commandTag.RowsAffected() != 1 {
				response.Message = "Failed to update attachment data"
				response.Error = true
			}
		}
		attachmentIds = append(attachmentIds, attachmentId)
	}
	if entity.ID == nil {
		if entity.IsActive == nil {
			if len(entity.CustomerGroup) > 0 {
				querystring := "INSERT INTO flash_bulletin (type, title, description, validity_date, attachments, recipients, is_active, created_by, sales_organisation, customer_group) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
				commandTag, err := tx.Exec(context.Background(), querystring, entity.TypeID, entity.Title, entity.Description, entity.ValidityDate, attachmentIds, entity.Recipients, true, entity.CreatedBy, entity.SalesOrgId, entity.CustomerGroup)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					response.Message = "Failed to create flash bulletin data"
					response.Error = true
					return
				} else if commandTag.RowsAffected() != 1 {
					response.Message = "Failed to create flash bulletin data"
					response.Error = true
					return
				} else {
					response.Message = "Flash Bulletin successfully inserted"
					response.Error = false
				}

			} else {
				querystring := "INSERT INTO flash_bulletin (type, title, description, validity_date, attachments, recipients, is_active, created_by, sales_organisation) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)"
				commandTag, err := tx.Exec(context.Background(), querystring, entity.TypeID, entity.Title, entity.Description, entity.ValidityDate, attachmentIds, entity.Recipients, true, entity.CreatedBy, entity.SalesOrgId)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					response.Message = "Failed to create flash bulletin data"
					response.Error = true
					return
				} else if commandTag.RowsAffected() != 1 {
					response.Message = "Failed to create flash bulletin data"
					response.Error = true
					return
				} else {
					response.Message = "Flash Bulletin successfully inserted"
					response.Error = false
				}
			}
		} else {
			if len(entity.CustomerGroup) > 0 {
				querystring := "INSERT INTO flash_bulletin (type, title, description, validity_date, attachments, recipients, is_active, created_by, sales_organisation, customer_group) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
				commandTag, err := tx.Exec(context.Background(), querystring, entity.TypeID, entity.Title, entity.Description, entity.ValidityDate, attachmentIds, entity.Recipients, entity.IsActive, entity.CreatedBy, entity.SalesOrgId, entity.CustomerGroup)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					response.Message = "Failed to create flash bulletin data"
					response.Error = true
					return
				} else if commandTag.RowsAffected() != 1 {
					response.Message = "Failed to create flash bulletin data"
					response.Error = true
					return
				} else {
					response.Message = "Flash Bulletin successfully inserted"
					response.Error = false
				}

			} else {
				querystring := "INSERT INTO flash_bulletin (type, title, description, validity_date, attachments, recipients, is_active, created_by, sales_organisation) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)"
				commandTag, err := tx.Exec(context.Background(), querystring, entity.TypeID, entity.Title, entity.Description, entity.ValidityDate, attachmentIds, entity.Recipients, entity.IsActive, entity.CreatedBy, entity.SalesOrgId)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					response.Message = "Failed to create flash bulletin data"
					response.Error = true
					return
				} else if commandTag.RowsAffected() != 1 {
					response.Message = "Failed to create flash bulletin data"
					response.Error = true
					return
				} else {
					response.Message = "Flash Bulletin successfully inserted"
					response.Error = false
				}
			}
		}

	} else if entity.IsDeleted != nil && *entity.IsDeleted {
		querystring := "UPDATE flash_bulletin SET is_deleted = true, is_active = false, modified_by = $2, last_modified = $3 WHERE id = $1"
		commandTag, err := pool.Exec(context.Background(), querystring, entity.ID, entity.ModifiedBy, timenow)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			response.Message = "Failed to delete in flash bulletin"
			response.Error = true
			return
		}
		if commandTag.RowsAffected() != 1 {
			response.Message = "Invalid flash bulletin data"
			response.Error = true
			return
		} else {
			response.Message = "Flash bulletin successfully deleted"
			response.Error = false
		}
	} else {
		if entity.IsActive == nil {
			if len(entity.CustomerGroup) > 0 {
				querystring := "UPDATE flash_bulletin SET type=$2, title=$3, description=$4, validity_date=$5, attachments=$6, recipients=$7, modified_by=$8, last_modified=$9, customer_group = $10 WHERE id=$1"
				commandTag, err := tx.Exec(context.Background(), querystring, entity.ID, entity.TypeID, entity.Title, entity.Description, entity.ValidityDate, attachmentIds, entity.Recipients, entity.ModifiedBy, timenow, entity.CustomerGroup)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					response.Message = "Failed to update flash bulletin data"
					response.Error = true
					return
				} else if commandTag.RowsAffected() != 1 {
					response.Message = "Failed to update flash bulletin data"
					response.Error = true
					return
				} else {
					response.Message = "Flash Bulletin successfully updated"
					response.Error = false
				}

			} else {
				querystring := "UPDATE flash_bulletin SET type=$2, title=$3, description=$4, validity_date=$5, attachments=$6, recipients=$7, modified_by=$8, last_modified=$9 WHERE id=$1"
				commandTag, err := tx.Exec(context.Background(), querystring, entity.ID, entity.TypeID, entity.Title, entity.Description, entity.ValidityDate, attachmentIds, entity.Recipients, entity.ModifiedBy, timenow)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					response.Message = "Failed to update flash bulletin data"
					response.Error = true
					return
				} else if commandTag.RowsAffected() != 1 {
					response.Message = "Failed to update flash bulletin data"
					response.Error = true
					return
				} else {
					response.Message = "Flash Bulletin successfully updated"
					response.Error = false
				}
			}
		} else {
			if len(entity.CustomerGroup) > 0 {
				querystring := "UPDATE flash_bulletin SET type=$2, title=$3, description=$4, validity_date=$5, attachments=$6, recipients=$7, modified_by=$8, last_modified=$9, is_active=$10, customer_group= $11 WHERE id=$1"
				commandTag, err := tx.Exec(context.Background(), querystring, entity.ID, entity.TypeID, entity.Title, entity.Description, entity.ValidityDate, attachmentIds, entity.Recipients, entity.ModifiedBy, timenow, entity.IsActive, entity.CustomerGroup)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					response.Message = "Failed to update flash bulletin data"
					response.Error = true
					return
				} else if commandTag.RowsAffected() != 1 {
					response.Message = "Failed to update flash bulletin data"
					response.Error = true
					return
				} else {
					response.Message = "Flash Bulletin successfully updated"
					response.Error = false
				}

			} else {
				querystring := "UPDATE flash_bulletin SET type=$2, title=$3, description=$4, validity_date=$5, attachments=$6, recipients=$7, modified_by=$8, last_modified=$9, is_active=$10 WHERE id=$1"
				commandTag, err := tx.Exec(context.Background(), querystring, entity.ID, entity.TypeID, entity.Title, entity.Description, entity.ValidityDate, attachmentIds, entity.Recipients, entity.ModifiedBy, timenow, entity.IsActive)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					response.Message = "Failed to update flash bulletin data"
					response.Error = true
					return
				} else if commandTag.RowsAffected() != 1 {
					response.Message = "Failed to update flash bulletin data"
					response.Error = true
					return
				} else {
					response.Message = "Flash Bulletin successfully updated"
					response.Error = false
				}
			}
		}
	}
	txErr := tx.Commit(context.Background())
	if txErr != nil {
		response.Message = "Failed to commit Flash Bulletin data"
		response.Error = true
	}
}

func AttachementBelongsToFlashBulletin(flashBulletinId *string, attachemtntId uuid.UUID) error {
	if pool == nil {
		pool = GetPool()
	}
	flashBulletinID := ""
	if flashBulletinId != nil {
		flashBulletinID = *flashBulletinId
	}
	var id string
	attachmentIdstring := "{" + util.UUIDV4ToString(attachemtntId) + "}"
	querystring := `select id from flash_bulletin where id=$1 and is_deleted=false and flash_bulletin.attachments @> $2`
	err := pool.QueryRow(context.Background(), querystring, flashBulletinID, attachmentIdstring).Scan(&id)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	return err
}

func IsLineOneManager(userID string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	querystring := `select distinct 1 from team t 
	inner join team_members tm on tm.team = t.id
	inner join code c on c.id = tm.approval_role
	inner join "user" u2 on u2.id = tm.employee
	where t.is_active = true and t.is_deleted = false 
	and tm.is_active = true and tm.is_deleted = false 
	and c.is_active = true and c.is_delete = false 
	and c.value = 'line1manager' 
	and c.category = 'ApprovalRole'
	and tm.employee = $1 
	and t.sales_organisation = u2.sales_organisation`
	var hasValue int
	err := pool.QueryRow(context.Background(), querystring, userID).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}
