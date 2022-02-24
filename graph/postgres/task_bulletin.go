package postgres

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/kpi/graph/util"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

func UpsertTaskBulletin(entity *entity.TaskBulletin, response *model.TaskBulletinUpsertResponse, loggedInUserEntity *entity.LoggedInUser) {
	if pool == nil {
		pool = GetPool()
	}
	tx, err := pool.Begin(context.Background())
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		response.Message = "Failed to begin transaction"
		response.Error = true
	}
	var teamMemberCustomer []string
	for _, value := range entity.TeamMemberCustomer {
		teamMemberCustomer = append(teamMemberCustomer, value)
	}
	defer tx.Rollback(context.Background())
	timenow := util.GetCurrentTime()
	var attachmentIds []string
	for _, attachment := range entity.Attachments {
		var attachmentId string
		if attachment.ID == nil {
			querystring := `INSERT INTO uploads (blob_url, blob_name, created_by, status) VALUES($1, $2, $3, $4) RETURNING(id)`
			err := tx.QueryRow(context.Background(), querystring, attachment.BlobUrl, attachment.BlobName, attachment.CreatedBy, "").Scan(&attachmentId)
			if err != nil {
				logengine.GetTelemetryClient().TrackException(err.Error())
				response.Message = "Failed to insert attachment data"
				response.Error = true
			}
		} else {
			attachmentId = attachment.ID.String()
			updateStmt := `UPDATE uploads SET blob_url = $2, blob_name = $3, modified_by = $4, last_modified = now() WHERE id = $1`
			commandTag, err := tx.Exec(context.Background(), updateStmt, attachmentId, attachment.BlobUrl, attachment.BlobName, attachment.ModifiedBy)
			if err != nil {
				logengine.GetTelemetryClient().TrackException(err.Error())
				response.Message = "Failed to update attachment data"
				response.Error = true
			}
			if commandTag.RowsAffected() != 1 {
				response.Message = "Failed to update attachment data"
				response.Error = true
			}
		}
		attachmentIds = append(attachmentIds, attachmentId)
	}
	if entity.ID == nil {
		querystring := `INSERT INTO task_bulletin (title, type, creation_date, target_date, team_member_customer, principal_name, description, attachments, sales_organisation, is_active, created_by, date_created) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		commandTag, err := tx.Exec(context.Background(), querystring, entity.Title, entity.TypeID, entity.CreationDate, entity.TargetDate, teamMemberCustomer, entity.PrincipalName, entity.Description, attachmentIds, loggedInUserEntity.SalesOrganisaton, true, entity.CreatedBy, timenow)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			response.Message = "Failed to create Task bulletin data"
			response.Error = true
			return
		} else if commandTag.RowsAffected() != 1 {
			response.Message = "Failed to create Task bulletin data"
			response.Error = true
			return
		} else {
			response.Message = "Task Bulletin successfully inserted"
			response.Error = false
		}

	} else if entity.IsDeleted != nil && *entity.IsDeleted {
		querystring := `UPDATE task_bulletin SET is_deleted = true, is_active = false, modified_by = $2, last_modified = $3 WHERE id = $1`
		commandTag, err := pool.Exec(context.Background(), querystring, entity.ID, entity.ModifiedBy, timenow)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			response.Message = "Failed to delete in Task bulletin"
			response.Error = true
			return
		}
		if commandTag.RowsAffected() != 1 {
			response.Message = "Invalid Task bulletin data"
			response.Error = true
			return
		} else {
			response.Message = "Task bulletin successfully deleted"
			response.Error = false
		}
	} else {
		querystring := `UPDATE task_bulletin SET title = $2,  target_date = $3, team_member_customer = $4,  description = $5, attachments = $6, sales_organisation = $7, modified_by = $8, last_modified = $9 WHERE id = $1`
		commandTag, err := tx.Exec(context.Background(), querystring, entity.ID, entity.Title, entity.TargetDate, teamMemberCustomer, entity.Description, attachmentIds, loggedInUserEntity.SalesOrganisaton, entity.ModifiedBy, timenow)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			response.Message = "Failed to update Task bulletin data"
			response.Error = true
			return
		} else if commandTag.RowsAffected() != 1 {
			response.Message = "Failed to update Task bulletin data"
			response.Error = true
			return
		} else {
			response.Message = "Task Bulletin successfully updated"
			response.Error = false
		}
	}
	txErr := tx.Commit(context.Background())
	if txErr != nil {
		response.Message = "Failed to commit Task Bulletin data"
		response.Error = true
	}
}

func AttachmentBelongsToTaskBulletin(taskBulletinId *string, attachmentId uuid.UUID) error {
	if pool == nil {
		pool = GetPool()
	}
	taskBulletinID := ""
	if taskBulletinId != nil {
		taskBulletinID = *taskBulletinId
	}
	var id string
	attachmentIdstring := "{" + util.UUIDV4ToString(attachmentId) + "}"
	querystring := `select id from task_bulletin where id=$1 and is_deleted=false and task_bulletin.attachments @> $2`
	err := pool.QueryRow(context.Background(), querystring, taskBulletinID, attachmentIdstring).Scan(&id)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	return err
}

func ValidateTaskBulletinType(typeValue string) (int, error) {
	if pool == nil {
		pool = GetPool()
	}
	querystring := `select id from code where category='TaskBulletin' and value = $1`
	var typeId int
	err := pool.QueryRow(context.Background(), querystring, typeValue).Scan(&typeId)
	if err != nil {
		err = errors.New("Invalid Type")
	}
	return typeId, err
}

func HasTitle(Title string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	queryString := `select 1 from code where code.title = $1 AND is_active = true AND is_delete = false`
	var hasValue int
	err := pool.QueryRow(context.Background(), queryString, Title).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}

func HasTeamMemberCustomerId(TeamMemberCustomerID string, salesorg string) (bool, error) {
	if pool == nil {
		pool = GetPool()
	}
	var isActive bool
	querystring := `select tmc.is_active from team_member_customer tmc 
	inner join team_members tm on tmc.team_member = tm.id 
	inner join team t on t.id = tm.team 
	where tmc.id = $1
	and t.sales_organisation = $2
	and tmc.is_active = true and tmc.is_deleted = false 
	and tm.is_active = true and tm.is_deleted = false 
	and t.is_active = true and t.is_deleted = false `
	err := pool.QueryRow(context.Background(), querystring, TeamMemberCustomerID, salesorg).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func HasTaskBulletinID(tasbuletinId string) (bool, error) {
	if pool == nil {
		pool = GetPool()
	}
	var isActive bool
	querystring := `select is_active from task_bulletin tb 
	where tb.id = $1
	and tb.is_active = true and tb.is_deleted = false `
	err := pool.QueryRow(context.Background(), querystring, tasbuletinId).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func TaskBulletinFeedbackDetails(input entity.TaskBulletinfeedbackInput) ([]entity.TaskBulletinFeedbackOutput, error) {
	result := []entity.TaskBulletinFeedbackOutput{}
	if pool == nil {
		pool = GetPool()
	}
	var inputArgs []interface{}
	query := `select c.title, c.value , ctf.remarks , to_char(ctf.date_created,'YYYY-MM-DD HH:MI:SS..MS') , ctf.attachments
    from customer_task_feedback ctf 
	inner join code c on c.id = ctf.status 
	where ctf.team_member_customer = ?
	and ctf.task_bulletin_id = ?
	and c.category = 'TaskFeedback'
	and ctf .is_active = true and ctf.is_deleted = false 
	and c.is_active = true and c.is_delete = false  
	order by ctf.date_created desc `

	inputArgs = append(inputArgs, input.TeamMemberCustomerId)
	inputArgs = append(inputArgs, input.TaskBulletinId)
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, _ := pool.Query(context.Background(), query, inputArgs...)
	defer rows.Close()
	for rows.Next() {
		res := entity.TaskBulletinFeedbackOutput{}
		err := rows.Scan(
			&res.StatusTitle,
			&res.StatusValue,
			&res.Remarks,
			&res.DateCreated,
			&res.Attachment,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return result, err
		}
		result = append(result, res)
	}
	return result, nil

}

func InsertCustomerTaskFeedBack(entity *entity.CustomerTaskFeedBack, response *model.CustomerTaskFeedBackResponse, loggedInUserEntity *entity.LoggedInUser) {
	if pool == nil {
		pool = GetPool()
	}
	tx, err := pool.Begin(context.Background())
	if err != nil {
		msg := "Failed to begin transaction"
		response.Message = &msg
		response.Error = true
	}
	defer tx.Rollback(context.Background())
	timenow := util.GetCurrentTime()
	var attachmentIds []string
	if entity.Attachments != nil {
		for _, attachment := range entity.Attachments {
			var attachmentId string
			if &attachment.Id == nil || attachment.Id == "" {
				querystring := `INSERT INTO uploads (blob_url, blob_name, created_by, status) VALUES($1, $2, $3, $4) RETURNING(id)`
				err := tx.QueryRow(context.Background(), querystring, attachment.BlobUrl, attachment.BlobName, loggedInUserEntity.ID, "").Scan(&attachmentId)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					msg := "Failed to insert attachment data"
					response.Message = &msg
					response.Error = true
				}
			} else {
				attachmentId = attachment.Id
				updateStmt := `UPDATE uploads SET blob_url = $2, blob_name = $3, modified_by = $4, last_modified = now() WHERE id = $1`
				commandTag, err := tx.Exec(context.Background(), updateStmt, attachmentId, attachment.BlobUrl, attachment.BlobName, loggedInUserEntity.ID)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					msg := "Failed to update attachment data"
					response.Message = &msg
					response.Error = true
				}
				if commandTag.RowsAffected() != 1 {
					msg := "Failed to update attachment data"
					response.Message = &msg
					response.Error = true
				}
			}
			attachmentIds = append(attachmentIds, attachmentId)
		}
	}

	queryContactString := `INSERT INTO customer_task_feedback (task_bulletin_id, team_member_customer, status, remarks, attachments, is_active, is_deleted, created_by, date_created) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	commandTag, err := tx.Exec(context.Background(), queryContactString, entity.TaskBulletinId, entity.TeamMemberCustomerId, entity.Status, entity.Remarks, attachmentIds, true, false, entity.CreatedBy, timenow)

	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		msg := "Failed to insert in customer task feedback"
		response.Message = &msg
		response.Error = true
	} else if commandTag.RowsAffected() != 1 {
		msg := "Failed to update customer task feedback"
		response.Message = &msg
		response.Error = true

	} else {
		msg := "Customer task feedback successfully inserted"
		response.Message = &msg
		response.Error = false
	}

	txErr := tx.Commit(context.Background())
	if txErr != nil {
		msg := "Failed to commit customer taskfeedback data"
		response.Message = &msg
		response.Error = true
	}
}
func TaskBulletinIdPresent(input string) (*int8, error) {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select id from taskBulletin where value=$1 AND is_active = true AND is_delete = false"
	var id int8
	err := pool.QueryRow(context.Background(), querystring, input).Scan((&id))
	if err != nil {
		return &id, errors.New("Invalid")
	}
	return &id, err
}

func GetCodeIdForStatus(input string, category string) (*int8, error) {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select id from code where value=$1 AND category = $2 AND is_active = true AND is_delete = false"
	var id int8
	err := pool.QueryRow(context.Background(), querystring, input, category).Scan((&id))
	if err != nil {
		return &id, errors.New("Invalid Status")
	}
	return &id, err
}

func TaskBulletinDetails(input entity.TaskBulletinInput) ([]entity.TaskBulletinOutput, int, error) {
	result := []entity.TaskBulletinOutput{}

	dt := time.Now()
	currentDate := (dt.Format("01-02-2006"))
	if pool == nil {
		pool = GetPool()
	}
	var totalRecords int
	limit := 0
	totalPages := 0
	var inputArgs []interface{}
	query := `with listData1 as(select tb.id as task_bulletin_id, c.title,c.value, tb.title as task_bulletin_title, tb.principal_name, t.id as team_id, t.team_name, tb.description, TO_CHAR(tb.creation_date ,'DD/MM/YYYY') as date_created,
	TO_CHAR(tb.target_date ,'DD/MM/YYYY') as target_date, tb.attachments as attachments, tm.id as team_member_id, u.id as user_id, u.first_name, u.last_name,
	u.active_directory, u.email, c2.title as approval_title, c2.value as approval_value, tmc.id as team_member_customer_id, c3.id as customer_id, c3."name", c3.sold_to , c3.ship_to
from task_bulletin tb 
inner join team_member_customer tmc on tmc.id = any(tb.team_member_customer)
inner join team_members tm  on tm.id = tmc.team_member 
inner join team t on tm.team = t.id 
inner join "user" u on u.id = tm.employee 
inner join customer c3 on c3.id = tmc.customer 
inner join code c on c.id = tb."type"
inner join code c2 on tm.approval_role = c2.id
where tb.sales_organisation = ? `

	inputArgs = append(inputArgs, input.Salesorg)
	query += `and tmc.is_active = true and tmc.is_deleted = false 
	and tm.is_active = true and tm.is_deleted = false 
	and t.is_active = true and t.is_deleted = false 
	and u.is_active = true and u.is_deleted = false 
	and c3.is_active = true and c3.is_deleted = false 
	and c2.is_active = true and c2.is_delete = false 
	and c.is_active = true and c.is_delete = false 
	and tb.is_active = true and tb.is_deleted = false`

	if input.Id != nil && *input.Id != "" {
		query += ` and tb.id = ? `
		inputArgs = append(inputArgs, input.Id)
	}
	if input.TeamMemberCustomerId != nil && *input.TeamMemberCustomerId != "" {
		query += ` and tmc.id = ? `
		inputArgs = append(inputArgs, input.TeamMemberCustomerId)
	}
	if input.TeamMemberId != nil && *input.TeamMemberId != "" {
		query += ` and tm.id = ? `
		inputArgs = append(inputArgs, input.TeamMemberId)
	}
	if input.CreationDate != nil && *input.CreationDate != "" {
		query += ` and tb.creation_date >= ? `
		inputArgs = append(inputArgs, *input.CreationDate)
	}
	if input.TargetDate != nil && *input.TargetDate != "" {
		query += ` and tb.target_date <= ? `
		inputArgs = append(inputArgs, *input.TargetDate)
	}
	if input.IsActive {
		query += ` and tb.target_date >= ? `
		inputArgs = append(inputArgs, currentDate)
	}
	if input.Type != nil {
		query += ` and c.value = ? `
		inputArgs = append(inputArgs, input.Type)
	}
	if input.SearchItem != nil && *input.SearchItem != "" {
		query += ` and (c3."name" ilike ?
			or t.team_name ilike ?
			or tb.title ilike ?
			or tb.principal_name ilike ?) `
		inputArgs = append(inputArgs, "%"+*input.SearchItem+"%", "%"+*input.SearchItem+"%", "%"+*input.SearchItem+"%", "%"+*input.SearchItem+"%")

	}

	query += `), listData2 as(select task_bulletin_id, title, value, task_bulletin_title, principal_name,
        team_id, team_name, description, listData1.date_created, listData1.target_date, attachment, u2.blob_url , u2.blob_name, team_member_id,
        user_id, first_name, last_name, active_directory, email, approval_title, approval_value, team_member_customer_id, customer_id, name, sold_to, ship_to
        from listData1 
        left join unnest(listData1.attachments) as attachment on true
        left join uploads u2 on u2.id  = attachment ),
 listData3 as (select distinct(task_bulletin_id) tbId from listData2  `

	if input.Limit != nil {
		limit = *input.Limit
		query += ` limit ? `
		inputArgs = append(inputArgs, *input.Limit)
	}

	if input.Offset != nil {
		query += ` offset ? `
		inputArgs = append(inputArgs, *input.Offset)
	}

	query += `),listData4 as (select count(distinct(task_bulletin_id)) total_count from listData2),
 listData5 as (select task_bulletin_id, title, value, task_bulletin_title, principal_name,
			   team_id, team_name, description, date_created, target_date, attachment, blob_url , blob_name, team_member_id,
			   user_id, first_name, last_name, active_directory, email, approval_title, approval_value, team_member_customer_id, customer_id, name, sold_to, ship_to from listData3 l3 inner join listData2 l2 on l2.task_bulletin_id = l3.tbId), 
listData6 as (select * from listData4 cross join listData5)     
                   select * from listData6 `
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, _ := pool.Query(context.Background(), query, inputArgs...)
	defer rows.Close()
	for rows.Next() {
		res := entity.TaskBulletinOutput{}
		err := rows.Scan(
			&totalRecords,
			&res.Id,
			&res.TypeTitle,
			&res.TypeValue,
			&res.Title,
			&res.PrincipalName,
			&res.TeamId,
			&res.TeamName,
			&res.Description,
			&res.CreationDate,
			&res.TargetDate,
			&res.BlobId,
			&res.BlobURL,
			&res.BlobName,
			&res.TeamMemberId,
			&res.UserId,
			&res.FirstName,
			&res.LastName,
			&res.ActiveDirectory,
			&res.Email,
			&res.ApprovalRoleTitle,
			&res.ApprovalRoleValues,
			&res.TeamMemberCustomerId,
			&res.CustomerID,
			&res.CustomerName,
			&res.SoldTo,
			&res.ShipTo,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return result, 0, err
		}
		result = append(result, res)
	}
	if limit > 0 {
		d := float64(totalRecords) / float64(limit)
		totalPages = int(math.Ceil(d))
	}
	return result, totalPages, nil
}

func IsValidTaskBulletinType(types string) (bool, error) {
	if pool == nil {
		pool = GetPool()
	}
	var isActive bool
	querystring := `select is_active from code where value = $1 and category = 'TaskBulletin'`
	err := pool.QueryRow(context.Background(), querystring, types).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func TeamToCustomerDropdownFetch(input *model.TaskBulletinInput, loggedIn *entity.LoggedInUser) ([]entity.TeamToCustomer, error) {
	if pool == nil {
		pool = GetPool()
	}
	var finalResponse []entity.TeamToCustomer

	var inputArgs []interface{}
	QueryString := `select 
	distinct(u.id) as user_Id,
	u.first_name,
	u.last_name, 
	u.active_directory,
	u.email,
	t.team_name as team_name,
	t.id as team_id,
	tmc.id as team_member_customer,
	tm.id as team_member_id,
	c2.id,
	c.value as approval_value,
	c.title as approval_title,
	c2.sold_to,
	c2.ship_to,
	c2."name" 
	from team t
inner join team_members tm on tm.team = t.id 
inner join code c on c.id = tm.approval_role 
inner join "user" u on u.id = tm.employee 
inner join team_member_customer tmc on tmc.team_member = tm.id 
inner join customer c2 on c2.id = tmc.customer 
where 
t.sales_organisation = ?`
	inputArgs = append(inputArgs, loggedIn.SalesOrganisaton)

	if input != nil {
		if input.TeamID != nil {
			QueryString += ` and t.id in (`
			for key, value := range input.TeamID {
				if key > 0 {
					QueryString += `,`
				}
				QueryString += `?`
				inputArgs = append(inputArgs, *value)
			}
			QueryString += `)`
		}

		if input.TeamMemberCustomerID != nil {
			QueryString += ` and tmc.id = ? `
			inputArgs = append(inputArgs, input.TeamMemberCustomerID)
		}

		if input.CustomerID != nil {
			QueryString += ` and c2.id = ? `
			inputArgs = append(inputArgs, input.CustomerID)
		}

		if input.TeamMemberID != nil {
			QueryString += ` and tm.id = ? `
			inputArgs = append(inputArgs, input.TeamMemberID)
		}

		if input.OnlySalesrep != nil {
			if *input.OnlySalesrep == true {
				QueryString += ` and tm.approval_role = 8 `
			}
		}
	}
	QueryString += `
	and tm.is_active = true 
	and tm.is_deleted = false 
	and c.is_active = true 
	and c.is_delete = false 
	and u.is_active = true 
	and u.is_deleted = false 
	and tmc.is_active = true 
	and tmc.is_deleted = false 
	and c2.is_active = true 
	and c2.is_deleted = false`

	QueryString = sqlx.Rebind(sqlx.DOLLAR, QueryString)
	rows, err := pool.Query(context.Background(), QueryString, inputArgs...)
	if err != nil {
		return finalResponse, err
	}

	for rows.Next() {
		var response entity.TeamToCustomer
		err = rows.Scan(
			&response.UserId,
			&response.FirstName,
			&response.LastName,
			&response.ActiveDirectory,
			&response.Email,
			&response.TeamName,
			&response.TeamId,
			&response.TeamMemberCustomerId,
			&response.TeamMemberId,
			&response.CustomerId,
			&response.ApproverTitle,
			&response.ApprovalValue,
			&response.CustomerSoldTo,
			&response.CustomerShipTo,
			&response.CustomerName,
		)
		finalResponse = append(finalResponse, response)
	}
	return finalResponse, nil
}

func CombiNationExistsForFeedback(TaskBulletinId, TeamMemberCustomerID string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	querystring := `with temp1 as (select unnest(tb.team_member_customer) as team_member_customer_id from task_bulletin tb 
	where tb.id = $1)
	select 1 from temp1 
	where temp1.team_member_customer_id = $2 `
	var hasValue int
	err := pool.QueryRow(context.Background(), querystring, TaskBulletinId, TeamMemberCustomerID).Scan(&hasValue)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	} else {
		result = true
	}
	return result
}
func HasTeamMemberId(TeamMemberID string, salesorg string) (bool, error) {
	if pool == nil {
		pool = GetPool()
	}
	var isActive bool
	querystring := `select tm.is_active
	from team_members tm 
	inner join team t on t.id = tm.team 
	where tm.id = $1 and t.sales_organisation = $2
	and tm.is_active = true and tm.is_deleted = false and t.is_active = true and t.is_deleted = false `
	err := pool.QueryRow(context.Background(), querystring, TeamMemberID, salesorg).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func HasTeamID(TeamID string, salesorg string) (bool, error) {
	if pool == nil {
		pool = GetPool()
	}
	var isActive bool
	querystring := `select t.is_active from team t
	where t.id = $1
	and t.sales_organisation = $2`
	err := pool.QueryRow(context.Background(), querystring, TeamID, salesorg).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func GetPriciPalNames(teamId, salesOrg string) ([]entity.PrincipalDropDownOutput, error) {
	result := []entity.PrincipalDropDownOutput{}
	if pool == nil {
		pool = GetPool()
	}
	query := `select distinct (p.principal_name) 
	from product p 
	inner join team_products tp on tp.material_code = p.id 
	where tp.team = ? and p.sales_organisation = ?
	and p.is_active = true and p.is_deleted = false 
	and tp.is_active = true and tp.is_delete = false `
	var inputArgs []interface{}
	inputArgs = append(inputArgs, teamId)
	inputArgs = append(inputArgs, salesOrg)
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, _ := pool.Query(context.Background(), query, inputArgs...)
	defer rows.Close()
	for rows.Next() {
		res := entity.PrincipalDropDownOutput{}
		err := rows.Scan(
			&res.PrincipalName,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return result, err
		}
		result = append(result, res)
	}
	return result, nil
}

func TypeOfTaskBulletin(team, salesorg string) ([]entity.Titlelues, error) {
	output := []entity.Titlelues{}
	if pool == nil {
		pool = GetPool()
	}
	var inputArgs []interface{}
	query := `select distinct(tb.title) 
	from task_bulletin tb 
	inner join team_member_customer tmc on tmc.id = any(tb.team_member_customer)
	inner join team_members tm on tm.id = tmc.team_member 
	where tb.sales_organisation = ?
	and tb.is_active = true and tb.is_deleted = false `
	inputArgs = append(inputArgs, salesorg)
	if team != "" {
		query += `and tm.team  = ? `
		inputArgs = append(inputArgs, team)
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, _ := pool.Query(context.Background(), query, inputArgs...)
	defer rows.Close()
	for rows.Next() {
		res := entity.Titlelues{}
		err := rows.Scan(
			&res.Type,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return output, err
		}
		output = append(output, res)
	}
	return output, nil

}
