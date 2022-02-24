package entity

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type FlashBulletinList struct {
	ID           sql.NullString `json:"id"`
	Title        sql.NullString `json:"title"`
	Description  sql.NullString `json:"description"`
	Status       sql.NullBool   `json:"status"`
	Attachments  sql.NullString `json:"attachments"`
	StartDate    string         `json:"startdate"`
	EndDate      string         `json:"enddate"`
	Type         sql.NullInt64  `json:"type"`
	CreatedDate  time.Time      `json:"createdDate"`
	ModifiedDate *time.Time     `json:"modifiedDate"`
}

type FlashBulletinListInput struct {
	Type                 int        `json:"type"`
	TypeValue            string     `json:"typeValue"`
	Status               *bool      `json:"status"`
	StartDate            string     `json:"startDate"`
	EndDate              string     `json:"endDate"`
	ReceipientID         *uuid.UUID `json:"receipientId"`
	TeamMemberCustomerID *uuid.UUID `json:"teamMemberCustomerId"`
	TeamID               *uuid.UUID `json:"teamID"`
	CustomerID           *uuid.UUID `json:"customerID"`
}
