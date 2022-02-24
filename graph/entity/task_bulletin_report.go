package entity

import (
	"database/sql"
	"time"
)

type ReportData struct {
	BulletinId      sql.NullString `json:"bulletinId"`
	BulletinTitle   sql.NullString `json:"bulletinTitle"`
	BulletinType    sql.NullString `json:"bulletinType"`
	PrincipalName   sql.NullString `json:"principalName"`
	CustomerId      sql.NullString `json:"customerId"`
	CustomerName    sql.NullString `json:"customerName"`
	TeamName        sql.NullString `json:"teamName"`
	UserName        sql.NullString `json:"userName"`
	ActiveDirectory sql.NullString `json:"activeDirectory"`
	CreationDate    sql.NullString `json:"creationDate"`
	TargetDate      sql.NullString `json:"targetDate"`
	Status          sql.NullString `json:"status"`
	Remarks         sql.NullString `json:"remarks"`
	Attachments     sql.NullString `json:"attachments"`
	FeedbackDate    time.Time
	WeekNumber      int `json:"weeknumber"`
	WeekDateValue   string
}

type LatestFeedBack struct {
	BulletinId string
	WeekNumber int
	Customer   string
}

type UniqueTaskBulletinReport struct {
	BulletinId      string
	BulletinTitle   string
	BulletinType    string
	PrincipalName   string
	CustomerId      string
	CustomerName    string
	TeamName        string
	UserName        string
	ActiveDirectory string
	CreationDate    string
	TargetDate      string
}

type UniqueFeedBack struct {
	BulletinId  string
	WeekNumber  int
	WeekDate    string
	Customer    string
	Status      string
	Remark      string
	Attachments string
}
