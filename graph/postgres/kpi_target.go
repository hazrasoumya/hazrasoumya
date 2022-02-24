package postgres

import (
	"context"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/kpi/graph/util"
	"github.com/gofrs/uuid"
	suuid "github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
)

func HasKpiTargetId(ID *suuid.UUID) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	queryString := "select 1 from kpi_targets where id = $1"
	var hasValue int
	err := pool.QueryRow(context.Background(), queryString, ID).Scan(&hasValue)
	if err == nil {
		result = false
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = true
	}
	return result
}

func ValidateYear(year int, id *uuid.UUID, salesRepId *uuid.UUID, teamId *uuid.UUID) (bool, string) {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	var message = ""
	var hasValue int
	queryString := ""
	if id == nil {
		if salesRepId != nil {
			queryString = "select 1 from kpi_targets where year = $1 AND sales_rep_id = $2"
			err := pool.QueryRow(context.Background(), queryString, year, salesRepId).Scan(&hasValue)
			if err == nil {
				result = true
				message = "Sales Representative"
			} else {
				logengine.GetTelemetryClient().TrackException(err.Error())
				result = false
				message = ""
			}
		} else {
			queryString = "select 1 from kpi_targets where year = $1 AND team_id = $2"
			err := pool.QueryRow(context.Background(), queryString, year, teamId).Scan(&hasValue)
			if err == nil {
				result = true
				message = "Team"
			} else {
				logengine.GetTelemetryClient().TrackException(err.Error())
				result = false
				message = ""
			}
		}
	} else {
		if salesRepId != nil {
			queryString = "select 1 from kpi_targets where year = $1 AND sales_rep_id = $2 AND id <> $3"
			err := pool.QueryRow(context.Background(), queryString, year, salesRepId, id).Scan(&hasValue)
			if err == nil {
				result = true
				message = "Sales Representative"
			} else {
				logengine.GetTelemetryClient().TrackException(err.Error())
				result = false
				message = ""
			}
		} else {
			queryString = "select 1 from kpi_targets where year = $1 AND team_id = $2 AND id <> $3"
			err := pool.QueryRow(context.Background(), queryString, year, teamId, id).Scan(&hasValue)
			if err == nil {
				result = true
				message = "Team"
			} else {
				logengine.GetTelemetryClient().TrackException(err.Error())
				result = false
				message = ""
			}
		}
	}
	return result, message
}

func UpsertKpiTarget(entity *entity.KpiTargetInput, response *model.KpiResponse, loggedInUserEntity *entity.LoggedInUser) {
	tx, err := pool.Begin(context.Background())
	if err != nil {
		response.Message = "Failed to begin transaction"
		response.Error = true
	}
	defer tx.Rollback(context.Background())
	timenow := util.GetCurrentTime()
	if entity.ID == nil {
		querystring := "INSERT INTO kpi_targets (year, status, target, sales_organisation, created_by, date_created, sales_rep_id, team_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8)"
		commandTag, err := tx.Exec(context.Background(), querystring, entity.Year, entity.Status, entity.Target, loggedInUserEntity.SalesOrganisaton, loggedInUserEntity.ID, timenow, entity.SalesRepID, entity.TeamID)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			response.Message = "Failed to insert Kpi target data"
			response.Error = true
			return
		}
		if commandTag.RowsAffected() != 1 {
			response.Message = "Invalid Kpi target data"
			response.Error = true
			return
		} else {
			response.Message = "Kpi target successfully inserted"
			response.Error = false
		}
	} else {
		querystring := "UPDATE kpi_targets SET year = $2, status = $3, target = $4, last_modified = $5, modified_by = $6, sales_rep_id = $7, team_id = $8 WHERE id = $1"
		commandTag, err := tx.Exec(context.Background(), querystring, entity.ID, entity.Year, entity.Status, entity.Target, timenow, loggedInUserEntity.ID, entity.SalesRepID, entity.TeamID)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			response.Message = "Failed to update Kpi target data"
			response.Error = true
			return
		}
		if commandTag.RowsAffected() != 1 {
			response.Message = "Invalid kpi target data"
			response.Error = true
			return
		} else {
			response.Message = "kpi target successfully updated"
			response.Error = false
		}
	}
	txErr := tx.Commit(context.Background())
	if txErr != nil {
		response.Message = "Failed to commit customer data"
		response.Error = true
	}
}

func GetKpiTargetDetails(id *string, teamId *string, salesRepId *string, status *string, salesOrganisaton string, year int, typeData string) ([]entity.KpiTargetData, error) {
	if pool == nil {
		pool = GetPool()
	}
	query := `select kt.id, kt.year as kt_year, u.active_directory, 
		(select team_name from team where id = ?) team_name, 
		u.region, so.country, so.currency, so.plants, so.bergu, kt.status, kt.target 
		from kpi_targets kt 
		inner join "user" u on u.id = kt.sales_rep_id
		inner join sales_organisation so on so.id = u.sales_organisation
		where u.is_active = true 
		and u.is_deleted = false
		and so.is_active = true 
		and so.is_deleted = false `

	var rows pgx.Rows
	var err error
	var inputArgs []interface{}

	if teamId != nil && *teamId != "" {
		inputArgs = append(inputArgs, teamId)
	} else {
		inputArgs = append(inputArgs, nil)
	}

	if typeData == "filter" {
		query = query + ` and kt.sales_organisation = ? and kt.year = ? `
		inputArgs = append(inputArgs, salesOrganisaton)
		inputArgs = append(inputArgs, year)
	}

	if id != nil && *id != "" {
		query = query + ` and kt.id = ?`
		inputArgs = append(inputArgs, *id)
	}
	if salesRepId != nil && *salesRepId != "" {
		query = query + ` and kt.sales_rep_id = ?`
		inputArgs = append(inputArgs, *salesRepId)
	}
	if teamId != nil && *teamId != "" {
		query = query + ` and kt.sales_rep_id in (select tm.employee from team_members tm 
			inner join team t on t.id = tm.team 
			inner join code c on c.id = tm.approval_role 
			where c.value = 'salesrep' and c.category = 'ApprovalRole'
			and tm.is_active = true and tm.is_deleted = false 
			and t.is_active = true and t.is_deleted = false
			and t.id = ?)`
		inputArgs = append(inputArgs, *teamId)
	}
	if status != nil && *status != "" {
		statusCode, err := GetCodeIdForKpi(*status, "KPITargetStatus")
		if err != nil {
			return []entity.KpiTargetData{}, err
		}
		query = query + ` and kt.status = ?`
		inputArgs = append(inputArgs, *statusCode)
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err = pool.Query(context.Background(), query, inputArgs...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.KpiTargetData{}, err
	}

	kpiTargets := []entity.KpiTargetData{}
	defer rows.Close()
	for rows.Next() {
		kpiTrg := entity.KpiTargetData{}
		err := rows.Scan(
			&kpiTrg.ID,
			&kpiTrg.Year,
			&kpiTrg.SalesRep,
			&kpiTrg.TeamName,
			&kpiTrg.Region,
			&kpiTrg.Country,
			&kpiTrg.Currency,
			&kpiTrg.Plants,
			&kpiTrg.Bergu,
			&kpiTrg.Status,
			&kpiTrg.Target,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return kpiTargets, err
		}
		kpiTargets = append(kpiTargets, kpiTrg)
	}

	return kpiTargets, err
}

func ActionKPITargetData(entity *entity.ActionKPITarget, response *model.KpiResponse, loggedInUserEntity *entity.LoggedInUser) {
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

	querystring := "UPDATE kpi_targets SET status=$2, modified_by=$3, last_modified=$4 WHERE id=$1 AND status <> $2"
	commandTag, err := tx.Exec(context.Background(), querystring, entity.ID, entity.Status, loggedInUserEntity.ID, timenow)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		response.Message = "Failed to take action against KPI target"
		response.Error = true
	} else if commandTag.RowsAffected() != 1 {
		response.Message = "Failed to take action against KPI target"
		response.Error = true
	} else {
		response.Message = "KPI Target action taken successfully"
		response.Error = false
	}
	txErr := tx.Commit(context.Background())
	if txErr != nil {
		response.Message = "Failed to commit kpi target data"
		response.Error = true
	}
}

func ActualWorkingDays(rtId string, year int, dataType string) ([]entity.CallPlanData, error) {
	if pool == nil {
		pool = GetPool()
	}
	query := `with tempData as (
		select 
			 sp."month" as call_month, 
			 row_number() over (partition by date(se.event_date) order by sp."month" desc) as row_number 
		from schedule_plan sp 
			inner join code c1 on c1.id = sp.status 
			inner join schedule_event se on se.schedule_plan_id = sp.id 
			inner join code c2 on c2.id = se."type" 
			inner join team_members tm on tm.id = sp.team_member 
		where c1.value = 'approved' and c1.category = 'ScheduleStatus'
			and sp.is_active = true and sp.is_deleted = false
			and c1.is_active = true and c1.is_delete = false 
			and se.is_active = true and se.is_deleted = false 
			and tm.is_active = true and tm.is_deleted = false 
			and c2.value = 'completed' and c2.category = 'ScheduleEventType'
			and se.check_in is not null and se.check_out is not null 
			and se.team_member_customer is not null
			and sp."year" = ? `

	var inputArgs []interface{}
	inputArgs = append(inputArgs, year)

	if dataType == "team" {
		query += ` and tm.team = ?`
		inputArgs = append(inputArgs, rtId)
	} else if dataType == "representative" {
		query += ` and tm.employee = ?`
		inputArgs = append(inputArgs, rtId)
	}

	query += ` ) select call_month, count(call_month) from tempData where row_number = 1 group by call_month order by call_month`

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := pool.Query(context.Background(), query, inputArgs...)
	dbIDs := []entity.CallPlanData{}
	defer rows.Close()
	for rows.Next() {
		dbID := entity.CallPlanData{}
		err := rows.Scan(
			&dbID.Month,
			&dbID.Value,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return dbIDs, err
		}
		dbIDs = append(dbIDs, dbID)
	}

	return dbIDs, err
}
