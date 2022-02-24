package model

type KPITargetInput struct {
	ID                  *string       `json:"id"`
	SalesRepID          *string       `json:"salesRepId"`
	TeamID              *string       `json:"teamId"`
	Year                int           `json:"year"`
	Target              []*KPITargets `json:"target"`
	AuthorActiveDirName string        `json:"authorActiveDirName"`
}

type KPITargets struct {
	KpiTitle string          `json:"kpiTitle"`
	Values   []*TargetValues `json:"values"`
}

type TargetValueRes struct {
	Month int64   `json:"month"`
	Value float64 `json:"value"`
}

type GetKpiTargetRequest struct {
	ID         *string `json:"id"`
	Year       *int    `json:"year"`
	Status     *string `json:"status"`
	SalesRepID *string `json:"salesRepId"`
	TeamID     *string `json:"teamId"`
}

type GetKpiTargetResponse struct {
	Error      bool         `json:"error"`
	Message    string       `json:"message"`
	GetTargets []*KpiTarget `json:"getTargets"`
}

type KpiTarget struct {
	ID       string         `json:"id"`
	Year     int64          `json:"year"`
	Region   string         `json:"region"`
	Country  string         `json:"country"`
	Currency string         `json:"currency"`
	Plants   int64          `json:"plants"`
	Bergu    string         `json:"bergu"`
	Status   string         `json:"status"`
	TeamName string         `json:"teamName"`
	SalesRep string         `json:"salesRep"`
	Target   []KpiTargetRes `json:"target"`
}

type KpiTargetRes struct {
	KpiTitle string           `json:"kpiTitle"`
	KpiValue string           `json:"kpiValue"`
	Values   []TargetValueRes `json:"values"`
}

type TargetValues struct {
	Month int64   `json:"month"`
	Value float64 `json:"value"`
}

type ActionKPITargetInput struct {
	ID                  string `json:"id"`
	Action              bool   `json:"action"`
	AuthorActiveDirName string `json:"authorActiveDirName"`
}
