package model

type PicturesInput struct {
	TeamID     string    `json:"teamID"`
	Type       string    `json:"type"`
	ProductID  []*string `json:"productID"`
	BrandID    []*string `json:"brandID"`
	CustomerID []*string `json:"customerID"`
	StartDate  string    `json:"startDate"`
	EndDate    string    `json:"endDate"`
}

type RetrievePicturesResponse struct {
	Error   bool    `json:"error"`
	Message string  `json:"message"`
	Data    []Teams `json:"data"`
}

type Teams struct {
	TeamID    string     `json:"teamID"`
	TeamName  string     `json:"teamName"`
	Customers []Customer `json:"customers"`
}

type Customer struct {
	CustomerID    string   `json:"customerID"`
	CustomerName  string   `json:"customerName"`
	Product       Images   `json:"product"`
	Brand         Images   `json:"brand"`
	Competitor    Images   `json:"competitor"`
	Survey        Images   `json:"survey"`
	Promotion     Images   `json:"promotion"`
	FlashBulletin []Images `json:"flashbulletin"`
}

type Images struct {
	Name string   `json:"name"`
	URL  []string `json:"url"`
}
