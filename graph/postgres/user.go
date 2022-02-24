package postgres

import (
	"context"

	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/util"
	suuid "github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
)

// Not In Use
func GetIDByActiveDirectoryName(activeDirName string) (*suuid.UUID, error) {
	querystring := `select id from "user" where active_directory = $1 AND is_active = true AND is_deleted = false`
	var uuid pgtype.UUID
	var userID suuid.UUID
	err := pool.QueryRow(context.Background(), querystring, activeDirName).Scan(&uuid)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return &userID, err
	} else {
		if uuid.Status == pgtype.Present {
			mybyte := uuid.Get().([16]byte)
			userID = util.BytesToUUIDV4(mybyte)
		}
	}
	return &userID, err
}

func IsSalesRep(userId string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	var id pgtype.UUID
	querystring := `select u.id from "user" u inner join code c on c.id = u.authorisation_role  
		where u.id=$1 and c.value = 'salesrep'
		and u.is_active=true and u.is_deleted=false 
		and c.is_active=true and c.is_delete=false`
	err := pool.QueryRow(context.Background(), querystring, userId).Scan(&id)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}
