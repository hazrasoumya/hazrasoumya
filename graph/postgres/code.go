package postgres

import (
	"context"
	"errors"

	"github.com/eztrade/kpi/graph/entity"
)

func GetCodeIdForKpi(input string, category string) (*int8, error) {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select id from code where value=$1 and category= $2 AND is_active = true AND is_delete = false"
	var id int8
	err := pool.QueryRow(context.Background(), querystring, input, category).Scan(&id)
	if err != nil {
		return &id, errors.New("Invalid type")
	}
	return &id, err
}

func GetCodeTitleFromValue(input string, category string) (*string, error) {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select title from code where value=$1 and category= $2 AND is_active = true AND is_delete = false"
	var title string
	err := pool.QueryRow(context.Background(), querystring, input, category).Scan(&title)
	if err != nil {
		return &title, errors.New("Invalid type")
	}
	return &title, err
}

func CheckCodeIDForKpi(input int64, category string) bool {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select id from code where id=$1 and category= $2 AND is_active = true AND is_delete = false"
	var id int8
	err := pool.QueryRow(context.Background(), querystring, input, category).Scan(&id)
	if err != nil {
		return false
	}
	return true
}

func GetCodeValueForKpi(input int8, category string) (*string, *string, error) {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select title, value from code where id=$1 and category= $2 AND is_active = true AND is_delete = false"
	var title string
	var value string
	err := pool.QueryRow(context.Background(), querystring, input, category).Scan(&title, &value)
	if err != nil {
		return &title, &value, errors.New("Invalid type")
	}
	return &title, &value, err
}

func GetCategoryFromCode(categoryID int64) string {
	if pool == nil {
		pool = GetPool()
	}
	querystring := "select value from code where id=$1 AND is_active = true AND is_delete = false"
	var value string
	err := pool.QueryRow(context.Background(), querystring, categoryID).Scan(&value)
	if err != nil {
		return "Category not found, category id invalid"
	}
	return value
}

func ValidateFlashBulletinType(typeId int) (string, error) {
	if pool == nil {
		pool = GetPool()
	}
	querystring := `select value as typeValue from code where id = $1 and category = $2 AND is_active = true AND is_delete = false`
	var typeValue string
	err := pool.QueryRow(context.Background(), querystring, typeId, "BulletinType").Scan(&typeValue)
	if err != nil {
		err = errors.New("Invalid Type")
	}
	return typeValue, err
}

func GetKpiTargetTitleInfo(category string) ([]entity.KpiTargetTitleList, error) {
	if pool == nil {
		pool = GetPool()
	}
	kpiTargetTitles := []entity.KpiTargetTitleList{}
	querystring := "select title,value,description from code where is_delete=false and category=$1"
	rows, err := pool.Query(context.Background(), querystring, category)
	if err != nil {
		return kpiTargetTitles, err
	}
	defer rows.Close()
	for rows.Next() {
		kpiTargetTitle := entity.KpiTargetTitleList{}
		err := rows.Scan(
			&kpiTargetTitle.Title,
			&kpiTargetTitle.Value,
			&kpiTargetTitle.Description,
		)
		if err != nil {
			return kpiTargetTitles, err
		}
		kpiTargetTitles = append(kpiTargetTitles, kpiTargetTitle)
	}
	return kpiTargetTitles, nil
}
