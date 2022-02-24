package entity

import (
	"database/sql"
	"strconv"

	"github.com/eztrade/kpi/graph/model"
	uuid "github.com/gofrs/uuid"
)

type Kpi struct {
	ID             *uuid.UUID
	Name           string
	TypeName       string
	TargetTeam     string
	TargetItems    []string
	EffectiveMonth int
	EffectiveYear  *int
	ExistingID     string
	IsPriority     bool
	Type           int8
	Design         []KpiDesign
	IsActive       bool
}

type KpiDesign struct {
	Name               string        `json:"name"`
	Active             bool          `json:"active"`
	Category           int64         `json:"category"`
	Type               string        `json:"type"`
	EffectiveStartDate string        `json:"effectiveStartDate"`
	EffectiveEndDate   string        `json:"effectiveEndDate"`
	Questions          []KpiQuestion `json:"questions"`
}

type UpsertKpi struct {
	ID             *uuid.UUID
	ProductId      string
	BrandId        string
	Name           string
	TypeName       string
	TargetTeam     string
	TargetProduct  []string
	TargetBrand    []string
	EffectiveMonth int
	EffectiveYear  *int
	ExistingID     string
	IsPriority     bool
	BrandDesign    []UpsertKpiDesign
	ProductDesign  []UpsertKpiDesign
	IsDeleted      bool
}

type UpsertKpiDesign struct {
	Name               string        `json:"name"`
	Category           int64         `json:"category"`
	Type               string        `json:"type"`
	Active             bool          `json:"active"`
	EffectiveStartDate string        `json:"effectiveStartDate"`
	EffectiveEndDate   string        `json:"effectiveEndDate"`
	Questions          []KpiQuestion `json:"questions"`
}

type KpiQuestion struct {
	QuestionNumber int      `json:"questionNumber"`
	Title          string   `json:"title"`
	Type           string   `json:"type"`
	OptionValues   []string `json:"optionValues"`
	Active         bool     `json:"active"`
	Required       bool     `json:"required"`
}

type KpiData struct {
	ParentKpiID         sql.NullString   `json:"parentKpiID"`
	ProductKpiID        sql.NullString   `json:"productKpiID"`
	BrandKpiID          sql.NullString   `json:"brandKpiID"`
	ProductKpiVersionID sql.NullString   `json:"productKpiVersionID"`
	BrandKpiVersionID   sql.NullString   `json:"brandKpiVersionID"`
	KpiName             sql.NullString   `json:"kpiName"`
	TargetTeamID        sql.NullString   `json:"targetTeamID"`
	TargetTeamName      sql.NullString   `json:"targetTeamName"`
	EffectiveMonth      sql.NullInt64    `json:"effectiveMonth"`
	EffectiveYear       sql.NullInt64    `json:"effectiveYear"`
	IsPriority          sql.NullBool     `json:"isPriority"`
	TargetProduct       []sql.NullString `json:"targetProduct"`
	TargetBrand         []sql.NullString `json:"targetBrand"`
	ProductDesign       sql.NullString   `json:"productDesign"`
	BrandDesign         sql.NullString   `json:"brandDesign"`
}

type KpiBrandProductData struct {
	BrandId             sql.NullString `json:"brand_id"`
	BrandName           sql.NullString `json:"brand_name"`
	TeamProductId       sql.NullString `json:"team_product_id"`
	ProductId           sql.NullString `json:"item_id"`
	PrincipalName       sql.NullString `json:"principal_name"`
	IsPriority          sql.NullBool   `json:"is_priority"`
	MaterialDescription sql.NullString `json:"material_description"`
	Type                sql.NullString `json:"type"`
}

type KpiBrandProductInput struct {
	TargetTeam    *string
	BrandID       *string
	TeamProductId *string
	ProductId     *string
	IsPriority    *bool
	SearchIteam   *string
	IsActive      *bool
	IsKpi         *bool
}

type KpiBrandItemData struct {
	BrandId       sql.NullString `json:"item_id"`
	BrandName     sql.NullString `json:"item_name"`
	PrincipalName sql.NullString `json:"principal_name"`
	IsPriority    sql.NullBool   `json:"is_priority"`
}

type KpiProductItemData struct {
	ProductId          sql.NullString `json:"item_id"`
	ProductName        sql.NullString `json:"item_name"`
	ProductDescription sql.NullString `json:"item_description"`
	IsPriority         sql.NullBool   `json:"is_priority"`
	Brand              sql.NullString `json:"brand"`
}

type UniqueBrand struct {
	BrandID   string
	BrandName string
}

type UniqueTeam struct {
	TeamID string
}

type UniqueCustomer struct {
	TeamID     string
	CustomerID string
}

type GetKpisInput struct {
	SalesOrgId    *string
	ParentKpiID   *string
	KpiID         *string
	KpiVersionID  *string
	TeamID        *string
	BrandID       *string
	TeamProductID *string
	Month         *int
	Year          *int
	SearchItem    *string
	Limit         *int
	Offset        *int
	PageNo        *int
}

type GetKpisResponse struct {
	ID             sql.NullString   `json:"id"`
	KpiVersionID   sql.NullString   `json:"kpi_version_id"`
	Name           sql.NullString   `json:"name"`
	TeamName       sql.NullString   `json:"teamName"`
	TeamID         sql.NullString   `json:"team_id"`
	TargetItems    []sql.NullString `json:"target_items"`
	Type           sql.NullString   `json:"type"`
	IsActive       sql.NullBool     `json:"isActive"`
	IsPriority     sql.NullBool     `json:"isPriority"`
	Design         sql.NullString   `json:"design"`
	EffectiveMonth sql.NullInt32    `json:"month"`
	EffectiveYear  sql.NullInt32    `json:"year"`
}

type KpiEvent struct {
	ParentKpiID         sql.NullString   `json:"parentKpiId"`
	ProductKpiId        sql.NullString   `json:"productKpiId"`
	BrandKpiId          sql.NullString   `json:"BrandKpiId"`
	ProductKpiVersionId sql.NullString   `json:"ProductkpiVersionId"`
	BrandKpiVersionId   sql.NullString   `json:"BrandkpiVersionId"`
	KpiName             sql.NullString   `json:"kpiName"`
	TeamID              sql.NullString   `json:"teamMemberID"`
	TeamName            sql.NullString   `json:"teamName"`
	EffectiveMonth      sql.NullInt32    `json:"effectiveMonth"`
	EffectiveYear       sql.NullInt32    `json:"effectiveYear"`
	IsPriority          sql.NullBool     `json:"isPriority"`
	TargetkpiProduct    []sql.NullString `json:"targetKpiProduct"`
	TargetkpiBrand      []sql.NullString `json:"targetKpiBrand"`
	ProductKpiDesign    sql.NullString   `json:"productKpiDesign"`
	BrandKpiDesign      sql.NullString   `json:"BrandKpiDesign"`
}

type UniqueKpiVersion struct {
	ParentKpiID         string
	ProductKpiId        string
	BrandKpiId          string
	ProductKpiVersionId string
	BrandKpiVersionId   string
	KpiName             string
	TeamID              string
	TeamName            string
	EffectiveMonth      int
	EffectiveYear       int
	IsPriority          bool
	TargetkpiProduct    []sql.NullString
	TargetkpiBrand      []sql.NullString
	ProductKpiDesign    string
	BrandKpiDesign      string
}

type KPIDesign struct {
	Name               *string        `json:"name"`
	Active             bool           `json:"active"`
	Category           int64          `json:"category"`
	Type               *string        `json:"type"`
	EffectiveStartDate string         `json:"effectiveStartDate"`
	EffectiveEndDate   string         `json:"effectiveEndDate"`
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

type KpiProductBrandAnswer struct {
	ScheduleEventId       sql.NullString `json:"scheduleEventId"`
	TeamMemberCustomer    sql.NullString `json:"teamMemberCustomer"`
	ItemId                sql.NullString `json:"itemId"`
	ItemName              sql.NullString `json:"itemName"`
	ItemDescription       sql.NullString `json:"itemDescription"`
	ProductIsPriority     sql.NullBool   `json:"productIsPriority"`
	ProductPrincipalName  sql.NullString `json:"productPrincipalName"`
	BrandName             sql.NullString `json:"brandName"`
	ProductId             sql.NullString `json:"productId"`
	KpiId                 sql.NullString `json:"kpiId"`
	KpiType               sql.NullString `json:"kpiType"`
	KpiVersionId          sql.NullString `json:"kpiVersionId"`
	KpiAnsId              sql.NullString `json:"kpiAnsId"`
	Answers               sql.NullString `json:"answers"`
	Category              sql.NullInt32  `json:"category"`
	TargetItem            sql.NullString `json:"targetItem"`
	AnsKpiVersionId       sql.NullString `json:"ansKpiVersionId"`
	AnsScheduleEventId    sql.NullString `json:"ansScheduleEventId"`
	AnsTeamMemberCustomer sql.NullString `json:"ansTeamMemberCustomer"`
	TeamProductID         sql.NullString `json:"teamProductID"`
	TeamID                sql.NullString `json:"teamID"`
	Month                 sql.NullString `json:"month"`
	Year                  sql.NullString `json:"year"`
}

type UniqueEvent struct {
	ScheduleEventId    string
	TeamMemberCustomer string
}

type UniqueBrandProduct struct {
	BrandId              string
	ProductId            string
	PrincipleName        string
	ProductIsPriority    bool
	ProductPrincipalName string
	MaterialDescription  string
	TeamProductID        string
	ProductKpiVersionId  string
	TeamId               string
	Category             int
	Month                string
	Year                 string
}

type UniqueBrandData struct {
	ItemId       string
	ItemName     string
	KpiVersionId string
}

type KPIAnswerStruct struct {
	QuestioNnumber int      `json:"questioNnumber"`
	Value          []string `json:"value"`
}

func (c *Kpi) ValidateData(row int, result *model.ValidationResult) {
	validationMessages := []*model.ValidationMessage{}
	if c.Name == "" {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ":  Name is blank!"}
		validationMessages = append(validationMessages, errorMessage)
	}
	if c.TypeName == "" {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ": Type name is blank!"}
		validationMessages = append(validationMessages, errorMessage)
	}
	if c.EffectiveMonth > 12 || c.EffectiveMonth < 1 {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ": Invalid Effective month!"}
		validationMessages = append(validationMessages, errorMessage)
	}
	if c.EffectiveYear == nil {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ":  EffectiveYear is blank!"}
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
