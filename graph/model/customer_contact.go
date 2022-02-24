package model

type CustomerContact struct {
	ID            string `json:"id"`
	ContactName   string `json:"contactName"`
	Designation   string `json:"designation"`
	ContactNumber string `json:"contactNumber"`
	ContactImage  string `json:"contactImage"`
	CustomerID    string `json:"customerId"`
	CustomerName  string `json:"customerName"`
	HasConsent    bool   `json:"hasConsent"`
	EmailID       string `json:"emailId"`
}

type CustomerContactRequest struct {
	ID                  *string `json:"id"`
	ContactName         string  `json:"contactName"`
	Designation         string  `json:"designation"`
	ContactNumber       string  `json:"contactNumber"`
	ContactImage        *string `json:"contactImage"`
	AuthorActiveDirName string  `json:"authorActiveDirName"`
	CustomerID          string  `json:"customerId"`
	EmailID             *string `json:"emailId"`
}

type CustomerContactDeleteRequest struct {
	ID                  string `json:"id"`
	AuthorActiveDirName string `json:"authorActiveDirName"`
}

type GetCustomerContactRequest struct {
	ID                     *string `json:"id"`
	CustomerID             *string `json:"customerId"`
	TeamMememberCustomerID *string `json:"teamMememberCustomerID"`
}

type GetCustomerContactResponse struct {
	GetCustomerContact []*CustomerContact `json:"getCustomerContact"`
}
