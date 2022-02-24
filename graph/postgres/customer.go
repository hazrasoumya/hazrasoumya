package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

func GetCustomerName(id string) string {
	var name string
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select name from customer where id=$1 AND is_active = true AND is_deleted = false"
	err := pool.QueryRow(context.Background(), querystring, id).Scan(&name)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return ""
	}
	return name
}

func CustomersExistBySalesOrg(ids []string, salesOrgId string) error {
	if pool == nil {
		pool = GetPool()
	}
	query := `select id from customer 
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
		err = errors.New("Invalid Customers")
	}
	return err
}

func CustomersExist(ids []string) error {
	if pool == nil {
		pool = GetPool()
	}
	query := `select id from customer 
	where is_active = true AND is_deleted = false
	AND id in (?)`
	var inputArgs []interface{}
	inputArgs = append(inputArgs, ids)
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
		err = errors.New("Invalid Customers")
	}
	return err
}

func GetCustomersByIDs(receipients string) ([]*model.Recipients, error) {
	if pool == nil {
		pool = GetPool()
	}
	receipients = strings.ReplaceAll(receipients, "{", "")
	receipients = strings.ReplaceAll(receipients, "}", "")
	receipientIDs := strings.Split(receipients, ",")
	sqlQuery := `select id, sold_to, ship_to, name from customer where id in (?)`

	var inputArgs []interface{}
	inputArgs = append(inputArgs, receipientIDs)
	sqlQuery, args, err := sqlx.In(sqlQuery, inputArgs...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
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
		var id, customerName string
		var soldTo, shipTo int
		err := rows.Scan(
			&id,
			&soldTo,
			&shipTo,
			&customerName,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return []*model.Recipients{}, err
		}
		receipient.ID = id
		receipient.Description = fmt.Sprintf("%d,%d: %s", soldTo, shipTo, customerName)
		receipientList = append(receipientList, &receipient)
	}
	return receipientList, nil
}

func GetCustomerGroups(receipients string) ([]*model.Recipients, error) {
	if pool == nil {
		pool = GetPool()
	}
	receipients = strings.ReplaceAll(receipients, "{", "")
	receipients = strings.ReplaceAll(receipients, "}", "")
	receipientIDs := strings.Split(receipients, ",")
	sqlQuery := `select distinct(c.industrycode2) from customer c where c.industrycode2 in (?)`

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
		var id string
		err := rows.Scan(
			&id,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return []*model.Recipients{}, err
		}
		receipient.ID = id
		receipientList = append(receipientList, &receipient)
	}
	return receipientList, nil
}

func GetCountryByID(salesorg_id string) (string, error) {
	if pool == nil {
		pool = GetPool()
	}
	var con string
	query := `SELECT country FROM sales_organisation WHERE id = $1`
	err := pool.QueryRow(context.Background(), query, salesorg_id).Scan(&con)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return "", err
	}
	return con, nil
}
