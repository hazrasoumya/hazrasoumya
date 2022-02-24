package postgres

import (
	"context"

	"github.com/eztrade/kpi/graph/logengine"
	suuid "github.com/gofrs/uuid"
)

func CheckScheduleEventId(seId *suuid.UUID) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	querystring := `select 1 from schedule_event where id = $1 and is_active = true and is_deleted = false`
	var hasValue int
	err := pool.QueryRow(context.Background(), querystring, seId).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}
