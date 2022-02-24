package entity

import (
	"strconv"

	"github.com/eztrade/kpi/graph/model"
	uuid "github.com/gofrs/uuid"
)

type FlashBulletin struct {
	ID            *uuid.UUID
	TypeID        int
	TypeValue     string
	Title         string
	Description   string
	ValidityDate  string
	Attachments   []Attachment
	Recipients    []string
	CreatedBy     uuid.UUID
	ModifiedBy    uuid.UUID
	IsDeleted     *bool
	IsActive      *bool
	SalesOrgId    uuid.UUID
	CustomerGroup []string
}

func (c *FlashBulletin) ValidateData(int, *model.ValidationResult) {
	var row int
	var result *model.ValidationResult
	validationMessages := []*model.ValidationMessage{}
	if c.TypeID == 0 {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ":  Type is blank!"}
		validationMessages = append(validationMessages, errorMessage)
	}
	if c.Title == "" {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ": Title is blank!"}
		validationMessages = append(validationMessages, errorMessage)
	}
	if c.Description == "" {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ": Description is blank!"}
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

type Attachment struct {
	ID         *uuid.UUID
	BlobName   string
	BlobUrl    string
	CreatedBy  uuid.UUID
	ModifiedBy uuid.UUID
}

func (c *Attachment) ValidateData(row int, result *model.ValidationResult) {
	validationMessages := []*model.ValidationMessage{}
	if c.BlobName == "" {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ": File name is blank!"}
		validationMessages = append(validationMessages, errorMessage)
	}
	if c.BlobUrl == "" {
		errorMessage := &model.ValidationMessage{Row: row, Message: "Row " + strconv.Itoa(row) + ": Url is blank!"}
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
