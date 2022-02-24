package postgres

import (
	"context"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	suuid "github.com/gofrs/uuid"
)

func CheckTeamMemberCustomerId(tmcId *suuid.UUID) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	querystring := `select 1 from team_member_customer
	left join team_members on team_member_customer.team_member = team_members.id
	left join team on team_members.team = team.id
	where team_member_customer.id = $1 and team_member_customer.is_active = true and team_member_customer.is_deleted = false`
	var hasValue int
	err := pool.QueryRow(context.Background(), querystring, tmcId).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}

func CheckTeamMemberCustomerID(entity *entity.FlashBulletinListInput) error {

	querystring := `SELECT team.id,customer from team_member_customer
	left join team_members on team_member_customer.team_member = team_members.id
	left join team on team_members.team = team.id
	where team_member_customer.id = $1 and team_member_customer.is_active = true and team_member_customer.is_deleted = false`

	err := pool.QueryRow(context.Background(), querystring, entity.TeamMemberCustomerID).Scan(&entity.TeamID, &entity.CustomerID)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	return err
}
