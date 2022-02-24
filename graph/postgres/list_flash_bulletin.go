package postgres

import (
	"context"
	"database/sql"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/util"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
)

func ListFlashBulletinData(entityInput *entity.FlashBulletinListInput, loggedInUserEntity *entity.LoggedInUser) ([]entity.FlashBulletinList, error) {
	if pool == nil {
		pool = GetPool()
	}
	query := `select flash_bulletin.id, flash_bulletin.title , flash_bulletin.description , flash_bulletin.is_active, flash_bulletin.attachments, flash_bulletin.validity_date, 
	flash_bulletin.type, flash_bulletin.date_created, flash_bulletin.last_modified from flash_bulletin where flash_bulletin.is_deleted = false and flash_bulletin.sales_organisation =?`
	var rows pgx.Rows
	var err error
	var inputArgs []interface{}
	inputArgs = append(inputArgs, loggedInUserEntity.SalesOrganisaton)
	if entityInput.Type > 0 {
		query = query + ` and flash_bulletin.type = ? `
		inputArgs = append(inputArgs, entityInput.Type)
	}

	if entityInput.StartDate != "" && entityInput.EndDate != "" {
		dateRange := "[" + entityInput.StartDate + "," + entityInput.EndDate + ")"
		query = query + ` and flash_bulletin.validity_date <@ ?`
		inputArgs = append(inputArgs, dateRange)
	}
	if entityInput.Status != nil {
		query = query + ` and flash_bulletin.is_active = ? `
		isActive := entityInput.Status
		inputArgs = append(inputArgs, isActive)
	}
	if entityInput.TeamMemberCustomerID != nil {
		switch entityInput.TypeValue {
		case "customer":
			query = query + ` and flash_bulletin.recipients @> ? `
			customerID := "{" + util.UUIDV4ToString(*entityInput.CustomerID) + "}"
			inputArgs = append(inputArgs, customerID)
		case "team":
			query = query + ` and flash_bulletin.recipients @> ? `
			teamID := "{" + util.UUIDV4ToString(*entityInput.TeamID) + "}"
			inputArgs = append(inputArgs, teamID)
		default:
			query = query + ` and (flash_bulletin.recipients @> ? `
			teamId := "{" + util.UUIDV4ToString(*entityInput.TeamID) + "}"
			inputArgs = append(inputArgs, teamId)
			query = query + ` or flash_bulletin.recipients @> ? `
			customerID := "{" + util.UUIDV4ToString(*entityInput.CustomerID) + "}"
			inputArgs = append(inputArgs, customerID)
			query = query + ` ) `
		}
	}
	if entityInput.ReceipientID != nil {
		query = query + ` and flash_bulletin.recipients @> ? `
		recepientID := util.UUIDV4ToString(*entityInput.ReceipientID)
		recepientID = "{" + recepientID + "}"
		inputArgs = append(inputArgs, recepientID)
	}
	if loggedInUserEntity.AuthRole == "cbm" {
		query = query + ` and flash_bulletin.created_by = ? `
		inputArgs = append(inputArgs, loggedInUserEntity.ID)
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err = pool.Query(context.Background(), query, inputArgs...)
	flashBulletinLists := make([]entity.FlashBulletinList, 0)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.FlashBulletinList{}, err
	}

	defer rows.Close()
	for rows.Next() {

		var dateStr sql.NullString
		flashBulletinList := entity.FlashBulletinList{}
		err = rows.Scan(
			&flashBulletinList.ID,
			&flashBulletinList.Title,
			&flashBulletinList.Description,
			&flashBulletinList.Status,
			&flashBulletinList.Attachments,
			&dateStr,
			&flashBulletinList.Type,
			&flashBulletinList.CreatedDate,
			&flashBulletinList.ModifiedDate,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return flashBulletinLists, err
		}
		flashBulletinList.StartDate, flashBulletinList.EndDate = util.GetStartDateEndDate(dateStr.String)
		flashBulletinLists = append(flashBulletinLists, flashBulletinList)
	}

	return flashBulletinLists, err
}

func ListFlashBulletinDataForLineManager(entityInput *entity.FlashBulletinListInput, loggedInUserEntity *entity.LoggedInUser) ([]entity.FlashBulletinList, error) {
	if pool == nil {
		pool = GetPool()
	}
	query := `
	select id,title,description,active,attachment,validity_date,fbtype,date_created,last_modified from (

	select * from (
with teams as (
	with team1 as (
	select unnest(fb.recipients) as receipt, fb.id as id, fb.title as title , fb.description as description, fb.is_active as active, 
	fb.attachments as attachment, fb.validity_date as validity_date , fb."type" as fbtype, fb.date_created as date_created, fb.last_modified as last_modified,
	fb.recipients as recipients, fb.sales_organisation as sales, fb.is_deleted as is_deleted
	from flash_bulletin fb
	where fb.is_deleted = false
	and fb."type" = (select id from code c2 where c2.category='BulletinType' and c2.value = 'team')
	) select distinct(td.id), td.title, td.description, td.active, td.attachment, td.validity_date, td.fbtype, td.date_created,
	td.last_modified, td.recipients, td.sales, td.is_deleted from team1 td
	inner join team_members tm on tm.team = receipt
	where tm.employee = ?
	and tm.is_active = true and tm.is_deleted = false
	), 
	customer as(
	with customer1 as (
	select unnest(fb.recipients) as receipt, fb.id as id, fb.title as title , fb.description as description, fb.is_active as active, 
	fb.attachments as attachment, fb.validity_date as validity_date , fb."type" as fbtype, fb.date_created as date_created, fb.last_modified as last_modified,
	fb.recipients as recipients, fb.sales_organisation as sales, fb.is_deleted as is_deleted
	from flash_bulletin fb
	where fb.is_active = true and fb.is_deleted = false
	and fb."type" = (select id from code c2 where c2.category='BulletinType' and c2.value = 'customer')
	) select distinct(cm.id),cm.title, cm.description, cm.active, cm.attachment, cm.validity_date, cm.fbtype, cm.date_created,
	cm.last_modified, cm.recipients ,cm.sales, cm.is_deleted from customer1 cm
	inner join team_member_customer tmc2 on tmc2.customer = receipt
	inner join team_members tm2 on tm2.id = tmc2.team_member
	where tm2.employee = ?
	and tm2.is_active = true  and tm2.is_deleted = false 
	and tmc2.is_active = true and tmc2.is_deleted = false 
	),customer_group as(
        with customer1 as (
        select unnest(fb.recipients) as receipt, fb.id as id, fb.title as title , fb.description as description, fb.is_active as active,
        fb.attachments as attachment, fb.validity_date as validity_date , fb."type" as fbtype, fb.date_created as date_created, fb.last_modified as last_modified,  
        fb.recipients as recipients, fb.sales_organisation as sales, fb.is_deleted as is_deleted
        from flash_bulletin fb
        where fb.is_active = true and fb.is_deleted = false
        and fb."type" = (select id from code c2 where c2.category='BulletinType' and c2.value = 'customergroup')
        ) select distinct(cm.id),cm.title, cm.description, cm.active, cm.attachment, cm.validity_date, cm.fbtype, cm.date_created,
        cm.last_modified, cm.recipients ,cm.sales, cm.is_deleted from customer1 cm
        inner join team_member_customer tmc2 on tmc2.customer = receipt
        inner join team_members tm2 on tm2.id = tmc2.team_member
        where tm2.employee = ?
        and tm2.is_active = true  and tm2.is_deleted = false
        and tmc2.is_active = true and tmc2.is_deleted = false
        ),salesOrg as (
	with sales as(
	select unnest(fb.recipients) as receipt, fb.id as id, fb.title as title , fb.description as description, fb.is_active as active, 
	fb.attachments as attachment, fb.validity_date as validity_date , fb."type" as fbtype, fb.date_created as date_created, fb.last_modified as last_modified,
	fb.recipients as recipients, fb.sales_organisation as sales, fb.is_deleted as is_deleted
	from flash_bulletin fb
	where fb.is_active = true and fb.is_deleted = false
	and fb."type" = (select id from code c2 where c2.category='BulletinType' and c2.value = 'salesorganisation')
	)
	select distinct(so.id), so.title , so.description, so.active, 
	so.attachment, so.validity_date , so.fbtype, so.date_created, so.last_modified, so.recipients, so.sales, so.is_deleted
	from sales so
	inner join team tm on tm.sales_organisation = so.sales 
	where tm.is_active = true and tm.is_deleted = false
	and tm.id in (select tm3.team from team_members tm3 where tm3.employee = ?)
	)
	select * from teams
	union 
	select * from customer
	union
	select * from salesOrg
	union
	select * from customer_group ) as uniondata
	where sales = ? and is_deleted = false`
	var rows pgx.Rows
	var err error
	var inputArgs []interface{}
	inputArgs = append(inputArgs, loggedInUserEntity.ID, loggedInUserEntity.ID, loggedInUserEntity.ID, loggedInUserEntity.ID, loggedInUserEntity.SalesOrganisaton)
	if entityInput.Type > 0 {
		query = query + ` and fbtype = ? `
		inputArgs = append(inputArgs, entityInput.Type)
	}
	if entityInput.StartDate != "" && entityInput.EndDate != "" {
		dateRange := "[" + entityInput.StartDate + "," + entityInput.EndDate + ")"
		query = query + ` and validity_date <@ ?`
		inputArgs = append(inputArgs, dateRange)
	}
	if entityInput.Status != nil {
		query = query + ` and active = ? `
		isActive := entityInput.Status
		inputArgs = append(inputArgs, isActive)
	}
	if entityInput.TeamMemberCustomerID != nil {
		switch entityInput.TypeValue {
		case "customer":
			query = query + ` and recipients @> ? `
			customerID := "{" + util.UUIDV4ToString(*entityInput.CustomerID) + "}"
			inputArgs = append(inputArgs, customerID)
		case "team":
			query = query + ` and recipients @> ? `
			teamID := "{" + util.UUIDV4ToString(*entityInput.TeamID) + "}"
			inputArgs = append(inputArgs, teamID)
		default:
			query = query + ` and (recipients @> ? `
			teamId := "{" + util.UUIDV4ToString(*entityInput.TeamID) + "}"
			inputArgs = append(inputArgs, teamId)
			query = query + ` or recipients @> ? `
			customerID := "{" + util.UUIDV4ToString(*entityInput.CustomerID) + "}"
			inputArgs = append(inputArgs, customerID)
			query = query + ` ) `
		}
	}
	if entityInput.ReceipientID != nil {
		query = query + ` and recipients @> ? `
		recepientID := util.UUIDV4ToString(*entityInput.ReceipientID)
		recepientID = "{" + recepientID + "}"
		inputArgs = append(inputArgs, recepientID)
	}
	query = query + ` ) as flashbulletindata`
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err = pool.Query(context.Background(), query, inputArgs...)
	flashBulletinLists := make([]entity.FlashBulletinList, 0)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.FlashBulletinList{}, err
	}

	defer rows.Close()
	for rows.Next() {

		var dateStr sql.NullString
		flashBulletinList := entity.FlashBulletinList{}
		err = rows.Scan(
			&flashBulletinList.ID,
			&flashBulletinList.Title,
			&flashBulletinList.Description,
			&flashBulletinList.Status,
			&flashBulletinList.Attachments,
			&dateStr,
			&flashBulletinList.Type,
			&flashBulletinList.CreatedDate,
			&flashBulletinList.ModifiedDate,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return flashBulletinLists, err
		}
		flashBulletinList.StartDate, flashBulletinList.EndDate = util.GetStartDateEndDate(dateStr.String)
		flashBulletinLists = append(flashBulletinLists, flashBulletinList)
	}

	return flashBulletinLists, err
}
