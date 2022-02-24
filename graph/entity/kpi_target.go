package entity

import (
	"database/sql"
	"strconv"

	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/kpi/graph/util"
	uuid "github.com/gofrs/uuid"
)

type KpiTargetInput struct {
	ID         *uuid.UUID
	SalesRepID *uuid.UUID
	TeamID     *uuid.UUID
	Year       int
	Status     *int8
	Target     []KpiTarget
}

type KpiTarget struct {
	KpiTitle *int8         `json:"kpiTitle"`
	Values   []TargetValue `json:"values"`
}

type TargetValue struct {
	Month int64   `json:"month"`
	Value float64 `json:"value"`
}

type KpiTargetData struct {
	ID       sql.NullString `json:"id"`
	Year     sql.NullInt64  `json:"year"`
	SalesRep sql.NullString `json:"salesrep"`
	TeamName sql.NullString `json:"teamname"`
	Region   sql.NullString `json:"region"`
	Country  sql.NullString `json:"country"`
	Currency sql.NullString `json:"currency"`
	Plants   sql.NullInt64  `json:"plants"`
	Bergu    sql.NullString `json:"bergu"`
	Status   sql.NullInt64  `json:"status"`
	Target   sql.NullString `json:"target"`
}

type ActionKPITarget struct {
	ID     uuid.UUID
	Status *int8
}

type KpiTargetTitleList struct {
	Title       sql.NullString `json:"title"`
	Value       sql.NullString `json:"value"`
	Description sql.NullString `json:"description"`
}

type CallPlanData struct {
	Month sql.NullInt64 `json:"month"`
	Value sql.NullInt64 `json:"value"`
}

func (c *KpiTargetInput) ValidateTargetData(row int, result *model.ValidationResult) {
	validationMessages := []*model.ValidationMessage{}
	if c.Year == 0 {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ":  Year can not be blank!"}
		validationMessages = append(validationMessages, errorMessage)
	}
	if c.Year < 0 {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ": Year is not valid"}
		validationMessages = append(validationMessages, errorMessage)
	}
	if util.GetIntegerLength(c.Year) > 4 || util.GetIntegerLength(c.Year) < 4 {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ": Length of year is invalid"}
		validationMessages = append(validationMessages, errorMessage)
	}
	if len(validationMessages) > 0 {
		if !result.Error {
			result.Error = true
			result.ValidationMessage = []*model.ValidationMessage{}
		}
		result.ValidationMessage = append(result.ValidationMessage, validationMessages...)
	}
}
