package postgres

import (
	"context"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/kpi/graph/util"
	suuid "github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
)

func UpsertCustomerContactData(entity *entity.CustomerContact, response *model.KpiResponse, loggedInUserEntity *entity.LoggedInUser) {
	if pool == nil {
		pool = GetPool()
	}
	tx, err := pool.Begin(context.Background())
	if err != nil {
		response.Message = "Failed to begin transaction"
		response.Error = true
	}
	defer tx.Rollback(context.Background())
	timenow := util.GetCurrentTime()

	if entity.ID == nil {
		var dbContactID suuid.UUID
		queryContactString := `INSERT INTO customer_contacts (contact_name, designation, contact_number, created_by, date_created, is_active, is_deleted, customer, contact_image, email) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)RETURNING(id)`
		err := tx.QueryRow(context.Background(), queryContactString, entity.ContactName, entity.Designation, entity.ContactNumber, loggedInUserEntity.ID, timenow, true, false, entity.CustomerID, entity.ContactImage, entity.EmailID).Scan(&dbContactID)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			response.Message = "Failed to insert in customer contact"
			response.Error = true
		} else {
			response.Message = "Customer Contact successfully inserted"
			response.Error = false
		}
	} else {
		querystring := "UPDATE customer_contacts SET contact_name=$2, designation=$3, contact_number=$4, modified_by=$5, last_modified=$6, customer=$7, contact_image=$8, email=$9 WHERE id=$1"
		commandTag, err := tx.Exec(context.Background(), querystring, entity.ID, entity.ContactName, entity.Designation, entity.ContactNumber, loggedInUserEntity.ID, timenow, entity.CustomerID, entity.ContactImage, entity.EmailID)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			response.Message = "Failed to update customer contact data"
			response.Error = true
		} else if commandTag.RowsAffected() != 1 {
			response.Message = "Failed to update customer contact data"
			response.Error = true
		} else {
			response.Message = "Customer Contact successfully Updated"
			response.Error = false
		}
	}
	txErr := tx.Commit(context.Background())
	if txErr != nil {
		response.Message = "Failed to commit customer contact data"
		response.Error = true
	}
}

func HasConsentFormSalesOrg(formType string, salesOrgId string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	var hasResult int
	queryString := `select 1 from form_version fv
	inner join form f on f.id = fv.form_id 
	inner join code c on c.id = f."type" 
	where f.is_active = true and f.is_deleted = false 
	and fv.is_active = true and fv.is_deleted = false
	and c.is_active = true and c.is_delete = false 
	and c.category = 'FormType' and c.value = $1
	and fv.sales_organisation = $2`
	err := pool.QueryRow(context.Background(), queryString, formType, salesOrgId).Scan(&hasResult)

	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}

func GetCustomerContactData(inputModel *model.GetCustomerContactRequest, loggedInUserEntity *entity.LoggedInUser) ([]entity.CustomerContactData, error) {
	if pool == nil {
		pool = GetPool()
	}
	query := `SELECT distinct(cc.id), cc.contact_name, cc.designation, cc.contact_number, cc.contact_image, cc.customer, c."name",cc.email,
	case 
		when cc.has_consent = false then false
		when
			(SELECT row_number FROM
				(
				SELECT fs.version_id,fv.id, fv.date_created, ROW_NUMBER () OVER (ORDER BY fv.date_created desc)
				FROM customer_contacts
				INNER JOIN (SELECT id, version_id FROM form_submission) fs
				ON fs.id = customer_contacts.consent_submission
				INNER JOIN (SELECT id, form_id, date_created, is_active FROM form_version) fv
				ON fv.form_id = (SELECT form_id FROM form_version WHERE id = fs.version_id)
				WHERE customer_contacts.id = cc.id AND ((fv.is_active = true AND fv.id <> fs.version_id ) OR fv.id = fs.version_id )
				) x
				WHERE version_id = id
			) = 1 then true ELSE false END
		as has_consent
	FROM customer_contacts cc 
		INNER JOIN customer c ON c.id = cc.customer	
		INNER JOIN team_member_customer tmc ON tmc.customer = c.id	
		INNER JOIN "user" u on u.id = cc.created_by 
	WHERE u.sales_organisation = ?
		AND cc.is_active = true 
		AND cc.is_deleted = false 
		AND c.is_active = true 
		AND c.is_deleted = false
		AND tmc.is_active = true 
		AND tmc.is_deleted = false`

	var rows pgx.Rows
	var err error
	var inputArgs []interface{}

	inputArgs = append(inputArgs, loggedInUserEntity.SalesOrganisaton)

	if inputModel != nil {
		if inputModel.ID != nil && *inputModel.ID != "" {
			query = query + ` and cc.id = ?`
			inputArgs = append(inputArgs, *inputModel.ID)
		}
		if inputModel.CustomerID != nil && *inputModel.CustomerID != "" {
			query = query + ` and cc.customer = ?`
			inputArgs = append(inputArgs, *inputModel.CustomerID)
		}
		if inputModel.TeamMememberCustomerID != nil && *inputModel.TeamMememberCustomerID != "" {
			query = query + ` and tmc.id = ?`
			inputArgs = append(inputArgs, *inputModel.TeamMememberCustomerID)
		}
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err = pool.Query(context.Background(), query, inputArgs...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.CustomerContactData{}, err
	}

	customerContacts := []entity.CustomerContactData{}
	defer rows.Close()
	for rows.Next() {
		contacts := entity.CustomerContactData{}
		err := rows.Scan(
			&contacts.ID,
			&contacts.ContactName,
			&contacts.Designation,
			&contacts.ContactNumber,
			&contacts.ContactImage,
			&contacts.CustomerID,
			&contacts.CustomerName,
			&contacts.EmailID,
			&contacts.HasConsent,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return customerContacts, err
		}
		customerContacts = append(customerContacts, contacts)
	}

	return customerContacts, err
}

func DeleteCustomerContactData(entity *entity.CustomerContactDelete, response *model.KpiResponse, loggedInUserEntity *entity.LoggedInUser) {
	if pool == nil {
		pool = GetPool()
	}
	tx, err := pool.Begin(context.Background())
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		response.Message = "Failed to begin transaction"
		response.Error = true
	}
	defer tx.Rollback(context.Background())
	timenow := util.GetCurrentTime()

	querystring := "UPDATE customer_contacts SET is_active=$2, is_deleted=$3, modified_by=$4, last_modified=$5 WHERE id=$1 AND is_active=true AND is_deleted=false"
	commandTag, err := tx.Exec(context.Background(), querystring, entity.ID, false, true, loggedInUserEntity.ID, timenow)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		response.Message = "Failed to delete customer contact data"
		response.Error = true
	} else if commandTag.RowsAffected() != 1 {
		response.Message = "Failed to delete customer contact data"
		response.Error = true
	} else {
		response.Message = "Customer Contact successfully Deleted"
		response.Error = false
	}
	txErr := tx.Commit(context.Background())
	if txErr != nil {
		response.Message = "Failed to commit customer contact data"
		response.Error = true
	}
}
