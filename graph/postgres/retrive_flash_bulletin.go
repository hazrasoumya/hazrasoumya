package postgres

import (
	"context"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	"github.com/jmoiron/sqlx"
)

func RetiveFlashBulletinDataById(inputModel model.RetriveInfoFlashBulletinInput) (entity.RetriveFlashBulletin, error) {
	if pool == nil {
		pool = GetPool()
	}
	query := `Select flash_bulletin.id as ID, code.title as Type, flash_bulletin.title as Title, flash_bulletin.description as Description, 
	flash_bulletin.validity_date as Validity_Date, flash_bulletin.attachments as Attachments, 
	case when code.title = 'Customer Group' then flash_bulletin.customer_group else flash_bulletin.recipients::text[] end as  Recipients
	from flash_bulletin flash_bulletin 
	inner join code on flash_bulletin.type = code.id
	where flash_bulletin.is_deleted = false and flash_bulletin.id = ?`
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := pool.Query(context.Background(), query, inputModel.BulletinID)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return entity.RetriveFlashBulletin{}, err
	}
	var flashBulletinRetriveData entity.RetriveFlashBulletin
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&flashBulletinRetriveData.ID,
			&flashBulletinRetriveData.Type,
			&flashBulletinRetriveData.Title,
			&flashBulletinRetriveData.Description,
			&flashBulletinRetriveData.ValidityDate,
			&flashBulletinRetriveData.Attachments,
			&flashBulletinRetriveData.Recipients,
		)
		if err != nil {
			return flashBulletinRetriveData, err
		}
	}
	return flashBulletinRetriveData, err
}

func UserIDChecking(input string, loggedInUserEntity *entity.LoggedInUser) bool {
	if pool == nil {
		pool = GetPool()
	}
	querystring := `
	with team as (
		with team1 as (
		select unnest(fb.recipients) as receipt from flash_bulletin fb
		where fb.id = $1
		and fb.is_deleted = false 
		) select 1 from team1
		inner join team_members tm on tm.team = receipt
		where tm.employee = $2
		and tm.is_active = true and tm.is_deleted = false
		), customer as(
		with customer1 as (
		select unnest(fb.recipients) as receipt from flash_bulletin fb
		where fb.id = $1
		and fb.is_deleted = false 
		) select 1 from customer1
		inner join team_member_customer tmc2 on tmc2.customer = receipt
		inner join team_members tm2 on tm2.id = tmc2.team_member
		where tm2.employee = $2
		and tm2.is_active = true  and tm2.is_deleted = false 
		and tmc2.is_active = true and tmc2.is_deleted = false 
		), salesOrg as (
		select 1 from "user" u 
		inner join flash_bulletin fb on fb.sales_organisation = u.sales_organisation 
		where u.id = $2 and
		fb.is_deleted = false 
		and u.is_active = true  and u.is_deleted = false and 
		fb.id = $1
		)
		select * from team
		union 
		select * from customer
		union
		select * from salesOrg`
	var result bool
	var hasValue int
	err := pool.QueryRow(context.Background(), querystring, input, loggedInUserEntity.ID).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}
