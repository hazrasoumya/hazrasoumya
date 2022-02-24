package entity

import (
	"database/sql"

	uuid "github.com/gofrs/uuid"
)

type TaskBulletin struct {
	ID                 *uuid.UUID
	Title              string
	TypeID             int
	TypeValue          string
	CreationDate       string
	TargetDate         string
	TeamMemberCustomer []string
	Description        string
	PrincipalName      string
	Attachments        []Attachment
	IsActive           *bool
	IsDeleted          *bool
	CreatedBy          uuid.UUID
	ModifiedBy         uuid.UUID
}

type TaskBulletinfeedbackInput struct {
	TeamMemberCustomerId uuid.UUID
	TaskBulletinId       uuid.UUID
}

type TaskBulletinFeedbackOutput struct {
	StatusTitle sql.NullString `json:"statusTitle"`
	StatusValue sql.NullString `json:"statusValue"`
	Remarks     sql.NullString `json:"remarks"`
	DateCreated sql.NullString `json:"dateCreated"`
	Attachment  sql.NullString `json:"attachment"`
}

type CustomerTaskFeedBack struct {
	TaskBulletinId       uuid.UUID
	TeamMemberCustomerId uuid.UUID
	Status               int8
	Remarks              string
	Attachments          []Attachmentdata
	IsDeleted            *bool
	IsActive             *bool
	CreatedBy            string
	ModifiedBy           string
}

type Attachmentdata struct {
	Id       string
	BlobName string
	BlobUrl  string
}

type TaskBulletinInput struct {
	Id                   *string
	Type                 *string
	IsActive             bool
	CreationDate         *string
	TargetDate           *string
	TeamMemberId         *string
	TeamMemberCustomerId *string
	PageNo               *int
	Limit                *int
	Offset               *int
	SearchItem           *string
	UserId               *string
	Salesorg             *string
}

type TaskBulletinOutput struct {
	Id                   sql.NullString `json:"id"`
	TypeTitle            sql.NullString `json:"typeTitle"`
	TypeValue            sql.NullString `json:"typeValue"`
	Title                sql.NullString `json:"title"`
	TeamId               sql.NullString `json:"teamId"`
	TeamName             sql.NullString `json:"teamName"`
	Description          sql.NullString `json:"description"`
	PrincipalName        sql.NullString `json:"principalName"`
	CreationDate         sql.NullString `json:"creationDate"`
	TargetDate           sql.NullString `json:"targetDate"`
	BlobId               sql.NullString `json:"blobId"`
	BlobURL              sql.NullString `json:"blobURL"`
	BlobName             sql.NullString `json:"blobName"`
	TeamMemberId         sql.NullString `json:"teamMemberId"`
	UserId               sql.NullString `json:"userId"`
	FirstName            sql.NullString `json:"firstName"`
	LastName             sql.NullString `json:"lastName"`
	ActiveDirectory      sql.NullString `json:"activeDirectory"`
	Email                sql.NullString `json:"email"`
	ApprovalRoleTitle    sql.NullString `json:"approvalRoleTitle"`
	ApprovalRoleValues   sql.NullString `json:"approvalRoleValues"`
	CustomerID           sql.NullString `json:"customerID"`
	TeamMemberCustomerId sql.NullString `json:"teamMemberCustomerId"`
	CustomerName         sql.NullString `json:"customerName"`
	SoldTo               int            `json:"soldTo"`
	ShipTo               int            `json:"shipTo"`
}

type UniqueTeammember struct {
	TaskBulletinId     string
	TeamId             string
	TeamMemberId       string
	UserId             string
	FirstName          string
	LastName           string
	ActiveDirectory    string
	Email              string
	ApprovalRoleTitle  string
	ApprovalRoleValues string
}

type UniqueTaskBulletin struct {
	Id            string
	TypeTitle     string
	TypeValue     string
	Title         string
	TeamId        string
	TeamName      string
	Description   string
	CreationDate  string
	TargetDate    string
	PrincipalName string
}

type UniqueAttachment struct {
	Id       string
	BlobId   string
	Url      string
	BlobName string
}

type Attachments struct {
	BlobId   string
	Url      string
	BlobName string
}

type UniqueCustomers struct {
	TaskBulletinId       string
	TeamMemberId         string
	CusomerId            string
	TeamMemberCustomerID string
	CustomerName         string
	SoldTo               int
	ShipTo               int
}

type TeamToCustomer struct {
	UserId               sql.NullString `json:"UserId"`
	FirstName            sql.NullString `json:"FirstName"`
	LastName             sql.NullString `json:"LastName"`
	ActiveDirectory      sql.NullString `json:"ActiveDirectory"`
	Email                sql.NullString `json:"Email"`
	TeamName             sql.NullString `json:"TeamName"`
	TeamId               sql.NullString `json:"TeamId"`
	TeamMemberCustomerId sql.NullString `json:"TeamMemberCustomerId"`
	TeamMemberId         sql.NullString `json:"TeamMemberId"`
	CustomerId           sql.NullString `json:"CustomerId"`
	ApprovalValue        sql.NullString `json:"ApprovalValue"`
	ApproverTitle        sql.NullString `json:"ApproverTitle"`
	CustomerSoldTo       sql.NullInt64  `json:"CustomerSoldTo"`
	CustomerShipTo       sql.NullInt64  `json:"CustomerShipTo"`
	CustomerName         sql.NullString `json:"CustomerName"`
}

type UniqueTeamEntity struct {
	TeamName string
	TeamId   string
}

type UniqueEmployeeEntity struct {
	TeamId            string
	TeamMemberID      string
	UserId            string
	FirstName         string
	LastName          string
	ActiveDirectory   string
	Email             string
	ApprovalRoleTitle string
	ApprovalRoleValue string
}

type UniqueCustomerEntity struct {
	TeamMemberId         string
	CustomerId           string
	TeamMemberCustomerId string
	CustomerName         string
	SoldTo               int
	ShipTo               int
}

type PrincipalDropDownInput struct {
	TeamId string
}
type PrincipalDropDownOutput struct {
	PrincipalName string
}

type Titlelues struct {
	Type string
}
