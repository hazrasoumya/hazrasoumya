package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
)

func GetTeamName(id string) string {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select team_name from team where id=$1 AND is_active = true AND is_deleted = false"
	var name string
	err := pool.QueryRow(context.Background(), querystring, id).Scan(&name)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return name
	}
	return name
}

func GetUserName(id string) string {
	if pool == nil {
		pool = GetPool()
	}
	querystring := `select active_directory from "user" where id = $1 and is_active = true and is_deleted = false`
	var name string
	err := pool.QueryRow(context.Background(), querystring, id).Scan(&name)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return name
	}
	return name
}

func TeamsExistBySalesOrg(ids []string, salesOrgId string) error {
	if pool == nil {
		pool = GetPool()
	}
	query := `select id from team 
	where is_active = true AND is_deleted = false
	AND id in (?) AND sales_organisation = ?`
	var inputArgs []interface{}
	inputArgs = append(inputArgs, ids, salesOrgId)
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
		err = errors.New("Invalid Teams")
	}
	return err
}

func HasTeamId(team string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	var id pgtype.UUID
	querystring := "select id from team where id=$1 and is_active=true and is_deleted=false"
	err := pool.QueryRow(context.Background(), querystring, team).Scan(&id)
	if err == nil {
		result = false
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = true
	}
	return result
}

func GetTeamsByIDs(receipients string) ([]*model.Recipients, error) {
	if pool == nil {
		pool = GetPool()
	}
	receipients = strings.ReplaceAll(receipients, "{", "")
	receipients = strings.ReplaceAll(receipients, "}", "")
	receipientIDs := strings.Split(receipients, ",")
	sqlQuery := `select id, team_name from team where id in (?)`

	var inputArgs []interface{}
	inputArgs = append(inputArgs, receipientIDs)
	sqlQuery, args, err := sqlx.In(sqlQuery, inputArgs...)
	if err != nil {
		return []*model.Recipients{}, err
	}
	sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
	rows, err := pool.Query(context.Background(), sqlQuery, args...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []*model.Recipients{}, err
	}
	defer rows.Close()
	receipientList := make([]*model.Recipients, 0)
	for rows.Next() {
		receipient := model.Recipients{}
		var id, teamName string
		err := rows.Scan(
			&id,
			&teamName,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return []*model.Recipients{}, err
		}
		receipient.ID = id
		receipient.Description = teamName
		receipientList = append(receipientList, &receipient)
	}
	return receipientList, nil
}

func GetTeamIDsBySalesOrg(loggedInUserEntity *entity.LoggedInUser) ([]uuid.UUID, error) {
	var inputArgs []interface{}
	query := `select t.id from team t 
	inner join team_members tm on tm.team = t.id 
	inner join code c on c.id = tm.approval_role and c.category = 'ApprovalRole'
	inner join sales_organisation so on so.id = t.sales_organisation 
	where so.id = ?`
	inputArgs = append(inputArgs, loggedInUserEntity.SalesOrganisaton)
	if loggedInUserEntity.AuthRole == "salesrep" {
		query = query + `and c.value in ('line1manager','line2manager','salesrep')
		and tm.employee = ?`
		inputArgs = append(inputArgs, loggedInUserEntity.ID)
	}
	query = query + `and so.is_active = true and so.is_deleted = false 
	and t.is_active = true and t.is_deleted = false 
	and tm.is_active = true and tm.is_deleted = false
	and c.is_active = true and c.is_delete = false
	and so.is_active = true and so.is_deleted = false`

	sqlQuery := sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := pool.Query(context.Background(), sqlQuery, inputArgs...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []uuid.UUID{}, err
	}
	teams := []uuid.UUID{}
	defer rows.Close()
	for rows.Next() {
		team := entity.TeamIDS{}
		err := rows.Scan(
			&team.TeamID,
		)
		teamUuid, err := uuid.FromString(team.TeamID.String)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return teams, err
		}
		teams = append(teams, teamUuid)
	}
	return teams, err
}

func HasTeam(teamId string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	queryString := `select 1 from team where id = $1 AND is_active = true AND is_deleted = false`
	var hasValue int
	err := pool.QueryRow(context.Background(), queryString, teamId).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}
