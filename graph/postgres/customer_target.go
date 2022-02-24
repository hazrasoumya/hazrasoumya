package postgres

import (
	"context"
	"math"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/kpi/graph/util"
)

func IdPresent(id string, salesorg string) bool {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select 1 from customer_target where id = $1 AND sales_organisation =$2 AND is_active = true AND is_deleted = false"
	var value int8
	_ = pool.QueryRow(context.Background(), querystring, id, salesorg).Scan(&value)
	if value == 1 {
		return true
	}
	return false
}

func ValidProductId(id string, salesorg string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var value int
	querystring := "select 1 from product where id = $1 AND sales_organisation =$2 AND is_active = true AND is_deleted = false"
	_ = pool.QueryRow(context.Background(), querystring, id, salesorg).Scan(&value)
	if value == 1 {
		return true
	}
	return false
}

func ValidBrandId(id string, salesorg string) bool {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select 1 from brand where id=$1 AND sales_organisation =$2 AND is_active = true AND is_deleted = false"
	var ids int8
	_ = pool.QueryRow(context.Background(), querystring, id, salesorg).Scan(&ids)
	if ids == 1 {
		return true
	}
	return false
}

func ValidBrandProductValue(types string, category string) bool {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select 1 from code where value=$1 and category= $2 AND is_active = true AND is_delete = false"
	var value int8
	_ = pool.QueryRow(context.Background(), querystring, types, category).Scan(&value)
	if value == 1 {
		return true
	}
	return false
}

func CustomerTarget(entity *entity.CustomerTarget, kpiResponse *model.KpiResponse, loggedInUserId string) {
	tx, err := pool.Begin(context.Background())
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		kpiResponse.Message = "Failed to begin transaction"
		kpiResponse.Error = true
		return
	}
	if pool == nil {
		pool = GetPool()
	}
	timenow := util.GetCurrentTime()
	defer tx.Rollback(context.Background())
	if entity.ID == nil {
		queryString := `insert into customer_target ("type", category, "year", targets, product_brand_id, sales_organisation, is_active, is_deleted, created_by, date_created) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
		commandTag, err := tx.Exec(context.Background(), queryString, entity.TypeID, entity.Category, entity.Year, entity.Targets, entity.ProductBraindId, entity.SalesOrgId, true, false, loggedInUserId, timenow)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			kpiResponse.Message = "Failed to insert Customer target data"
			kpiResponse.Error = true
			return
		}
		if commandTag.RowsAffected() != 1 {
			kpiResponse.Message = "Invalid Customer target data"
			kpiResponse.Error = true
			return
		} else {
			kpiResponse.Message = "Customer target successfully inserted"
			kpiResponse.Error = false
		}
	} else {
		queryString := `update customer_target set targets = $1, last_modified = $2, modified_by = $3 `

		if entity.IsDeleted == true {
			queryString += `, is_deleted= true, is_active = false `
		}
		queryString += `where id = $4`
		commandTag, err := tx.Exec(context.Background(), queryString, entity.Targets, timenow, loggedInUserId, entity.ID)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			kpiResponse.Message = "Failed to Update Customer target data"
			kpiResponse.Error = true
			return
		}
		if commandTag.RowsAffected() != 1 {
			kpiResponse.Message = "Invalid Customer target data"
			kpiResponse.Error = true
			return
		} else {
			kpiResponse.Message = "Customer target successfully Updated"
			kpiResponse.Error = false
		}
	}

	txErr := tx.Commit(context.Background())
	if txErr != nil {
		kpiResponse.Message = "Failed to commit customer data"
		kpiResponse.Error = true
	}

}

func AlreadyExists(input *entity.CustomerTarget, salesorg string) bool {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select 1 from customer_target where type =$1 and category = $2 and year =$3 and sales_organisation = $4 and product_brand_id = $5 and is_active = true and is_deleted = false"
	var value int8
	_ = pool.QueryRow(context.Background(), querystring, input.TypeID, input.Category, input.Year, salesorg, input.ProductBraindId).Scan(&value)
	if value == 1 {
		return true
	}
	return false
}

func IsCombinationExist(input *entity.CustomerTarget, salesorg string, inputId *string) bool {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select 1 from customer_target where type =$1 and category = $2 and year =$3 and sales_organisation = $4 and product_brand_id = $5 and id = $6 and is_active = true and is_deleted = false"
	var value int8
	_ = pool.QueryRow(context.Background(), querystring, input.TypeID, input.Category, input.Year, salesorg, input.ProductBraindId, inputId).Scan(&value)
	if value == 1 {
		return true
	}
	return false
}

func HasTargetCustomerId(customerTarget string) (bool, error) {
	if pool == nil {
		pool = GetPool()
	}

	var isActive bool
	querystring := "select is_active from customer_target ct where ct.id = $1 and ct.is_active = true and ct.is_deleted = false"
	err := pool.QueryRow(context.Background(), querystring, customerTarget).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func HasProductBrandId(productBrand string, targetType string) (bool, error) {
	if pool == nil {
		pool = GetPool()
	}

	var isActive bool
	var querystring string

	if strings.EqualFold(targetType, "targetbrand") {
		querystring = "select is_active from brand b where b.id = $1 and b.is_active = true and b.is_deleted = false"
	} else {
		querystring = "select is_active from product p where p.id = $1 and p.is_active = true and p.is_deleted = false"
	}

	err := pool.QueryRow(context.Background(), querystring, productBrand).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func GetTargetCustomer(input entity.TargetCustomerInput) ([]entity.TargetCustomerResponse, int, error) {
	var entities []entity.TargetCustomerResponse
	totalPages := 0
	limit := 0
	var totalRecords int
	var query string
	var inputArgs []interface{}

	if strings.EqualFold(*input.Type, "targetbrand") {
		query = `select
				ct.id customer_target_id,
				c1.value "type",
				c2.value category,
				b.id brand_id,
				b.brand_name,
				ct.targets,
				count(b.id) over()
			from
				brand b
			cross join (select id from code c where c.category = 'MerchandisingType') as category
			left join customer_target ct on
				b.id = ct.product_brand_id
				and category.id = ct.category
				and ct."year" = ?
				and ct.is_active = true
				and ct.is_deleted = false
			left join code c1 on
				c1.id = (select c3.id from code c3 where c3.value = 'targetbrand' and c3.category = 'TargetCustomerType')
			left join code c2 on
				c2.id = category.id
			where
				b.sales_organisation = ?`
		inputArgs = append(inputArgs, *input.Year, *input.SalesOrgId)

		if input.TargetCustomerId != nil && *input.TargetCustomerId != "" {
			query += ` and ct.id = ?`
			inputArgs = append(inputArgs, *input.TargetCustomerId)
		}

		if input.Category != nil && *input.Category != "" {
			query += ` and c2.value = ?`
			inputArgs = append(inputArgs, *input.Category)
		}

		if input.ProductBrandId != nil && *input.ProductBrandId != "" {
			query += ` and b.id = ?`
			inputArgs = append(inputArgs, *input.ProductBrandId)
		}

		if input.ProductBrandName != nil && *input.ProductBrandName != "" {
			query += ` and b.brand_name ilike ?`
			inputArgs = append(inputArgs, "%"+*input.ProductBrandName+"%")
		}

		query += ` and b.is_active = true
				and b.is_deleted = false
			group by
				ct.id,
				b.id,
				c1.value,
				c2.value
			order by
				b.id,
				c2.value`

		if input.Limit != nil {
			limit = *input.Limit

			query = query + ` limit ?`
			inputArgs = append(inputArgs, *input.Limit)
		}

		if input.Offset != nil {
			query = query + ` offset ?`
			inputArgs = append(inputArgs, *input.Offset)
		}
	} else {
		query = `select
				ct.id customer_target_id,
				c1.value "type",
				c2.value category,
				p.id product_id,
				concat(p.principal_name, '|', p.material_description) product_name,
				ct.targets,
				count(p.id) over()
			from
				product p
			cross join (select id from code c where c.category = 'MerchandisingType' and c.value in('distribution', 'promotionexecution', 'posmexecution')) as category
			left join customer_target ct on
				p.id = ct.product_brand_id
				and category.id = ct.category
				and ct."year" = ?
				and ct.is_active = true
				and ct.is_deleted = false
			left join code c1 on
				c1.id = (select c3.id from code c3 where c3.value = 'targetproduct' and c3.category = 'TargetCustomerType')
			left join code c2 on
				c2.id = category.id
			where
				p.sales_organisation = ?`
		inputArgs = append(inputArgs, *input.Year, *input.SalesOrgId)

		if input.TargetCustomerId != nil && *input.TargetCustomerId != "" {
			query += ` and ct.id = ?`
			inputArgs = append(inputArgs, *input.TargetCustomerId)
		}

		if input.Category != nil && *input.Category != "" {
			query += ` and c2.value = ?`
			inputArgs = append(inputArgs, *input.Category)
		}

		if input.ProductBrandId != nil && *input.ProductBrandId != "" {
			query += ` and p.id = ?`
			inputArgs = append(inputArgs, *input.ProductBrandId)
		}

		if input.ProductBrandName != nil && *input.ProductBrandName != "" {
			query += ` and (p.principal_name ilike ?
				or p.material_description ilike ?)`
			inputArgs = append(inputArgs, "%"+*input.ProductBrandName+"%", "%"+*input.ProductBrandName+"%")
		}

		query += ` and p.is_active = true
				and p.is_deleted = false
			group by
				ct.id,
				p.id,
				c1.value,
				c2.value
			order by
				p.id,
				c2.value`

		if input.Limit != nil {
			limit = *input.Limit

			query = query + ` limit ?`
			inputArgs = append(inputArgs, *input.Limit)
		}

		if input.Offset != nil {
			query = query + ` offset ?`
			inputArgs = append(inputArgs, *input.Offset)
		}
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	var rows pgx.Rows
	var err error

	rows, err = pool.Query(context.Background(), query, inputArgs...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.TargetCustomerResponse{}, 0, err
	}

	defer rows.Close()
	for rows.Next() {
		brandEntity := entity.TargetCustomerResponse{}
		err := rows.Scan(
			&brandEntity.CustomerTargetId,
			&brandEntity.Type,
			&brandEntity.Category,
			&brandEntity.ProductBrandId,
			&brandEntity.ProductBrandName,
			&brandEntity.Targets,
			&totalRecords,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return []entity.TargetCustomerResponse{}, 0, err
		}
		entities = append(entities, brandEntity)
	}

	if limit > 0 {
		d := float64(totalRecords) / float64(limit)
		totalPages = int(math.Ceil(d))
	}

	return entities, totalPages, nil
}

func GetValueAndCategory(types string, category string) int {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select id from code where value = $1 and category = $2 and is_active = true and is_delete = false"
	var value int
	err := pool.QueryRow(context.Background(), querystring, types, category).Scan(&value)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	return value
}

func GetCustomerGroup(userEntity *entity.LoggedInUser, data []string, input *model.CustomerGroupInput) ([]entity.CustomerGroupResponse, error) {
	var entities []entity.CustomerGroupResponse
	var query string
	var inputArgs []interface{}
	query = `select 
	        distinct (c.id),
			c.industrycode2,
			c.name,
			c.sold_to,
			c.ship_to 
		from 
			customer c 
		inner join 
			team_member_customer tmc on tmc.customer = c.id 
		inner join 
			team_members tm on tm.id = tmc.team_member
		inner join 
			team t on t.id = tm.team
		where t.sales_organisation = ?
		and c.sales_organisation = ?`
	inputArgs = append(inputArgs, userEntity.SalesOrganisaton, userEntity.SalesOrganisaton)
	if len(data) > 0 {
		query = query + ` and t.id in (`
		for key, value := range data {
			if key == 0 {
				query = query + `'` + value + `'`
			} else {
				query = query + `, '` + value + `'`
			}
		}
		query = query + `)`
	}
	if input != nil {
		if len(input.CustomerGroup) > 0 {
			query = query + ` and c.industrycode2 in (`
			for key, value := range input.CustomerGroup {
				if key == 0 {
					query = query + `'` + *value + `'`
				} else {
					query = query + `, '` + *value + `'`
				}
			}
			query = query + `)`

		}
	}
	query += `and c.is_active = true and c.is_deleted = false 
	and tm.is_active = true and tm.is_deleted = false 
	and t.is_active = true and t.is_deleted = false`
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	var rows pgx.Rows
	var err error

	rows, err = pool.Query(context.Background(), query, inputArgs...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.CustomerGroupResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		brandEntity := entity.CustomerGroupResponse{}
		err := rows.Scan(
			&brandEntity.CustomerId,
			&brandEntity.IndustrialCode,
			&brandEntity.CustomerName,
			&brandEntity.SoldTo,
			&brandEntity.ShipTo,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return []entity.CustomerGroupResponse{}, err
		}
		entities = append(entities, brandEntity)
	}
	return entities, nil
}

func GetLineOneManager(userId string) ([]entity.CustomerLineManager, error) {
	var entities []entity.CustomerLineManager
	var query string
	var inputArgs []interface{}
	query = `select 
			tm.team 
		from team_members tm 
		inner join "user" u on tm.employee = u.id 
		where tm.approval_role = (select c.id from code c where c.value = 'line1manager' and c.category = 'ApprovalRole')
		and u.id = ?
		and u.is_active = true and u.is_deleted = false 
    	and tm.is_active = true and tm.is_deleted = false`
	inputArgs = append(inputArgs, userId)

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	var rows pgx.Rows
	var err error

	rows, err = pool.Query(context.Background(), query, inputArgs...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.CustomerLineManager{}, err
	}

	defer rows.Close()
	for rows.Next() {
		brandEntity := entity.CustomerLineManager{}
		err := rows.Scan(
			&brandEntity.TeamId,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return []entity.CustomerLineManager{}, err
		}
		entities = append(entities, brandEntity)
	}

	return entities, nil
}
