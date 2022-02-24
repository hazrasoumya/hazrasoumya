package entity

import (
	"database/sql"
)

type PictureData struct {
	Type     sql.NullString `json:"type"`
	Name     sql.NullString `json:"name"`
	Url      sql.NullString `json:"blob_url"`
	Team     sql.NullString `json:"team"`
	Customer sql.NullString `json:"customer"`
}

type CustomerList struct {
	CustomerId sql.NullString `json:"CustomerId"`
	Name       sql.NullString `json:"name"`
	Url        sql.NullString `json:"blob_url"`
}

type ProductList struct {
	ProductId sql.NullString `json:"ProductId"`
	Name      sql.NullString `json:"name"`
	Url       sql.NullString `json:"blob_url"`
}

type BrandList struct {
	BrandId sql.NullString `json:"BrandId"`
	Name    sql.NullString `json:"name"`
	Url     sql.NullString `json:"blob_url"`
}

type TeamIDS struct {
	TeamID sql.NullString `json:teamID`
}
