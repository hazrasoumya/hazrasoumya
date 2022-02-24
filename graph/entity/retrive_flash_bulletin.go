package entity

import (
	"database/sql"
)

type RetriveFlashBulletin struct {
	ID           sql.NullString `json:"id"`
	Type         sql.NullString `json:"type"`
	Title        sql.NullString `json:"title"`
	Description  sql.NullString `json:"description"`
	ValidityDate sql.NullString `json:"validity_date"`
	Attachments  sql.NullString `json:"attachments"`
	Recipients   sql.NullString `json:"recipients"`
}
