package entity

import (
	uuid "github.com/gofrs/uuid"
)

type KpiAnswer struct {
	Answers              []KpiAnswerStruct
	Category             int
	KpiVersionId         *uuid.UUID
	TeamMemberCustomerID *uuid.UUID
	TargetItem           *uuid.UUID
	ScheduleEvent        *uuid.UUID
}

type KpiAnswerStruct struct {
	QuestioNnumber int
	Value          []string
}

type KpiAnswerData struct {
	ID                   string `json:"id"`
	Answer               string `json:"name"`
	KpiId                string `json:"kpi_id"`
	KpiVersionId         string `json:"kpi_version_id"`
	Category             int64  `json:"category"`
	TeamMemberCustomerID string `json:"teamMemberCustomerId"`
	ScheduleEventID      string `json:"scheduleEventID`
	TargetItem           string `json:"targetItem"`
}
