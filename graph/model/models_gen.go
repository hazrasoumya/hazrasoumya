// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Attachment struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

type AttachmentUpsertInput struct {
	ID       *string `json:"id"`
	Filename string  `json:"filename"`
	URL      string  `json:"url"`
}

type BrandList struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BrandListResponse struct {
	Error   bool         `json:"error"`
	Message string       `json:"message"`
	Data    []*BrandList `json:"data"`
}

type CustomerData struct {
	CustomerID           *string `json:"customerID"`
	TeamMemberCustomerID *string `json:"teamMemberCustomerId"`
	CustomerName         *string `json:"customerName"`
	SoldTo               *int    `json:"soldTo"`
	ShipTo               *int    `json:"shipTo"`
}

type CustomerDataDropDown struct {
	CustomerID           *string `json:"customerID"`
	TeamMemberCustomerID *string `json:"teamMemberCustomerId"`
	CustomerName         *string `json:"customerName"`
	SoldTo               *int    `json:"soldTo"`
	ShipTo               *int    `json:"shipTo"`
}

type CustomerGroup struct {
	InDusTrialCode string              `json:"inDusTrialCode"`
	CustomeDetails []*CustomerResponse `json:"customeDetails"`
}

type CustomerGroupInput struct {
	TeamID        *string   `json:"teamId"`
	CustomerGroup []*string `json:"customerGroup"`
}

type CustomerGroupResponse struct {
	Error        bool             `json:"error"`
	Message      string           `json:"message"`
	CustoMerData []*CustomerGroup `json:"custoMerData"`
}

type CustomerList struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CustomerListResponse struct {
	Error   bool            `json:"error"`
	Message string          `json:"message"`
	Data    []*CustomerList `json:"data"`
}

type CustomerTargetInput struct {
	ID             *string         `json:"id"`
	Type           string          `json:"type"`
	Category       string          `json:"category"`
	ProductBrandID string          `json:"productBrandId"`
	Year           int             `json:"year"`
	Answers        []*TargetValues `json:"answers"`
	IsDeleted      *bool           `json:"isDeleted"`
}

type CustomerTaskFeedBackInput struct {
	TaskBulletinID       string         `json:"taskBulletinId"`
	TeamMemberCustomerID string         `json:"teamMemberCustomerId"`
	Status               string         `json:"status"`
	Remarks              *string        `json:"remarks"`
	Attachments          []*Attachments `json:"attachments"`
}

type CustomerTaskFeedBackResponse struct {
	Error            bool                 `json:"error"`
	Message          *string              `json:"message"`
	ValidationErrors []*ValidationMessage `json:"validationErrors"`
}

type GetBrandProductRequest struct {
	TargetTeam    *string `json:"targetTeam"`
	BrandID       *string `json:"brandId"`
	TeamProductID *string `json:"teamProductId"`
	ProductID     *string `json:"productId"`
	IsPriority    *bool   `json:"isPriority"`
	SearchItem    *string `json:"searchItem"`
	IsActive      bool    `json:"isActive"`
	IsKpi         *bool   `json:"isKpi"`
}

type GetKpiBrandProductResponse struct {
	Error     bool            `json:"error"`
	Message   string          `json:"message"`
	ErrorCode int             `json:"errorCode"`
	Brands    []*KpiBrandItem `json:"brands"`
}

type GetKpiOffline struct {
	ParentKpiID         string          `json:"parentKpiId"`
	ProductKpiID        string          `json:"productKpiId"`
	BrandKpiID          string          `json:"brandKpiId"`
	ProductKpiVersionID string          `json:"productKpiVersionId"`
	BrandKpiVersionID   string          `json:"brandKpiVersionId"`
	KpiName             string          `json:"kpiName"`
	TargetTeamID        string          `json:"targetTeamId"`
	TargetTeamName      string          `json:"targetTeamName"`
	EffectiveMonth      int             `json:"effectiveMonth"`
	EffectiveYear       int             `json:"effectiveYear"`
	IsPriority          bool            `json:"isPriority"`
	TargetProduct       []string        `json:"targetProduct"`
	TargetBrand         []string        `json:"targetBrand"`
	ProductDesign       []*KPIDesignRes `json:"productDesign"`
	BrandDesign         []*KPIDesignRes `json:"brandDesign"`
}

type GetKpiResponse struct {
	Error     bool      `json:"error"`
	Message   *string   `json:"message"`
	ErrorCode int       `json:"errorCode"`
	TotalPage int       `json:"totalPage"`
	Data      []*GetKpi `json:"data"`
}

type GetTargetCustomerRequest struct {
	IsExcel          bool    `json:"isExcel"`
	ID               *string `json:"id"`
	Type             *string `json:"type"`
	Category         *string `json:"category"`
	ProductBrandID   *string `json:"productBrandId"`
	ProductBrandName *string `json:"productBrandName"`
	Year             *int    `json:"year"`
	PageNo           *int    `json:"pageNo"`
	Limit            *int    `json:"limit"`
}

type GetTargetCustomerResponse struct {
	URL       string            `json:"url"`
	Error     bool              `json:"error"`
	Message   *string           `json:"message"`
	TotalPage int               `json:"totalPage"`
	Data      []*TargetCustomer `json:"data"`
}

type KpiBrandItem struct {
	BrandID   string            `json:"brandId"`
	BrandName string            `json:"brandName"`
	Products  []*KpiProductItem `json:"products"`
}

type KpiBrandItemOffline struct {
	BrandID           string                   `json:"brandId"`
	BrandName         string                   `json:"brandName"`
	BrandKpiVersionID *string                  `json:"brandKpiVersionId"`
	BrandKpiAnswer    []*KpiAnswer             `json:"brandKpiAnswer"`
	Products          []*KpiProductItemOffline `json:"products"`
}

type KpiOfflineInput struct {
	StartDate *int `json:"startDate"`
	EndDate   *int `json:"endDate"`
}

type KpiOfflineResponse struct {
	Error                 bool                     `json:"error"`
	ErrorCode             int                      `json:"errorCode"`
	Message               string                   `json:"message"`
	GetKpiOffline         []*GetKpiOffline         `json:"getKpiOffline"`
	KpiProductBrandAnswer []*KpiProductBrandAnswer `json:"kpiProductBrandAnswer"`
}

type KpiProductBrandAnswer struct {
	EventID        string                 `json:"eventID"`
	TeamCustomerID string                 `json:"teamCustomerID"`
	Brands         []*KpiBrandItemOffline `json:"brands"`
}

type KpiProductItem struct {
	TeamProductID       string `json:"teamProductId"`
	ProductID           string `json:"productId"`
	PrincipalName       string `json:"principalName"`
	MaterialDescription string `json:"materialDescription"`
	IsPriority          bool   `json:"isPriority"`
}

type KpiProductItemOffline struct {
	TeamID              string       `json:"teamId"`
	TeamProductID       string       `json:"teamProductId"`
	ProductID           string       `json:"productId"`
	PrincipalName       string       `json:"principalName"`
	MaterialDescription string       `json:"materialDescription"`
	IsPriority          bool         `json:"isPriority"`
	ProductKpiVersionID *string      `json:"productKpiVersionId"`
	ProductKpiAnswer    []*KpiAnswer `json:"productKpiAnswer"`
}

type KpiTaregetTitleResponse struct {
	Error           bool              `json:"error"`
	Message         string            `json:"message"`
	KpiTargetTitles []*KpiTargetTitle `json:"kpiTargetTitles"`
}

type KpiTargetTitle struct {
	Title       string `json:"title"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type ListInput struct {
	TeamID *string `json:"teamID"`
}

type PictureZipInput struct {
	Selections []string `json:"selections"`
	Type       string   `json:"type"`
}

type PrincipalDropDownInput struct {
	TeamID string `json:"teamID"`
}

type PrincipalDropDownResponse struct {
	Error   bool                     `json:"error"`
	Message *string                  `json:"message"`
	Data    []*PrincipalDropDownData `json:"data"`
}

type ProductList struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductListResponse struct {
	Error   bool           `json:"error"`
	Message string         `json:"message"`
	Data    []*ProductList `json:"data"`
}

type Recipients struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type RetrievePictureZip struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

type SalesRepData struct {
	TeamMemberID       *string         `json:"teamMemberId"`
	UserID             *string         `json:"userId"`
	FirstName          *string         `json:"firstName"`
	LastName           *string         `json:"lastName"`
	ActiveDirectory    *string         `json:"activeDirectory"`
	Email              *string         `json:"email"`
	ApprovalRoleTitle  *string         `json:"approvalRoleTitle"`
	ApprovalRoleValues *string         `json:"approvalRoleValues"`
	Customers          []*CustomerData `json:"customers"`
}

type SalesRepDataDropDown struct {
	TeamMemberID       *string                 `json:"teamMemberId"`
	UserID             *string                 `json:"userId"`
	FirstName          *string                 `json:"firstName"`
	LastName           *string                 `json:"lastName"`
	ActiveDirectory    *string                 `json:"activeDirectory"`
	Email              *string                 `json:"email"`
	ApprovalRoleTitle  *string                 `json:"approvalRoleTitle"`
	ApprovalRoleValues *string                 `json:"approvalRoleValues"`
	Customers          []*CustomerDataDropDown `json:"customers"`
}

type Target struct {
	Month *int     `json:"month"`
	Value *float64 `json:"value"`
}

type TargetCustomer struct {
	CustomerTargetID *string   `json:"customerTargetId"`
	Type             *string   `json:"type"`
	Category         *string   `json:"category"`
	ProductBrandID   *string   `json:"productBrandId"`
	ProductBrandName *string   `json:"productBrandName"`
	Year             *int      `json:"year"`
	Targets          []*Target `json:"targets"`
}

type TaskBulletinInput struct {
	TeamID               []*string `json:"teamID"`
	TeamMemberID         *string   `json:"teamMemberId"`
	TeamMemberCustomerID *string   `json:"teamMemberCustomerId"`
	CustomerID           *string   `json:"customerID"`
	OnlySalesrep         *bool     `json:"onlySalesrep"`
}

type TaskBulletinReportInput struct {
	IsExcel   bool     `json:"isExcel"`
	Tittle    []string `json:"tittle"`
	DateRange []string `json:"dateRange"`
}

type TaskBulletinResponse struct {
	Error    bool                  `json:"error"`
	Message  string                `json:"message"`
	DropDown []*TeamMemberDropdown `json:"dropDown"`
}

type TaskBulletinTitleInput struct {
	TeamID *string `json:"teamID"`
}

type TaskBulletinTitleResponse struct {
	Error       bool          `json:"error"`
	Message     *string       `json:"message"`
	TypeDetails []*TitleValue `json:"typeDetails"`
}

type TaskBulletinUpsertInput struct {
	ID                 *string        `json:"id"`
	CreationDate       string         `json:"creationDate"`
	TargetDate         string         `json:"targetDate"`
	TeamMemberCustomer []*string      `json:"teamMemberCustomer"`
	Type               string         `json:"type"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	PrincipalName      string         `json:"principalName"`
	IsDeleted          *bool          `json:"isDeleted"`
	Attachments        []*Attachments `json:"attachments"`
}

type TaskBulletinUpsertResponse struct {
	Error            bool                 `json:"error"`
	Message          string               `json:"message"`
	ValidationErrors []*ValidationMessage `json:"validationErrors"`
}

type TaskReportOutput struct {
	Error   bool          `json:"error"`
	Message string        `json:"message"`
	URL     string        `json:"url"`
	Values  []*TaskReport `json:"values"`
}

type TeamMemberDropdown struct {
	TeamID   *string                 `json:"teamId"`
	TeamName *string                 `json:"teamName"`
	Employee []*SalesRepDataDropDown `json:"employee"`
}

type WeeklyFeedback struct {
	WeekNumber    *int          `json:"weekNumber"`
	WeekDateValue *string       `json:"weekDateValue"`
	Status        *string       `json:"status"`
	Remarks       *string       `json:"remarks"`
	Attachments   []*Attachment `json:"attachments"`
}

type Attachments struct {
	ID       *string `json:"id"`
	URL      string  `json:"url"`
	Filename string  `json:"filename"`
}

type CustomerFeedback struct {
	StatusTitle string        `json:"statusTitle"`
	StatusValue string        `json:"statusValue"`
	Remarks     string        `json:"remarks"`
	DateCreated string        `json:"dateCreated"`
	Attachments []*Attachment `json:"attachments"`
}

type CustomerResponse struct {
	CustoMerID   string `json:"custoMerId"`
	CustoMerName string `json:"custoMerName"`
	SoldTo       string `json:"soldTo"`
	ShipTo       string `json:"shipTo"`
}

type FetchCustomerFeedbackInput struct {
	TeamMemberCustomerID string `json:"teamMemberCustomerId"`
	TaskBulletinID       string `json:"taskBulletinId"`
}

type FetchCustomerFeedbackResponse struct {
	Error            bool                `json:"error"`
	Message          string              `json:"message"`
	CustomerFeedback []*CustomerFeedback `json:"customerFeedback"`
}

type FlashBulletin struct {
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	ValidityDate string        `json:"validityDate"`
	Attachments  []*Attachment `json:"attachments"`
	Recipients   []*Recipients `json:"recipients"`
}

type FlashBulletinData struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Status       bool     `json:"status"`
	StartDate    string   `json:"startDate"`
	EndDate      string   `json:"endDate"`
	Type         string   `json:"type"`
	CreatedDate  string   `json:"createdDate"`
	ModifiedDate string   `json:"modifiedDate"`
	Attachments  []string `json:"attachments"`
}

type FlashBulletinResponse struct {
	Error            bool                 `json:"error"`
	Message          string               `json:"message"`
	ValidationErrors []*ValidationMessage `json:"validationErrors"`
}

type FlashBulletinUpsertInput struct {
	ActiveDirName     string                   `json:"activeDirName"`
	ID                *string                  `json:"id"`
	Type              int                      `json:"type"`
	Title             string                   `json:"title"`
	Description       string                   `json:"description"`
	ValidityDateStart string                   `json:"validity_date_start"`
	ValidityDateEnd   string                   `json:"validity_date_end"`
	Attachments       []*AttachmentUpsertInput `json:"attachments"`
	Recipients        []string                 `json:"recipients"`
	IsDeleted         *bool                    `json:"isDeleted"`
	IsActive          *bool                    `json:"isActive"`
	CustomerGroup     []*string                `json:"customerGroup"`
}

type FlashBulletinUpsertResponse struct {
	Error            bool                 `json:"error"`
	Message          string               `json:"message"`
	ValidationErrors []*ValidationMessage `json:"validationErrors"`
}

type GetKpiInput struct {
	ParentKpiID   *string `json:"parentKpiId"`
	KpiID         *string `json:"kpiId"`
	KpiVersionID  *string `json:"kpiVersionId"`
	TeamID        *string `json:"teamId"`
	BrandID       *string `json:"brandId"`
	TeamProductID *string `json:"teamProductId"`
	Month         *int    `json:"month"`
	Year          *int    `json:"year"`
	SearchItem    *string `json:"searchItem"`
	Limit         *int    `json:"limit"`
	PageNo        *int    `json:"pageNo"`
}

type KpiResponse struct {
	Error            bool                 `json:"error"`
	Message          string               `json:"message"`
	ErrorCode        int                  `json:"errorCode"`
	ValidationErrors []*ValidationMessage `json:"validationErrors"`
}

type ListFlashBulletinInput struct {
	Type                 *int    `json:"type"`
	IsActive             *bool   `json:"isActive"`
	StartDate            *string `json:"startDate"`
	EndDate              *string `json:"endDate"`
	ReceipientID         *string `json:"receipientId"`
	TeamMemberCustomerID *string `json:"teamMemberCustomerId"`
}

type ListFlashBulletinResponse struct {
	Error          bool                 `json:"error"`
	Message        string               `json:"message"`
	FlashBulletins []*FlashBulletinData `json:"flashBulletins"`
}

type ListTaskBulletinInput struct {
	ID                   *string `json:"id"`
	Type                 *string `json:"type"`
	IsActive             *bool   `json:"isActive"`
	CreationDate         *string `json:"creationDate"`
	TargetDate           *string `json:"targetDate"`
	TeamMemberID         *string `json:"teamMemberId"`
	TeamMemberCustomerID *string `json:"teamMemberCustomerId"`
	PageNo               *int    `json:"pageNo"`
	Limit                *int    `json:"limit"`
	SearchItem           *string `json:"searchItem"`
}

type ListTaskBulletinResponse struct {
	Error         bool                `json:"error"`
	Message       string              `json:"message"`
	TotalPages    int                 `json:"totalPages"`
	TaskBulletins []*TaskBulletinData `json:"taskBulletins"`
}

type PrincipalDropDownData struct {
	PrincipalName string `json:"principalName"`
}

type RetriveInfoFlashBulletinInput struct {
	BulletinID string `json:"bulletinID"`
}

type RetriveInfoFlashBulletinleResponse struct {
	Error             bool           `json:"error"`
	Message           string         `json:"message"`
	FlashBulletinData *FlashBulletin `json:"flashBulletinData"`
}

type TaskBulletinData struct {
	ID            string          `json:"id"`
	CreationDate  string          `json:"creationDate"`
	TargetDate    string          `json:"targetDate"`
	TypeTitle     string          `json:"typeTitle"`
	TypeValue     string          `json:"typeValue"`
	PrincipalName string          `json:"principalName"`
	TeamID        *string         `json:"teamId"`
	TeamName      string          `json:"teamName"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	SalesRep      []*SalesRepData `json:"salesRep"`
	Attachments   []*Attachment   `json:"attachments"`
}

type TaskReport struct {
	BulletinTitle   *string           `json:"bulletinTitle"`
	BulletinType    *string           `json:"bulletinType"`
	PrincipalName   *string           `json:"principalName"`
	CustomerName    *string           `json:"customerName"`
	TeamName        *string           `json:"teamName"`
	UserName        *string           `json:"userName"`
	ActiveDirectory *string           `json:"activeDirectory"`
	CreationDate    *string           `json:"creationDate"`
	TargetDate      *string           `json:"targetDate"`
	WeeklyFeedback  []*WeeklyFeedback `json:"weeklyFeedback"`
}

type TitleValue struct {
	Type string `json:"type"`
}

type ValidationMessage struct {
	Row       int    `json:"row"`
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

type ValidationResult struct {
	Error               bool                 `json:"error"`
	ValidationTimeTaken string               `json:"validationTimeTaken"`
	ValidationMessage   []*ValidationMessage `json:"validationMessage"`
}
