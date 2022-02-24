package model

type GetKpiAnswerRequest struct {
	KpiID                *string `json:"kpiId"`
	KpiVersionID         *string `json:"kpiVersionId"`
	Category             *int    `json:"category"`
	TeamMemberCustomerID *string `json:"teamMemberCustomerId"`
	ScheduleEvent        *string `json:"scheduleEvent"`
	TargetItem           *string `json:"targetItem"`
}

type GetKpiAnswerResponse struct {
	Error      		bool         `json:"error"`
	Message    		string       `json:"message"`
	ErrorCode  		int          `json:"errorCode"`
	IsOldAnswer 	bool         `json:"isOldAnswer"`
	IsProposedStock bool         `json:"isProposedStock"`
	GetAnswers 		[]*KpiAnswer `json:"getAnswers"`
}

type KPIAnswerRes struct {
	QuestioNnumber int      `json:"questioNnumber"`
	Value          []string `json:"value"`
}

type KPIAnswerStruct struct {
	QuestioNnumber int      `json:"questioNnumber"`
	Value          []string `json:"value"`
}

type KpiAnswer struct {
	ID                   string         `json:"id"`
	Answer               []KPIAnswerRes `json:"answer"`
	KpiID                string         `json:"kpiId"`
	KpiVersionID         string         `json:"kpiVersionId"`
	Category             int            `json:"category"`
	TeamMemberCustomerID string         `json:"teamMemberCustomerId"`
	ScheduleEvent        string         `json:"scheduleEvent"`
	TargetItem           string         `json:"targetItem"`
}

type KpiAnswerRequest struct {
	KpiVersionID         string             `json:"kpiVersionId"`
	Category             int                `json:"category"`
	Answers              []*KPIAnswerStruct `json:"answers"`
	AuthorActiveDirName  string             `json:"authorActiveDirName"`
	TeamMemberCustomerID string             `json:"teamMemberCustomerId"`
	TargetItem           string             `json:"targetItem"`
	ScheduleEvent        string             `json:"scheduleEvent"`
}
