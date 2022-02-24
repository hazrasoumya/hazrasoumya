package postgres

import (
	"context"
	"errors"

	"github.com/eztrade/kpi/graph/logengine"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

func ChannelsExist(ids []string) error {
	if pool == nil {
		pool = GetPool()
	}
	query := `select id from channel 
	where is_active = true AND is_deleted = false
	AND id in (?)`
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
		err = errors.New("Invalid Channels")
	}
	return err
}
