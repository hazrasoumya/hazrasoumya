package postgres

import (
	"context"
	"strings"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/jmoiron/sqlx"
)

func GetTaskBulletinReportValues(salesorg string, title []string) ([]entity.ReportData, error) {
	result := []entity.ReportData{}

	var inputArgs []interface{}
	query := `select tb.id as task_buletin_id, tb.title as task_bulletin_title, c.title as type,
	tb.principal_name, c3.name as customer_name, t.team_name as team_name, (case when u.last_name is null then u.first_name else concat(u.first_name, '-',u.last_name)end ) as user_name, 
	u.active_directory,TO_CHAR(tb.creation_date ,'DD/MM/YYYY'),TO_CHAR(tb.target_date ,'DD/MM/YYYY'),  
	c2.title as status, ctf.remarks,
	ctf.date_created , c3.id, ctf.attachments
	from task_bulletin tb
	inner join customer_task_feedback ctf on tb.id = ctf.task_bulletin_id
	inner join team_member_customer tmc on ctf.team_member_customer = tmc.id
	inner join team_members tm on tmc.team_member = tm.id
	inner join team t on tm.team = t.id
	inner join code c on c.id = tb.type
	inner join code c2 on c2.id = ctf.status
	inner join customer c3 on tmc.customer = c3.id
	inner join "user" u on u.id = tm.employee 
	where u.sales_organisation = ? `
	var IDs []string
	if len(title) > 0 {
		for _, tmId := range title {
			IDs = append(IDs, "'"+tmId+"'")
		}
		stringTmIds := strings.Trim(strings.Join(IDs, ", "), ", ")
		query = query + " AND tb.title IN(" + stringTmIds + ") "
	}
	inputArgs = append(inputArgs, salesorg)
	query += `and tb.is_active = true and tb.is_deleted = false 
and tmc.is_active = true and tmc.is_deleted = false 
and tm.is_active = true and tm.is_deleted = false 
and t.is_active = true and t.is_deleted = false 
and c3.is_active = true and c3.is_deleted = false 
and u.is_active = true and u.is_deleted = false 
and c.is_active = true and c.is_delete = false 
and c2.is_active = true and c2.is_delete = false
order by(ctf.date_created)`

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, _ := pool.Query(context.Background(), query, inputArgs...)
	defer rows.Close()
	for rows.Next() {
		res := entity.ReportData{}
		err := rows.Scan(
			&res.BulletinId,
			&res.BulletinTitle,
			&res.BulletinType,
			&res.PrincipalName,
			&res.CustomerName,
			&res.TeamName,
			&res.UserName,
			&res.ActiveDirectory,
			&res.CreationDate,
			&res.TargetDate,
			&res.Status,
			&res.Remarks,
			&res.FeedbackDate,
			&res.CustomerId,
			&res.Attachments,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return result, err
		}
		result = append(result, res)

	}
	return result, nil
}
