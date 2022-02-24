package postgres

import (
	"context"

	"github.com/eztrade/kpi/graph/logengine"
	"github.com/jackc/pgtype"
)

func HasBrandId(brand string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	var id pgtype.UUID
	querystring := "select id from brand where id=$1 and is_active=true and is_deleted=false"
	err := pool.QueryRow(context.Background(), querystring, brand).Scan(&id)
	if err == nil {
		result = false
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = true
	}
	return result
}

func HasBrand(brandId string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	queryString := `select 1 from brand where id = $1 AND is_active = true AND is_deleted = false`
	var hasValue int
	err := pool.QueryRow(context.Background(), queryString, brandId).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}
