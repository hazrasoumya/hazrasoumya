package entity

import (
	"database/sql"

	uuid "github.com/gofrs/uuid"
)

type CustomerTarget struct {
	ID              *uuid.UUID
	TypeID          int
	Category        int
	ProductBraindId string
	Targets         []TargetValue
	CreatedBy       uuid.UUID
	ModifiedBy      uuid.UUID
	IsDeleted       bool
	IsActve         bool
	SalesOrgId      string
	Year            int
}

type TargetCustomerInput struct {
	SalesOrgId       *string
	TargetCustomerId *string
	Type             *string
	Category         *string
	ProductBrandId   *string
	ProductBrandName *string
	Year             *int
	Limit            *int
	Offset           *int
	PageNo           *int
}

type TargetCustomerResponse struct {
	CustomerTargetId sql.NullString `json:"customerTargetId"`
	Type             sql.NullString `json:"type"`
	Category         sql.NullString `json:"category"`
	ProductBrandId   sql.NullString `json:"productBrandId"`
	ProductBrandName sql.NullString `json:"productBrandName"`
	Targets          sql.NullString `json:"targets"`
}

type CustomerGroupResponse struct {
	IndustrialCode sql.NullString `json:"inDusTrialCode"`
	CustomerId     sql.NullString `json:"customerId"`
	CustomerName   sql.NullString `json:"customerName"`
	SoldTo         sql.NullString `json:"soldTo"`
	ShipTo         sql.NullString `json:"shipTo"`
}

type CustomerTargetExcelInterface struct {
	Data []interface{}
}

type UniqueIndustryCode struct {
	IndustryCode sql.NullString `jsaon:"industryCode"`
}

type CustomerLineManager struct {
	TeamId sql.NullString `json:"teamId"`
}
