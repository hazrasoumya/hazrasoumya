package postgres

import (
	"context"

	"github.com/eztrade/kpi/graph/logengine"
	"github.com/jackc/pgtype"
)

func HasProductId(product string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	var id pgtype.UUID
	querystring := "select id from product where id=$1 and is_active=true and is_deleted=false"
	err := pool.QueryRow(context.Background(), querystring, product).Scan(&id)
	if err == nil {
		result = false
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = true
	}
	return result
}

func HasTeamProductId(teamProductID string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	var id pgtype.UUID
	querystring := "select id from team_products where id= $1 and is_active=true and is_delete=false"
	err := pool.QueryRow(context.Background(), querystring, teamProductID).Scan(&id)
	if err == nil {
		result = false
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = true
	}
	return result
}

func HasTeamProduct(teamProductId string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	queryString := `select 1 from team_products where id = $1 AND is_active = true AND is_delete = false`
	var hasValue int
	err := pool.QueryRow(context.Background(), queryString, teamProductId).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}
