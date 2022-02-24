package model

type UpsertKpiRequest struct {
	ID                  *string  `json:"id"`
	Name                string   `json:"name"`
	TargetTeam          string   `json:"targetTeam"`
	TargetProducts      []string `json:"targetProducts"`
	TargetBrand         []string `json:"targetBrand"`
	Type                string   `json:"type"`
	EffectiveMonth      int      `json:"effectiveMonth"`
	EffectiveYear       int      `json:"effectiveYear"`
	IsPriority          bool     `json:"isPriority"`
	BrandDesign         []*KPI   `json:"brandDesign"`
	ProductDesign       []*KPI   `json:"productDesign"`
	AuthorActiveDirName string   `json:"authorActiveDirName"`
	IsDeleted           bool     `json:"isDeleted"`
}

type KPI struct {
	Name               *string        `json:"name"`
	Category           int64          `json:"category"`
	Type               *string        `json:"type"`
	Active             bool           `json:"active"`
	EffectiveStartDate string         `json:"effectiveStartDate"`
	EffectiveEndDate   string         `json:"effectiveEndDate"`
	Questions          []*KPIQuestion `json:"questions"`
}

type KPIDesign struct {
	Name               *string        `json:"name"`
	Active             *bool          `json:"active"`
	Category           int64          `json:"category"`
	Type               *string        `json:"type"`
	EffectiveStartDate *string        `json:"effectiveStartDate"`
	EffectiveEndDate   *string        `json:"effectiveEndDate"`
	Questions          []*KPIQuestion `json:"questions"`
}

type KPIQuestion struct {
	QuestionNumber int      `json:"questionNumber"`
	Title          string   `json:"title"`
	Type           string   `json:"type"`
	OptionValues   []string `json:"optionValues"`
	Active         *bool    `json:"active"`
	Required       *bool    `json:"required"`
}

type KpiDeleteRequest struct {
	ID                  string `json:"id"`
	AuthorActiveDirName string `json:"authorActiveDirName"`
}

type GetKpi struct {
	ParentKpiID				string			`json:"parentKpiID"`
	ProductKpiID			string			`json:"productKpiID"`
	BrandKpiID				string			`json:"brandKpiID"`
	ProductKpiVersionID		string			`json:"productKpiVersionID"`
	BrandKpiVersionID		string			`json:"brandKpiVersionID"`
	KpiName					string			`json:"kpiName"`
	TargetTeamID			string			`json:"targetTeamID"`
	TargetTeamName			string			`json:"targetTeamName"`
	EffectiveMonth			int				`json:"effectiveMonth"`
	EffectiveYear			int				`json:"effectiveYear"`
	IsPriority				bool			`json:"isPriority"`
	TargetProduct			[]string		`json:"targetProduct"`
	TargetBrand				[]string		`json:"targetBrand"`
	ProductDesign			[]KPIDesignRes	`json:"productDesign"`
	BrandDesign				[]KPIDesignRes	`json:"brandDesign"`
}

type KPIDesignRes struct {
	Name               string           `json:"name"`
	Active             bool             `json:"active"`
	CategoryID         int              `json:"categoryId"`
	Category           string           `json:"category"`
	Type               string           `json:"type"`
	EffectiveStartDate string           `json:"effectiveStartDate"`
	EffectiveEndDate   string           `json:"effectiveEndDate"`
	Questions          []KPIQuestionRes `json:"questions"`
}

type KPIQuestionRes struct {
	QuestionNumber int      `json:"questionNumber"`
	Title          string   `json:"title"`
	Type           string   `json:"type"`
	OptionValues   []string `json:"optionValues"`
	Active         bool     `json:"active"`
	Required       bool     `json:"required"`
}
