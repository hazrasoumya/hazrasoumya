package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

func UploadsExist(ids []string) error {
	if pool == nil {
		pool = GetPool()
	}
	query := `select id from uploads 
	where id in (?)`
	var inputArgs []interface{}
	inputArgs = append(inputArgs, ids)
	query, args, err := sqlx.In(query, inputArgs...)
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := pool.Query(context.Background(), query, args...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return err
	}
	defer rows.Close()
	rowCount := 0
	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(
			&id,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return err
		}
		rowCount++
	}
	if len(ids) != rowCount {
		err = errors.New("Invalid Uploads")
	}
	return err
}

func GetUploadsByIDList(ids []string) ([]string, error) {
	var response []string
	var err error
	if pool == nil {
		pool = GetPool()
	}
	query := `select blob_url from uploads 
	where id in (?)`
	var inputArgs []interface{}
	inputArgs = append(inputArgs, ids)
	query, args, err := sqlx.In(query, inputArgs...)
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := pool.Query(context.Background(), query, args...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return response, err
	}
	defer rows.Close()
	for rows.Next() {
		var blobURL string
		err = rows.Scan(
			&blobURL,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return response, err
		}
		response = append(response, blobURL)
	}
	return response, err
}

func GetAttachmentsByIDs(attachments string) ([]*model.Attachment, error) {
	if pool == nil {
		pool = GetPool()
	}
	attachments = strings.ReplaceAll(attachments, "{", "")
	attachments = strings.ReplaceAll(attachments, "}", "")
	attachmentIDs := strings.Split(attachments, ",")
	sqlQuery := `select id, blob_url, blob_name from uploads where id in (?)`

	var inputArgs []interface{}
	inputArgs = append(inputArgs, attachmentIDs)
	sqlQuery, args, err := sqlx.In(sqlQuery, inputArgs...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []*model.Attachment{}, err
	}
	sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
	rows, err := pool.Query(context.Background(), sqlQuery, args...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []*model.Attachment{}, err
	}
	defer rows.Close()
	attachmentList := make([]*model.Attachment, 0)
	for rows.Next() {
		attachment := model.Attachment{}
		err := rows.Scan(
			&attachment.ID,
			&attachment.URL,
			&attachment.Filename,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return []*model.Attachment{}, err
		}
		attachmentList = append(attachmentList, &attachment)
	}
	return attachmentList, nil
}
