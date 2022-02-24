package postgres

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/kpi/graph/util"
	suuid "github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func HasKpiId(kpiID *suuid.UUID) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	queryString := "select 1 from kpi where id = $1 AND is_deleted = false"
	var hasValue int
	err := pool.QueryRow(context.Background(), queryString, kpiID).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}

func IsKPIVersionId(kpiVerId *suuid.UUID) (bool, bool) {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	queryString := `select k.is_active from kpi k
	inner join kpi_versions kv on kv.kpi_id=k.id
	where kv.is_active = true and kv.is_deleted=false and
	k.is_deleted = false
	and kv.id = $1`
	var isActive bool
	err := pool.QueryRow(context.Background(), queryString, kpiVerId).Scan(&isActive)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result, isActive
}

func CheckNameForSameTeamAndMonth(name string, teamID string, month int) (bool, string) {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	queryString := "select id from kpi where name = $1 and target_team= $2 and month = $3 AND is_deleted = false"
	var dbID string
	err := pool.QueryRow(context.Background(), queryString, name, teamID, month).Scan(&dbID)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result, dbID
}

func GetIDForSameTeamAndMonth(teamID string, month int, year int) (string, error) {
	if pool == nil {
		pool = GetPool()
	}
	var dbID string
	queryString := "select id from kpi where target_team = $1 and month = $2 and year = $3 and is_active = true and is_deleted = false"
	err := pool.QueryRow(context.Background(), queryString, teamID, month, year).Scan(&dbID)
	if err == nil {
		return dbID, nil
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return dbID, err
	}
}

func GetIDProductBrandIdFromKpiParentId(parentId, value string) string {
	if pool == nil {
		pool = GetPool()
	}
	var ID string
	queryString := `select id from kpi where parent_kpi = $1 and 
	"type" = (select id from code where category = 'KPIType' and value = $2 and is_active = true and is_delete = false)
	and is_active = true and is_deleted = false`
	err := pool.QueryRow(context.Background(), queryString, parentId, value).Scan(&ID)
	if err == nil {
		return ID
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return ID
	}
}

func KpiCombinationExist(teamID string, month int, year int) (bool, error) {
	if pool == nil {
		pool = GetPool()
	}
	var dbID int
	queryString := "select 1 from kpi where target_team = $1 and month = $2 and year = $3 and is_active = true and is_deleted = false"
	err := pool.QueryRow(context.Background(), queryString, teamID, month, year).Scan(&dbID)
	if err == nil {
		if dbID == 1 {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return false, err
	}
}

func IsMonthyearValid(is_priority bool, month int, year int, parentId string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var dbID int
	queryString := "select 1 from kpi where month = $1 and year = $2 and is_priority = $3 and id = $4 and is_active = true and is_deleted = false"
	err := pool.QueryRow(context.Background(), queryString, month, year, is_priority, parentId).Scan(&dbID)
	if err == nil {
		if dbID == 1 {
			return true
		} else {
			return false
		}
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return false
	}
}

func ProductBelongInAnyBrand(BrandID string, TargetProduct []string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var dbID int
	var IDs []string
	for _, productId := range TargetProduct {
		IDs = append(IDs, "'"+productId+"'")
	}
	stringTeamIds := strings.Trim(strings.Join(IDs, ", "), ", ")
	queryString := `select 1 from team_products tp 
	inner join product p2 on p2.id = tp.material_code
	where p2.brand = $1 and tp.id in (` + stringTeamIds + `)
	and tp.is_active = true and tp.is_delete = false 
    and p2.is_active = true and p2.is_deleted = false`
	err := pool.QueryRow(context.Background(), queryString, BrandID).Scan(&dbID)
	if err == nil {
		if dbID == 1 {
			return true
		} else {
			return false
		}
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return false
	}
}
func IsValidCategory(input string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var value int
	var result bool
	queryString := `select 1 from code where value = $1 and category ='MerchandisingType'`
	err := pool.QueryRow(context.Background(), queryString, input).Scan(&value)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}
func ValidateKpi(entity *entity.Kpi, validationResult *model.ValidationResult) {
	validationMessages := []*model.ValidationMessage{}
	if entity.ID != nil {
		if !HasKpiId(entity.ID) {
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Kpi ID does not exists!"}
			validationMessages = append(validationMessages, errorMessage)
		}
	}
	if strings.ToLower(entity.TypeName) == "product" {
		for key, item := range entity.TargetItems {
			if &item != nil {
				if HasTeamProductId(item) {
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Positon " + strconv.Itoa(key+1) + " product id is not valid"}
					validationMessages = append(validationMessages, errorMessage)
				}
			}
		}
	}
	if strings.ToLower(entity.TypeName) == "brand" {
		for key, item := range entity.TargetItems {
			if &item != nil {
				if HasBrandId(item) {
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Positon " + strconv.Itoa(key+1) + " brand id is not valid"}
					validationMessages = append(validationMessages, errorMessage)
				}
			}
		}
	}
	if entity.TargetTeam != "" && HasTeamId(entity.TargetTeam) {
		errorMessage := &model.ValidationMessage{Row: 0, Message: "team id is not valid"}
		validationMessages = append(validationMessages, errorMessage)
	}
	if len(validationMessages) > 0 {
		if !validationResult.Error {
			validationResult.Error = true
			validationResult.ValidationMessage = []*model.ValidationMessage{}
		}
		validationResult.ValidationMessage = append(validationResult.ValidationMessage, validationMessages...)
	}
}

func IsTargetItemValidForKPIAnswer(kpiVerId *suuid.UUID, targetItem *suuid.UUID) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	querystring := `WITH expand AS (
		SELECT UNNEST(k.target_items) AS item_id 
		FROM kpi k
			INNER JOIN kpi_versions kv ON kv.kpi_id = k.id	
		WHERE 
			k.is_active = TRUE AND k.is_deleted = FALSE AND
			kv.is_active = TRUE AND kv.is_deleted = FALSE AND
			kv.id = $1
	), lookup AS (
		SELECT 1 FROM expand e WHERE e.item_id = $2
		  )
	SELECT * FROM lookup`
	var hasValue int
	err := pool.QueryRow(context.Background(), querystring, kpiVerId, targetItem).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}

func SaveKpiInformation(entity entity.UpsertKpi, loggedInUserEntity *entity.LoggedInUser) model.KpiResponse {
	var response model.KpiResponse
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
	var ParentKpiId string
	var productId string

	if entity.ID == nil {
		querystring := `INSERT INTO kpi (name, type, target_team, target_items, is_active,is_deleted,created_by, date_created, month, year, is_priority) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)RETURNING(id)`
		err = tx.QueryRow(context.Background(), querystring, entity.Name, 27, entity.TargetTeam, entity.TargetBrand, true, false, loggedInUserEntity.ID, timenow, entity.EffectiveMonth, entity.EffectiveYear, entity.IsPriority).Scan(&ParentKpiId)
		if err != nil {
			if !response.Error {
				response.Error = true
				response.ValidationErrors = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to insert in kpi"}
			response.ValidationErrors = append(response.ValidationErrors, errorMessage)
			return response
		}

		updateBrandKPI := `update kpi set parent_kpi = $2 where id = $1`
		commandTag, err := tx.Exec(context.Background(), updateBrandKPI, ParentKpiId, ParentKpiId)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			if !response.Error {
				response.Error = true
				response.ValidationErrors = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
			response.ValidationErrors = append(response.ValidationErrors, errorMessage)
			return response
		} else if commandTag.RowsAffected() != 1 {
			logengine.GetTelemetryClient().TrackException(err.Error())
			if !response.Error {
				response.Error = true
				response.ValidationErrors = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
			response.ValidationErrors = append(response.ValidationErrors, errorMessage)
			return response
		} else {
			response.Message = "Kpi successfully Updated"
			response.Error = false
		}

		querystring = `INSERT INTO kpi (name, type, target_team, target_items, is_active,is_deleted,created_by, date_created, month, year, is_priority, parent_kpi) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)RETURNING(id)`
		err = tx.QueryRow(context.Background(), querystring, entity.Name, 26, entity.TargetTeam, entity.TargetProduct, true, false, loggedInUserEntity.ID, timenow, entity.EffectiveMonth, entity.EffectiveYear, entity.IsPriority, ParentKpiId).Scan(&productId)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			if !response.Error {
				response.Error = true
				response.ValidationErrors = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to insert in kpi"}
			response.ValidationErrors = append(response.ValidationErrors, errorMessage)
			return response
		}

		querystring = `INSERT INTO kpi_versions (design, kpi_id, is_active,is_deleted,created_by, date_created) VALUES($1, $2, $3, $4, $5, $6)`
		commandTag, err = tx.Exec(context.Background(), querystring, entity.BrandDesign, ParentKpiId, true, false, loggedInUserEntity.ID, timenow)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			if !response.Error {
				response.Error = true
				response.ValidationErrors = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
			response.ValidationErrors = append(response.ValidationErrors, errorMessage)
			return response
		} else if commandTag.RowsAffected() != 1 {
			if !response.Error {
				response.Error = true
				response.ValidationErrors = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
			response.ValidationErrors = append(response.ValidationErrors, errorMessage)
			return response
		} else {
			response.Message = "Kpi successfully Updated"
			response.Error = false
		}

		querystring = `INSERT INTO kpi_versions (design, kpi_id, is_active,is_deleted,created_by, date_created) VALUES($1, $2, $3, $4, $5, $6)`
		commandTag, err = tx.Exec(context.Background(), querystring, entity.ProductDesign, productId, true, false, loggedInUserEntity.ID, timenow)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			if !response.Error {
				response.Error = true
				response.ValidationErrors = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
			response.ValidationErrors = append(response.ValidationErrors, errorMessage)
			return response
		} else if commandTag.RowsAffected() != 1 {
			if !response.Error {
				response.Error = true
				response.ValidationErrors = []*model.ValidationMessage{}
			}
			errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
			response.ValidationErrors = append(response.ValidationErrors, errorMessage)
			return response
		} else {
			response.Message = "Kpi successfully Inserted!"
			response.Error = false
		}
	} else {
		if entity.ID != nil && entity.IsDeleted == false {
			if entity.TargetBrand != nil {

				if entity.BrandId == "" {
					querystring := `INSERT INTO kpi (name, type, target_team, target_items, is_active,is_deleted,created_by, date_created, month, year, is_priority) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)RETURNING(id)`
					err = tx.QueryRow(context.Background(), querystring, entity.Name, 27, entity.TargetTeam, entity.TargetBrand, true, false, loggedInUserEntity.ID, timenow, entity.EffectiveMonth, entity.EffectiveYear, entity.IsPriority).Scan(&ParentKpiId)
					if err != nil {
						logengine.GetTelemetryClient().TrackException(err.Error())
						if !response.Error {
							response.Error = true
							response.ValidationErrors = []*model.ValidationMessage{}
						}
						errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to insert in kpi"}
						response.ValidationErrors = append(response.ValidationErrors, errorMessage)
						return response
					}

					updateBrandKPI := `update kpi set parent_kpi = $2 where id = $1`
					commandTag, err := tx.Exec(context.Background(), updateBrandKPI, ParentKpiId, ParentKpiId)
					if err != nil {
						logengine.GetTelemetryClient().TrackException(err.Error())
						if !response.Error {
							response.Error = true
							response.ValidationErrors = []*model.ValidationMessage{}
						}
						errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
						response.ValidationErrors = append(response.ValidationErrors, errorMessage)
						return response
					} else if commandTag.RowsAffected() != 1 {
						if !response.Error {
							response.Error = true
							response.ValidationErrors = []*model.ValidationMessage{}
						}
						errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
						response.ValidationErrors = append(response.ValidationErrors, errorMessage)
						return response
					} else {
						response.Message = "Kpi successfully Updated"
						response.Error = false
					}

					querystring = `INSERT INTO kpi_versions (design, kpi_id, is_active,is_deleted,created_by, date_created) VALUES($1, $2, $3, $4, $5, $6)`
					commandTag, err = tx.Exec(context.Background(), querystring, entity.BrandDesign, ParentKpiId, true, false, loggedInUserEntity.ID, timenow)
					if err != nil {
						logengine.GetTelemetryClient().TrackException(err.Error())
						if !response.Error {
							response.Error = true
							response.ValidationErrors = []*model.ValidationMessage{}
						}
						errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
						response.ValidationErrors = append(response.ValidationErrors, errorMessage)
						return response
					} else if commandTag.RowsAffected() != 1 {
						if !response.Error {
							response.Error = true
							response.ValidationErrors = []*model.ValidationMessage{}
						}
						errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
						response.ValidationErrors = append(response.ValidationErrors, errorMessage)
						return response
					} else {
						response.Message = "Kpi successfully Updated"
						response.Error = false
					}

				} else {
					updateBrandKPIquery := `update kpi set name = $2, last_modified = $3, modified_by = $4, target_items = $5  where parent_kpi = $1 and is_active = true and is_deleted = false and
				 "type" = (select id from code where category = 'KPIType' and value = 'brand' and is_active = true and is_delete = false)`
					commandTag, err := tx.Exec(context.Background(), updateBrandKPIquery, entity.ID, entity.Name, timenow, loggedInUserEntity.ID, entity.TargetBrand)
					if err != nil {
						logengine.GetTelemetryClient().TrackException(err.Error())
						if !response.Error {
							response.Error = true
							response.ValidationErrors = []*model.ValidationMessage{}
						}
						errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update brand kpi data"}
						response.ValidationErrors = append(response.ValidationErrors, errorMessage)
						return response
					}
					if commandTag.RowsAffected() == 0 {
						logengine.GetTelemetryClient().TrackException(err.Error())
						if !response.Error {
							response.Error = true
							response.ValidationErrors = []*model.ValidationMessage{}
						}
						errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update brand kpi data"}
						response.ValidationErrors = append(response.ValidationErrors, errorMessage)
						return response
					}
				}
			}

			if entity.TargetProduct != nil {

				if entity.ProductId == "" {
					querystring := `INSERT INTO kpi (name, type, target_team, target_items, is_active,is_deleted,created_by, date_created, month, year, is_priority, parent_kpi) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)RETURNING(id)`
					err = tx.QueryRow(context.Background(), querystring, entity.Name, 26, entity.TargetTeam, entity.TargetProduct, true, false, loggedInUserEntity.ID, timenow, entity.EffectiveMonth, entity.EffectiveYear, entity.IsPriority, entity.ID).Scan(&productId)
					if err != nil {
						logengine.GetTelemetryClient().TrackException(err.Error())
						if !response.Error {
							response.Error = true
							response.ValidationErrors = []*model.ValidationMessage{}
						}
						errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to insert in kpi"}
						response.ValidationErrors = append(response.ValidationErrors, errorMessage)
						return response
					}

					querystring = `INSERT INTO kpi_versions (design, kpi_id, is_active,is_deleted,created_by, date_created) VALUES($1, $2, $3, $4, $5, $6)`
					commandTag, err := tx.Exec(context.Background(), querystring, entity.ProductDesign, productId, true, false, loggedInUserEntity.ID, timenow)
					if err != nil {
						logengine.GetTelemetryClient().TrackException(err.Error())
						if !response.Error {
							response.Error = true
							response.ValidationErrors = []*model.ValidationMessage{}
						}
						errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
						response.ValidationErrors = append(response.ValidationErrors, errorMessage)
						return response
					} else if commandTag.RowsAffected() != 1 {
						logengine.GetTelemetryClient().TrackException(err.Error())
						if !response.Error {
							response.Error = true
							response.ValidationErrors = []*model.ValidationMessage{}
						}
						errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi data"}
						response.ValidationErrors = append(response.ValidationErrors, errorMessage)
						return response
					} else {
						response.Message = "Kpi successfully Updated"
						response.Error = false
					}

				} else {
					if ParentKpiId != "" {
						updateProdustKPIquery := `update kpi set name = $2, last_modified = $3, modified_by = $4, target_items = $5, parent_kpi = $6  where parent_kpi = $1 and is_active = true and is_deleted = false and
						"type" = (select id from code where category = 'KPIType' and value = 'product' and is_active = true and is_delete = false)`
						commandTag, err := tx.Exec(context.Background(), updateProdustKPIquery, entity.ID, entity.Name, timenow, loggedInUserEntity.ID, entity.TargetProduct, ParentKpiId)
						if err != nil {
							logengine.GetTelemetryClient().TrackException(err.Error())
							if !response.Error {
								response.Error = true
								response.ValidationErrors = []*model.ValidationMessage{}
							}
							errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update product kpi data"}
							response.ValidationErrors = append(response.ValidationErrors, errorMessage)
							return response
						}
						if commandTag.RowsAffected() != 1 {
							logengine.GetTelemetryClient().TrackException(err.Error())
							if !response.Error {
								response.Error = true
								response.ValidationErrors = []*model.ValidationMessage{}
							}
							errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update product kpi data"}
							response.ValidationErrors = append(response.ValidationErrors, errorMessage)
							return response
						}
					} else {
						updateProdustKPIquery := `update kpi set name = $2, last_modified = $3, modified_by = $4, target_items = $5, parent_kpi = $6  where parent_kpi = $1 and is_active = true and is_deleted = false and
						"type" = (select id from code where category = 'KPIType' and value = 'product' and is_active = true and is_delete = false)`
						commandTag, err := tx.Exec(context.Background(), updateProdustKPIquery, entity.ID, entity.Name, timenow, loggedInUserEntity.ID, entity.TargetProduct, entity.BrandId)
						if err != nil {
							logengine.GetTelemetryClient().TrackException(err.Error())
							if !response.Error {
								response.Error = true
								response.ValidationErrors = []*model.ValidationMessage{}
							}
							errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update product kpi data"}
							response.ValidationErrors = append(response.ValidationErrors, errorMessage)
							return response
						}
						if commandTag.RowsAffected() != 1 {
							logengine.GetTelemetryClient().TrackException(err.Error())
							if !response.Error {
								response.Error = true
								response.ValidationErrors = []*model.ValidationMessage{}
							}
							errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update product kpi data"}
							response.ValidationErrors = append(response.ValidationErrors, errorMessage)
							return response
						}
					}

				}
			}

			if entity.BrandId != "" {
				UpdateBrandKpiVersionquery := "update kpi_versions set is_active = $1, is_deleted = $2, last_modified = $3, modified_by = $4 where kpi_id = $5 and is_active = true and is_deleted = false"
				commandTag, err := tx.Exec(context.Background(), UpdateBrandKpiVersionquery, false, true, timenow, loggedInUserEntity.ID, entity.BrandId)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to delete brand kpi version data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				}
				if commandTag.RowsAffected() != 1 {
					logengine.GetTelemetryClient().TrackException(err.Error())
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to delete brand kpi version data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				}

				querystring := `INSERT INTO kpi_versions (design, kpi_id, is_active,is_deleted,created_by, date_created) VALUES($1, $2, $3, $4, $5, $6)`
				commandTag, err = tx.Exec(context.Background(), querystring, entity.BrandDesign, entity.BrandId, true, false, loggedInUserEntity.ID, timenow)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi version data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				} else if commandTag.RowsAffected() != 1 {
					logengine.GetTelemetryClient().TrackException(err.Error())
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update version kpi data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				} else {
					response.Message = "Kpi successfully Updated"
					response.Error = false
				}

			}

			if entity.ProductId != "" {
				UpdateProductKpiVersionquery := "update kpi_versions set is_active = $1, is_deleted = $2, last_modified = $3, modified_by = $4 where kpi_id = $5 and is_active = true and is_deleted = false"
				commandTag, err := tx.Exec(context.Background(), UpdateProductKpiVersionquery, false, true, timenow, loggedInUserEntity.ID, entity.ProductId)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update product kpi version data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				}
				if commandTag.RowsAffected() != 1 {
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update product kpi version data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				}

				querystring := `INSERT INTO kpi_versions (design, kpi_id, is_active,is_deleted,created_by, date_created) VALUES($1, $2, $3, $4, $5, $6)`
				commandTag, err = tx.Exec(context.Background(), querystring, entity.ProductDesign, entity.ProductId, true, false, loggedInUserEntity.ID, timenow)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi version data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				} else if commandTag.RowsAffected() != 1 {
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update kpi version data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				}

			}

		} else if entity.ID != nil && entity.IsDeleted == true {
			query := "update kpi set is_active = false, is_deleted = true, deleted_by = $2, date_deleted = $3 where parent_kpi = $1 and is_active = true and is_deleted = false"
			commandTag, err := tx.Exec(context.Background(), query, entity.ID, loggedInUserEntity.ID, timenow)
			if err != nil {
				logengine.GetTelemetryClient().TrackException(err.Error())
				if !response.Error {
					response.Error = true
					response.ValidationErrors = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to delete kpi data"}
				response.ValidationErrors = append(response.ValidationErrors, errorMessage)
				return response
			}
			if commandTag.RowsAffected() < 1 {
				if !response.Error {
					response.Error = true
					response.ValidationErrors = []*model.ValidationMessage{}
				}
				errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to delete kpi data"}
				response.ValidationErrors = append(response.ValidationErrors, errorMessage)
				return response
			} else {
				response.Message = "Kpi successfully deleted"
				response.Error = false
			}

			if entity.BrandId != "" {
				UpdateBrandKpiVersionquery := "update kpi_versions set is_active = false, is_deleted = true, last_modified =$2, modified_by =$3 where kpi_id = $1 and is_active = true and is_deleted = false"
				commandTag, err = tx.Exec(context.Background(), UpdateBrandKpiVersionquery, entity.BrandId, timenow, loggedInUserEntity.ID)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update brand kpi version data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				}
				if commandTag.RowsAffected() != 1 {
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update brand kpi version data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				}
			}

			if entity.ProductId != "" {
				UpdateProductKpiVersionquery := "update kpi_versions set is_active = false, is_deleted = true, last_modified =$2, modified_by =$3 where kpi_id = $1 and is_active = true and is_deleted = false"
				commandTag, err = tx.Exec(context.Background(), UpdateProductKpiVersionquery, entity.ProductId, timenow, loggedInUserEntity.ID)
				if err != nil {
					logengine.GetTelemetryClient().TrackException(err.Error())
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update product kpi version data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				}
				if commandTag.RowsAffected() != 1 {
					if !response.Error {
						response.Error = true
						response.ValidationErrors = []*model.ValidationMessage{}
					}
					errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to update product kpi version data"}
					response.ValidationErrors = append(response.ValidationErrors, errorMessage)
					return response
				}
			}
		}
	}
	txErr := tx.Commit(context.Background())
	if txErr != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		if !response.Error {
			response.Error = true
			response.ValidationErrors = []*model.ValidationMessage{}
		}
		errorMessage := &model.ValidationMessage{Row: 0, Message: "Failed to commit Kpi data"}
		response.ValidationErrors = append(response.ValidationErrors, errorMessage)
		return response
	}
	return response
}

func GetKpiInfo(input entity.GetKpisInput, loggedInUserEntity *entity.LoggedInUser, teams []string) ([]entity.KpiData, int, error) {
	if pool == nil {
		pool = GetPool()
	}

	totalPages := 0
	limit := 0
	var totalRecords int
	var rows pgx.Rows
	var err error
	var inputArgs []interface{}

	query := `with lookup as (
		select
			k.parent_kpi parent_kpi_id,
			(case
				when c.value = 'product' then k.id::text
				else null
			end) product_kpi_id,
			(case
				when c.value = 'brand' then k.id::text
				else null
			end) brand_kpi_id,
			(case
				when c.value = 'product' then kv.id::text
				else null
			end) product_kpi_version_id,
			(case
				when c.value = 'brand' then kv.id::text
				else null
			end) brand_kpi_version_id,
			k."name" kpi_name,
			k.target_team target_team_id,
			t.team_name target_team_name,
			k."month" effective_month,
			k."year" effective_year,
			k.is_priority,
			(case
				when c.value = 'product' then k.target_items
				else null
			end) target_product,
			(case
				when c.value = 'brand' then k.target_items
				else null
			end) target_brand,
			(case
				when c.value = 'product' then kv.design
				else null
			end) product_design,
			(case
				when c.value = 'brand' then kv.design
				else null
			end) brand_design,
			c.value as "type",
		    row_number() over(partition by k.target_team,
			c.value
		order by
			k."year" desc,
			k."month" desc) row_number
		from
			kpi k
		left join team t on
			t.id = k.target_team
		left join code c on
			k.type = c.id
		inner join kpi_versions kv on
			kv.kpi_id = k.id
		inner join "user" u on
			u.id = k.created_by
		left join sales_organisation so on
			so.id = u.sales_organisation
		where
			kv.is_active = true
			and kv.is_deleted = false
			and k.is_active = true
			and k.is_deleted = false
			and u.sales_organisation = ?`

	inputArgs = append(inputArgs, loggedInUserEntity.SalesOrganisaton)

	if teams != nil && len(teams) > 0 {
		var IDs []string
		for _, repId := range teams {
			IDs = append(IDs, "'"+repId+"'")
		}
		stringRepIds := strings.Trim(strings.Join(IDs, ", "), ", ")
		query = query + ` AND t.id IN (` + stringRepIds + `) `
	}

	if input.Month != nil && *input.Month > 0 && *input.Month < 13 {
		query += ` and k."month" <= ?`
		inputArgs = append(inputArgs, *input.Month)
	}

	if input.Year != nil && *input.Year != 0 {
		query += ` and k."year" <= ?`
		inputArgs = append(inputArgs, *input.Year)
	}

	if input.ParentKpiID != nil && *input.ParentKpiID != "" {
		query += ` and k.parent_kpi = ?`
		inputArgs = append(inputArgs, *input.ParentKpiID)
	}

	if input.TeamID != nil && *input.TeamID != "" {
		query += ` and k.target_team = ?`
		inputArgs = append(inputArgs, *input.TeamID)
	}

	if input.SearchItem != nil && *input.SearchItem != "" {
		query += ` and (k."name" ilike ?
			or t.team_name ilike ?)`
		inputArgs = append(inputArgs, "%"+*input.SearchItem+"%", "%"+*input.SearchItem+"%")
	}

	query += ` order by
			k.date_created desc,
			parent_kpi_id,
			brand_kpi_id,
			product_kpi_id),
		brand_kpis as (
		select
			*
		from
			lookup l
		where
			l."type" = 'brand'`

	if input.Month != nil && *input.Month > 0 && *input.Month < 13 && input.Year != nil && *input.Year != 0 {
		query += ` and l.row_number = 1`
	}

	query += `),
		product_kpis as (
		select
			*
		from
			lookup l
		where
			l."type" = 'product'`

	if input.Month != nil && *input.Month > 0 && *input.Month < 13 && input.Year != nil && *input.Year != 0 {
		query += ` and l.row_number = 1`
	}

	query += `),
		merged_kpis_brand as (
		select
			bk.parent_kpi_id,
			pk.product_kpi_id,
			bk.brand_kpi_id,
			pk.product_kpi_version_id,
			bk.brand_kpi_version_id,
			bk.kpi_name,
			bk.target_team_id,
			bk.target_team_name,
			bk.effective_month,
			bk.effective_year,
			bk.is_priority,
			pk.target_product,
			bk.target_brand,
			pk.product_design,
			bk.brand_design
		from
			brand_kpis bk
		left join product_kpis pk on
			pk.parent_kpi_id = bk.parent_kpi_id),
		merged_kpis_product as (
		select
			pk.parent_kpi_id,
			pk.product_kpi_id,
			bk.brand_kpi_id,
			pk.product_kpi_version_id,
			bk.brand_kpi_version_id,
			pk.kpi_name,
			pk.target_team_id,
			pk.target_team_name,
			pk.effective_month,
			pk.effective_year,
			pk.is_priority,
			pk.target_product,
			bk.target_brand,
			pk.product_design,
			bk.brand_design
		from
			product_kpis pk
		left join brand_kpis bk on
			bk.parent_kpi_id = pk.parent_kpi_id),
		merged_kpis as (
		select
			*
		from
			merged_kpis_brand
		union
		select
			*
		from
			merged_kpis_product ),
		merged_kpis_final as (
		select
			*
		from
			merged_kpis
		order by
			effective_year desc,
			effective_month desc,
			parent_kpi_id,
			brand_kpi_id,
			product_kpi_id)
		select
			*,
			count(mkf.parent_kpi_id) over()
		from
			merged_kpis_final mkf
		where
			1 = 1`

	if input.KpiID != nil && *input.KpiID != "" {
		query += ` and (mkf.product_kpi_id = ?
			or mkf.brand_kpi_id = ?)`
		inputArgs = append(inputArgs, *input.KpiID, *input.KpiID)
	}
	if input.KpiVersionID != nil && *input.KpiVersionID != "" {
		query += ` and (mkf.product_kpi_version_id = ?
			or mkf.brand_kpi_version_id = ?)`
		inputArgs = append(inputArgs, *input.KpiVersionID, *input.KpiVersionID)
	}

	if input.BrandID != nil && *input.BrandID != "" {
		query += ` and ? = any(mkf.target_brand)`
		inputArgs = append(inputArgs, *input.BrandID)
	}

	if input.TeamProductID != nil && *input.TeamProductID != "" {
		query += ` and ? = any(mkf.target_product)`
		inputArgs = append(inputArgs, *input.TeamProductID)
	}

	if input.Limit != nil {
		limit = *input.Limit

		query += ` limit ?`
		inputArgs = append(inputArgs, *input.Limit)
	}

	if input.Offset != nil {
		query += ` offset ?`
		inputArgs = append(inputArgs, *input.Offset)
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err = pool.Query(context.Background(), query, inputArgs...)

	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.KpiData{}, 0, err
	}
	kpis := []entity.KpiData{}
	defer rows.Close()
	for rows.Next() {
		kpi := entity.KpiData{}
		err := rows.Scan(
			&kpi.ParentKpiID,
			&kpi.ProductKpiID,
			&kpi.BrandKpiID,
			&kpi.ProductKpiVersionID,
			&kpi.BrandKpiVersionID,
			&kpi.KpiName,
			&kpi.TargetTeamID,
			&kpi.TargetTeamName,
			&kpi.EffectiveMonth,
			&kpi.EffectiveYear,
			&kpi.IsPriority,
			pq.Array(&kpi.TargetProduct),
			pq.Array(&kpi.TargetBrand),
			&kpi.ProductDesign,
			&kpi.BrandDesign,
			&totalRecords,
		)

		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return kpis, 0, err
		}

		kpis = append(kpis, kpi)
	}

	if len(kpis) < 1 {
		return kpis, 0, errors.New("No data found")
	}

	if limit > 0 {
		d := float64(totalRecords) / float64(limit)
		totalPages = int(math.Ceil(d))
	}

	return kpis, totalPages, err
}

func GetPictureRetrieveCustomerList(input *model.ListInput) ([]entity.CustomerList, error) {
	if pool == nil {
		pool = GetPool()
	}
	query := `select c.id,c."name", ka.answers from customer c 
	inner join team_member_customer tmc on tmc.customer = c.id 
	inner join kpi_answers ka on ka.team_member_customer = tmc.id
	inner join team_members tm on tm.id = tmc.team_member
	where c.is_active = true and c.is_deleted = false 
	and tmc.is_active = true and tmc.is_deleted = false 
	and tm.is_active = true and tm.is_deleted = false`
	var rows pgx.Rows
	var err error
	var inputArgs []interface{}
	if input != nil {
		if input.TeamID != nil && *input.TeamID != "" {
			query = query + ` and tm.team = ?`
			inputArgs = append(inputArgs, *input.TeamID)
		}
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err = pool.Query(context.Background(), query, inputArgs...)

	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.CustomerList{}, err
	}
	customers := []entity.CustomerList{}
	defer rows.Close()
	for rows.Next() {
		customer := entity.CustomerList{}
		err := rows.Scan(
			&customer.CustomerId,
			&customer.Name,
			&customer.Url,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return customers, err
		}
		customers = append(customers, customer)
	}
	return customers, err
}

func GetPictureRetrieveProductList(input *model.ListInput) ([]entity.ProductList, error) {
	if pool == nil {
		pool = GetPool()
	}
	query := `select p.id,p.material_description,ka.answers from product p
	inner join team_products tp on tp.material_code = p.id 
	inner join kpi_answers ka on ka.target_item = tp.id
	inner join team_member_customer tmc on tmc.id = ka.team_member_customer 
	inner join team_members tm on tm.id = tmc.team_member
	where p.is_active = true and p.is_deleted = false 
	and tp.is_active = true and tp.is_delete = false 
	and tmc.is_active = true and tmc.is_deleted = false
	and tm.is_active = true and tm.is_deleted = false`
	var rows pgx.Rows
	var err error
	var inputArgs []interface{}
	if input != nil {
		if input.TeamID != nil && *input.TeamID != "" {
			query = query + ` and tm.team = ?`
			inputArgs = append(inputArgs, *input.TeamID)
		}
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err = pool.Query(context.Background(), query, inputArgs...)

	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.ProductList{}, err
	}
	products := []entity.ProductList{}
	defer rows.Close()
	for rows.Next() {
		Product := entity.ProductList{}
		err := rows.Scan(
			&Product.ProductId,
			&Product.Name,
			&Product.Url,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return products, err
		}
		products = append(products, Product)
	}
	return products, err
}

func GetPictureRetrieveBrandList(input *model.ListInput) ([]entity.BrandList, error) {
	if pool == nil {
		pool = GetPool()
	}
	query := `select b.id,b.brand_name,ka.answers from brand b
	inner join kpi_answers ka on ka.target_item = b.id
	inner join team_member_customer tmc on tmc.id = ka.team_member_customer 
	inner join team_members tm on tm.id = tmc.team_member
	where b.is_active = true and b.is_deleted = false 
	and tmc.is_active = true and tmc.is_deleted = false 
	and tm.is_active = true and tm.is_deleted = false`
	var rows pgx.Rows
	var err error
	var inputArgs []interface{}
	if input != nil {
		if input.TeamID != nil && *input.TeamID != "" {
			query = query + ` and tm.team = ?`
			inputArgs = append(inputArgs, *input.TeamID)
		}
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err = pool.Query(context.Background(), query, inputArgs...)

	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.BrandList{}, err
	}
	brands := []entity.BrandList{}
	defer rows.Close()
	for rows.Next() {
		brand := entity.BrandList{}
		err := rows.Scan(
			&brand.BrandId,
			&brand.Name,
			&brand.Url,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return brands, err
		}
		brands = append(brands, brand)
	}
	return brands, err
}
func GetKpiBrandProduct(inputModel *entity.KpiBrandProductInput, loggedInUserEntity *entity.LoggedInUser) ([]entity.KpiBrandProductData, error) {
	if pool == nil {
		pool = GetPool()
	}
	var rows pgx.Rows
	var err error
	var inputArgs []interface{}
	date := time.Now()
	formatdate := date.Format("01/02/2006")
	month := formatdate[0:2]
	year := formatdate[6:10]
	if *inputModel.IsKpi == false {

		query := `with lookup as (
		select b.id as brand_id, b.brand_name, tp.id as team_products_id,  p.id as product_id, p.principal_name, p.material_description,
		(CASE WHEN tp.is_priority THEN tp.is_priority WHEN NOT tp.is_priority THEN p.is_priority END) as is_priority, tp.team, b.sales_organisation as salesorg
		from team_products tp 
		inner join product p on p.id = tp.material_code  
		inner join brand b on b.id = p.brand and p.sales_organisation = b.sales_organisation
		where b.sales_organisation = ?`
		inputArgs = append(inputArgs, loggedInUserEntity.SalesOrganisaton)

		if *inputModel.IsActive {
			query = query + ` and tp.is_active = true and tp.is_delete = false 
			and p.is_active = true and p.is_deleted = false 
			and b.is_active = true and b.is_deleted = false`
		}

		query = query + `)
		select brand_id, brand_name, team_products_id, product_id, principal_name, material_description, is_priority, ''as type
		from lookup where lookup.salesorg =?`

		inputArgs = append(inputArgs, loggedInUserEntity.SalesOrganisaton)

		if inputModel.TargetTeam != nil && *inputModel.TargetTeam != "" {
			query = query + ` and lookup.team = ?`
			inputArgs = append(inputArgs, *inputModel.TargetTeam)
		}
		if inputModel.BrandID != nil && *inputModel.BrandID != "" {
			query = query + ` and lookup.brand_id = ?`
			inputArgs = append(inputArgs, *inputModel.BrandID)
		}
		if inputModel.ProductId != nil && *inputModel.ProductId != "" {
			query = query + ` and lookup.product_id = ?`
			inputArgs = append(inputArgs, *inputModel.ProductId)
		}
		if inputModel.TeamProductId != nil && *inputModel.TeamProductId != "" {
			query = query + ` and lookup.team_products_id = ?`
			inputArgs = append(inputArgs, *inputModel.TeamProductId)
		}
		if inputModel.IsPriority != nil {
			query = query + ` and lookup.is_priority = ?`
			inputArgs = append(inputArgs, *inputModel.IsPriority)
		}
		if inputModel.SearchIteam != nil && *inputModel.SearchIteam != "" {
			query += ` and (lookup.principal_name ilike ?
			or lookup.material_description ilike ?
			or lookup.brand_name ilike ?)`
			inputArgs = append(inputArgs, "%"+*inputModel.SearchIteam+"%", "%"+*inputModel.SearchIteam+"%", "%"+*inputModel.SearchIteam+"%")
		}

		query = query + ` order by lookup.brand_id`
		query = sqlx.Rebind(sqlx.DOLLAR, query)
		rows, err = pool.Query(context.Background(), query, inputArgs...)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return []entity.KpiBrandProductData{}, err
		}
	} else {
		query := `with kpis as(
			select unnest (target_items) as item_id, target_team ,created_by, is_priority,kpi.date_created , c.value as kpi_type ,
			row_number() over(partition by kpi.target_team, c.value order by kpi."year" desc, kpi."month" desc) as row_number
			from kpi 
			inner join code c on c.id = kpi."type" and c.category = 'KPIType'
			where kpi.target_team = ? `
		inputArgs = append(inputArgs, inputModel.TargetTeam)
		query += `and kpi."month" <= ?
		and kpi."year" <= ?`
		inputArgs = append(inputArgs, month)
		inputArgs = append(inputArgs, year)

		if *inputModel.IsActive {
			query += ` and kpi.is_active = true and kpi.is_deleted = false`
		}

		query += `), expand as
			(select * 
			from kpis where kpis.row_number = 1),
			
			product_brand as
			(
			select  e.item_id, e.kpi_type,'' as material_description , '' principal_name , b.id as brand_id , b.brand_name, '00000000-0000-0000-0000-000000000000' as team_product_id, '00000000-0000-0000-0000-000000000000' as product_id,
			e.kpi_type as type,	(CASE WHEN tp.is_priority THEN tp.is_priority WHEN NOT tp.is_priority THEN p.is_priority END) as product_priority,
			 row_number() over (partition by e.item_id ) as row_number
			from expand e 
				inner join brand b on
					b.id = e.item_id
				inner join product p on
					p.brand = b.id
				inner join team_products tp on
					tp.material_code = p.id
				where tp.team = ?`
		inputArgs = append(inputArgs, inputModel.TargetTeam)
		if *inputModel.IsActive {
			query = query + ` and tp.is_active = true
				and tp.is_delete = false
				and p.is_active = true
				and p.is_deleted = false
				and b.is_active = true
				and b.is_deleted = false`
		}
		query += ` and e.kpi_type = 'brand'
		and row_number = 1
		union all                        
		select 
		e.item_id, e.kpi_type,p.material_description , p.principal_name , b.id as brand_id , b.brand_name, tp.id as team_product_id, p.id as product_id,
		e.kpi_type as "type", (CASE WHEN tp.is_priority THEN tp.is_priority WHEN NOT tp.is_priority THEN p.is_priority END) as product_priority,
		 row_number() over (partition by e.item_id ) as row_number
		from expand e 		
			inner join team_products  tp on tp.id = e.item_id
			inner join product p on p.id = tp.material_code 
			inner join brand b on b.id = p.brand 
			where tp.team = ?`

		inputArgs = append(inputArgs, inputModel.TargetTeam)
		if *inputModel.IsActive {
			query = query + ` and tp.is_active = true and tp.is_delete = false
							and p.is_active = true and p.is_deleted = false
							and b.is_active = true and b.is_deleted = false
							and e.kpi_type = 'product'
							and row_number = 1`
		}
		query += `)select  brand_id, brand_name, team_product_id, product_id, principal_name, material_description, product_priority, "type"
		from product_brand where row_number = 1`
		if inputModel.BrandID != nil && *inputModel.BrandID != "" {
			query = query + ` and product_brand.brand_id = ?`
			inputArgs = append(inputArgs, *inputModel.BrandID)
		}
		if inputModel.ProductId != nil && *inputModel.ProductId != "" {
			query = query + ` and product_brand.product_id = ?`
			inputArgs = append(inputArgs, *inputModel.ProductId)
		}
		if inputModel.TeamProductId != nil && *inputModel.TeamProductId != "" {
			query = query + ` and product_brand.team_products_id = ?`
			inputArgs = append(inputArgs, *inputModel.TeamProductId)
		}
		if inputModel.IsPriority != nil {
			query = query + ` and product_brand.product_priority = ?`
			inputArgs = append(inputArgs, *inputModel.IsPriority)
		}
		if inputModel.SearchIteam != nil && *inputModel.SearchIteam != "" {
			query += ` and (product_brand.principal_name ilike ?
			or product_brand.material_description ilike ?
			or product_brand.brand_name ilike ?)`
			inputArgs = append(inputArgs, "%"+*inputModel.SearchIteam+"%", "%"+*inputModel.SearchIteam+"%", "%"+*inputModel.SearchIteam+"%")
		}

		query = sqlx.Rebind(sqlx.DOLLAR, query)
		rows, err = pool.Query(context.Background(), query, inputArgs...)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return []entity.KpiBrandProductData{}, err
		}

	}

	kpiProducts := []entity.KpiBrandProductData{}
	defer rows.Close()
	for rows.Next() {
		kpiProd := entity.KpiBrandProductData{}
		err := rows.Scan(
			&kpiProd.BrandId,
			&kpiProd.BrandName,
			&kpiProd.TeamProductId,
			&kpiProd.ProductId,
			&kpiProd.PrincipalName,
			&kpiProd.MaterialDescription,
			&kpiProd.IsPriority,
			&kpiProd.Type,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return kpiProducts, err
		}
		kpiProducts = append(kpiProducts, kpiProd)
	}
	return kpiProducts, err

}
func HasTargetTeam(teamId string) (bool, error) {
	if pool == nil {
		pool = GetPool()
	}

	var isActive bool
	querystring := "select is_active from team t where t.id = $1 and t.is_active = true and t.is_deleted = false"
	err := pool.QueryRow(context.Background(), querystring, teamId).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func IsValidTeamProduct(teamproductid string) (bool, error) {
	if pool == nil {
		pool = GetPool()
	}
	var isActive bool
	querystring := "select is_active from team_products tp where tp.id = $1 and tp.is_active =true and tp.is_delete =false"
	err := pool.QueryRow(context.Background(), querystring, teamproductid).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func HasKpi(kpiID string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	queryString := "select 1 from kpi where id = $1 AND is_active = true AND is_deleted = false"
	var hasValue int
	err := pool.QueryRow(context.Background(), queryString, kpiID).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}

func GetValueFromCode(id int, category string) (string, error) {
	querystring := "SELECT value FROM code WHERE id = $1 and category = $2"
	var codeValue string
	err := pool.QueryRow(context.Background(), querystring, id, category).Scan(&codeValue)

	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return "", err
	}
	return codeValue, nil
}

func HasKPIVersion(kpiVerId string) bool {
	if pool == nil {
		pool = GetPool()
	}
	var result bool
	queryString := `select 1 from kpi_versions where id = $1 AND is_active = true AND is_deleted = false`
	var hasValue int
	err := pool.QueryRow(context.Background(), queryString, kpiVerId).Scan(&hasValue)
	if err == nil {
		result = true
	} else {
		logengine.GetTelemetryClient().TrackException(err.Error())
		result = false
	}
	return result
}

func GetKpiEvents(startDate, endDate, userID, sales_organisation string) ([]entity.KpiData, error) {
	if pool == nil {
		pool = GetPool()
	}
	var events []entity.KpiData
	var inputArgs []interface{}
	var err error
	query := `with lookup as (
		select
			k.parent_kpi parent_kpi_id,
			(case
				when c.value = 'product' then k.id::text
				else null
			end) product_kpi_id,
			(case
				when c.value = 'brand' then k.id::text
				else null
			end) brand_kpi_id,
			(case
				when c.value = 'product' then kv.id::text
				else null
			end) product_kpi_version_id,
			(case
				when c.value = 'brand' then kv.id::text
				else null
			end) brand_kpi_version_id,
			k.name kpi_name,
			k.target_team target_team_id,
			t.team_name target_team_name,
			k.month effective_month,
			k.year effective_year,
			k.is_priority,
			(case
				when c.value = 'product' then k.target_items
				else null
			end) target_product,
			(case
				when c.value = 'brand' then k.target_items
				else null
			end) target_brand,
			(case
				when c.value = 'product' then kv.design
				else null
			end) product_design,
			(case
				when c.value = 'brand' then kv.design
				else null
			end) brand_design,
			c.value as "type"
		from
			kpi k
		left join team t on
			t.id = k.target_team
		left join code c on
			k.type = c.id
		inner join kpi_versions kv on
			kv.kpi_id = k.id
		inner join "user" u on
			u.id = k.created_by
		left join sales_organisation so on
			so.id = u.sales_organisation
		inner join team_members tm on
		tm.team = k.target_team 
		inner join schedule_plan sp on 
		sp.team_member = tm.id 
		inner join schedule_event se 
		on se.schedule_plan_id = sp.id
		where
			kv.is_active = true
			and kv.is_deleted = false
			and k.is_active = true
			and k.is_deleted = false
			and u.sales_organisation = ?`

	inputArgs = append(inputArgs, sales_organisation)
	if startDate != "" && endDate != "" {
		if startDate != endDate {
			query = query + ` and (case when se.rescheduled_date is null then date(se.event_date) else date(se.rescheduled_date) end) between ? and ?`
			inputArgs = append(inputArgs, startDate, endDate)
		} else {
			query = query + ` and (case when se.rescheduled_date is null then date(se.event_date) else date(se.rescheduled_date) end) = ?`
			inputArgs = append(inputArgs, startDate)
		}
	} else {
		if startDate != "" {
			query = query + ` and (case when se.rescheduled_date is null then date(se.event_date) else date(se.rescheduled_date) end) >= ?`
			inputArgs = append(inputArgs, startDate)
		}
		if endDate != "" {
			query = query + ` and (case when se.rescheduled_date is null then date(se.event_date) else date(se.rescheduled_date) end) <= ?`
			inputArgs = append(inputArgs, endDate)
		}
	}

	query = query + `  order by
		k.date_created desc,
		parent_kpi_id,
		brand_kpi_id,
		product_kpi_id
		),
	brand_kpis as (
	select
		*
	from
		lookup l
	where
		l."type" = 'brand' ),
	product_kpis as (
	select
		*
	from
		lookup l
	where
		l."type" = 'product' ),
	merged_kpis_brand as (
	select
		bk.parent_kpi_id,
		pk.product_kpi_id,
		bk.brand_kpi_id,
		pk.product_kpi_version_id,
		bk.brand_kpi_version_id,
		bk.kpi_name,
		bk.target_team_id,
		bk.target_team_name,
		bk.effective_month,
		bk.effective_year,
		bk.is_priority,
		pk.target_product,
		bk.target_brand,
		pk.product_design,
		bk.brand_design
	from
		brand_kpis bk
	left join product_kpis pk on
		pk.parent_kpi_id = bk.parent_kpi_id),
	merged_kpis_product as (
	select
		pk.parent_kpi_id,
		pk.product_kpi_id,
		bk.brand_kpi_id,
		pk.product_kpi_version_id,
		bk.brand_kpi_version_id,
		pk.kpi_name,
		pk.target_team_id,
		pk.target_team_name,
		pk.effective_month,
		pk.effective_year,
		pk.is_priority,
		pk.target_product,
		bk.target_brand,
		pk.product_design,
		bk.brand_design
	from
		product_kpis pk
	left join brand_kpis bk on
		bk.parent_kpi_id = pk.parent_kpi_id),
	merged_kpis as (
	select
		*
	from
		merged_kpis_brand
	union
	select
		*
	from
		merged_kpis_product ),
	merged_kpis_final as (
	select
		*
	from
		merged_kpis
	order by
		parent_kpi_id,
		brand_kpi_id,
		product_kpi_id)
	select
		*
	from
		merged_kpis_final mkf
	where
		1 = 1`

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := pool.Query(context.Background(), query, inputArgs...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return events, err
	}
	kpis := []entity.KpiData{}
	defer rows.Close()
	for rows.Next() {
		kpi := entity.KpiData{}
		err := rows.Scan(
			&kpi.ParentKpiID,
			&kpi.ProductKpiID,
			&kpi.BrandKpiID,
			&kpi.ProductKpiVersionID,
			&kpi.BrandKpiVersionID,
			&kpi.KpiName,
			&kpi.TargetTeamID,
			&kpi.TargetTeamName,
			&kpi.EffectiveMonth,
			&kpi.EffectiveYear,
			&kpi.IsPriority,
			pq.Array(&kpi.TargetProduct),
			pq.Array(&kpi.TargetBrand),
			&kpi.ProductDesign,
			&kpi.BrandDesign,
		)

		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return kpis, err
		}

		kpis = append(kpis, kpi)
	}
	return kpis, err
}

func GetKpiProductBrandAnswer(startDate string, endDate string, userID string, endMonth, year string) ([]entity.KpiProductBrandAnswer, error) {
	if pool == nil {
		pool = GetPool()
	}

	var events []entity.KpiProductBrandAnswer
	var err error
	currentMonth := util.GetFirstDayOfCurrentMonth()
	nextMonth := currentMonth.AddDate(0, 1, 0)

	query := `with _kpi as (
			select
				unnest(target_items) as item_id,
				se.id as schedule_event_id,
				se.team_member_customer,
				k.id as kpi_id,
    			k.parent_kpi as parent_kpi_id,
				c.value as kpi_type,
    			k."month",
				k."year",
				tm.team as team_id,
				row_number() over (partition by unnest(k.target_items),
				se.id
			order by
				k."month" desc) as row_number
			from
				schedule_event se
			inner join schedule_plan sp on
				sp.id = se.schedule_plan_id
		    inner join code c2 on
				c2.id = sp.status
			inner join team_member_customer tmc on
				tmc.id = se.team_member_customer
			inner join team_members tm on
				tm.id = tmc.team_member
			inner join kpi k on
				k.target_team = tm.team
			inner join code c on
				c.id = k."type"
				and c.category = 'KPIType'
			inner join kpi_versions kv on
				kv.kpi_id = k.id
			where
				tm.employee = ?`

	var inputArgs []interface{}

	inputArgs = append(inputArgs, userID)

	if startDate != "" && endDate != "" {
		if startDate != endDate {
			query = query + ` and (case when se.rescheduled_date is null then date(se.event_date) else date(se.rescheduled_date) end) between ? and ?`
			inputArgs = append(inputArgs, startDate, endDate)
		} else {
			query = query + ` and (case when se.rescheduled_date is null then date(se.event_date) else date(se.rescheduled_date) end) = ?`
			inputArgs = append(inputArgs, startDate)
		}
	} else {
		if startDate != "" {
			query = query + ` and (case when se.rescheduled_date is null then date(se.event_date) else date(se.rescheduled_date) end) >= ?`
			inputArgs = append(inputArgs, startDate)
		}
		if endDate != "" {
			query = query + ` and (case when se.rescheduled_date is null then date(se.event_date) else date(se.rescheduled_date) end) <= ?`
			inputArgs = append(inputArgs, endDate)
		}
	}

	query += ` and se.team_member_customer is not null
				and se.is_active = true
				and se.is_deleted = false
				and sp.is_active = true
				and sp.is_deleted = false
				and tmc.is_active = true
				and tmc.is_deleted = false
				and tm.is_active = true
				and tm.is_deleted = false
				and k.is_active = true
				and k.is_deleted = false
				and c.is_active = true
				and c.is_delete = false
				and c2.is_active = true
				and c2.is_delete = false
				and kv.is_active = true
				and kv.is_deleted = false
				and c2.value = 'approved'
				and c2.category = 'ScheduleStatus'
				and k."month" <= ?
				and k."year" <= ?
			order by
				schedule_event_id,
				kpi_type),`

	inputArgs = append(inputArgs, endMonth)
	inputArgs = append(inputArgs, year)

	query += `preScheduleKpi as (
			select
				schedule_event_id,
				parent_kpi_id,
				kpi_id,
				kpi_type,
				"month",
				"year",
				row_number() over(partition by schedule_event_id,
				kpi_type
			order by
				"year" desc,
				"month" desc)
			from
				_kpi
			group by
				schedule_event_id,
				kpi_id,
				parent_kpi_id,
				kpi_type,
				"month",
				"year" ),
			brandScheduleKpi as (
			select
				schedule_event_id,
				kpi_id,
				parent_kpi_id,
				kpi_type,
				"month",
				"year"
			from
				preScheduleKpi
			where
				kpi_type = 'brand'
				and row_number = 1 ),
			prodScheduleKpi as (
			select
				schedule_event_id,
				kpi_id,
				parent_kpi_id,
				kpi_type,
				"month",
				"year"
			from
				preScheduleKpi
			where
				kpi_type = 'product'
				and row_number = 1 ),
			merged_kpis_brand as (
			select
				s1.schedule_event_id,
				s1.kpi_id brand_kpi_id,
				s2.kpi_id product_kpi_id,
				s1.parent_kpi_id,
				s1."month",
				s1."year"
			from
				brandScheduleKpi s1
			left join prodScheduleKpi s2 on
				s2.schedule_event_id = s1.schedule_event_id
				and s2."month" = s1."month"
				and s2."year" = s1."year"),
			merged_kpis_product as (
			select
				s1.schedule_event_id,
				s2.kpi_id brand_kpi_id,
				s1.kpi_id product_kpi_id,
				s1.parent_kpi_id,
				s1."month",
				s1."year"
			from
				prodScheduleKpi s1
			left join brandScheduleKpi s2 on
				s2.schedule_event_id = s1.schedule_event_id
				and s2."month" = s1."month"
				and s2."year" = s1."year"),
			unifiedMergedKpis as (
			select
				*
			from
				merged_kpis_brand
			union
			select
				*
			from
				merged_kpis_product ),
			preMergedKpis as (
			select
				schedule_event_id,
				brand_kpi_id,
				product_kpi_id,
				row_number() over (partition by schedule_event_id
			order by
				"year" desc,
				"month" desc)
			from
				unifiedMergedKpis ),
			mergedKpis as (
			select
				schedule_event_id,
				brand_kpi_id,
				product_kpi_id
			from
				preMergedKpis
			where
				row_number = 1 ),
			_productBrand as (
			select
				e.schedule_event_id,
				e.team_member_customer,
				b.id item_id,
				b.brand_name as item_name,
				p.material_description as item_description,
				p.principal_name as product_principal_name,
				b.brand_name as brand_name,
				(case
					when tp.is_priority then tp.is_priority
					when not tp.is_priority then p.is_priority
				end) as product_priority,
				p.id as product_id,
				tp.team as team_id,
				tp.id as team_product_id,
				e.kpi_id,
				e.kpi_type,
				(
				select
					kv.id
				from
					kpi k
				inner join team t on
					t.id = k.target_team
				inner join kpi_versions kv on
					kv.kpi_id = k.id
				where
					kv.is_active = true
					and kv.is_deleted = false
					and t.is_active = true
					and t.is_deleted = false
					and k.is_active = true
					and k.is_deleted = false
					and t.id = e.team_id
					and e.item_id = any(k.target_items)
					and to_Date(k.month::varchar || ' ' || k.year::varchar, 'mm YYYY') < ?
				order by
					year desc,
					month desc
				limit 1) as kpi_version_id,
				ka.id as kpi_ans_id,
				ka.answers,
				ka.category as kpi_answer_category,
				ka.target_item,
				ka.kpi_version_id as ans_kpi_version_id,
				ka.schedule_event_id as ans_schedule_event_id,
				ka.team_member_customer as ans_team_member_customer,
				e."month" as months,
				e."year" as years,
				row_number() over (partition by e.schedule_event_id,
				tp.id,
				ka.category
			order by
				ka.date_created desc) as row_number
			from
				_kpi e
			inner join team_products tp on
				tp.id = e.item_id
			inner join product p on
				p.id = tp.material_code
			inner join brand b on
				b.id = p.brand
			inner join team_members tm2 on
				tm2.team = tp.team
			left join kpi_answers ka on
				ka.schedule_event_id = e.schedule_event_id
				and ka.target_item = tp.id
			where
				tp.is_active = true
				and tp.is_delete = false
				and p.is_active = true
				and p.is_deleted = false
				and b.is_active = true
				and b.is_deleted = false
				and e.kpi_type = 'product'
				and tm2.employee = ?`

	inputArgs = append(inputArgs, nextMonth)
	inputArgs = append(inputArgs, userID)

	query += ` and row_number = 1 ),
			prodAnswers as (
			select
				schedule_event_id,
				team_member_customer,
				item_id,
				item_name,
				item_description,
				product_priority,
				product_principal_name,
				brand_name,
				product_id::text,
				team_product_id::text,
				team_id,
				kpi_id,
				'product' kpi_type,
				kpi_version_id,
				kpi_ans_id,
				answers,
				kpi_answer_category,
				target_item,
				ans_kpi_version_id,
				ans_schedule_event_id,
				ans_team_member_customer,
				months,
				years
			from
				_productBrand
			where
				row_number = 1),
			prodAnswersFinal as (
			select
				p.schedule_event_id,
				team_member_customer,
				item_id,
				item_name,
				item_description,
				product_priority,
				product_principal_name,
				brand_name,
				product_id,
				team_product_id,
				team_id,
				kpi_id,
				kpi_type,
				(case
					when m.product_kpi_id is not null then kpi_version_id::text
					else kpi_version_id::text
				end) kpi_version_id,
				kpi_ans_id,
				answers,
				kpi_answer_category,
				target_item,
				ans_kpi_version_id,
				ans_schedule_event_id,
				ans_team_member_customer,
				months,
				years
			from
				prodAnswers p
			left join mergedKpis m on
				m.schedule_event_id = p.schedule_event_id
				and m.product_kpi_id = p.kpi_id),
			_brandProduct as (
			select
				e.schedule_event_id,
				e.team_member_customer,
				e.item_id,
				b.brand_name as item_name,
				'' as item_description,
				'' as product_principal_name,
				b.brand_name as brand_name,
				false as product_priority,
				'' as product_id,
				tm2.team as team_id,
				'' as team_product_id,
				e.kpi_id,
				e.kpi_type,
				(
				select
					kv.id
				from
					kpi k
				inner join team t on
					t.id = k.target_team
				inner join kpi_versions kv on
					kv.kpi_id = k.id
				where
					kv.is_active = true
					and kv.is_deleted = false
					and t.is_active = true
					and t.is_deleted = false
					and k.is_active = true
					and k.is_deleted = false
					and t.id = e.team_id
					and e.item_id = any(k.target_items)
					and to_Date(k.month::varchar || ' ' || k.year::varchar, 'mm YYYY') < ?
				order by
					year desc,
					month desc
				limit 1) as kpi_version_id,
				ka.id as kpi_ans_id,
				ka.answers,
				ka.category as kpi_answer_category,
				ka.target_item,
				ka.kpi_version_id as ans_kpi_version_id,
				ka.schedule_event_id as ans_schedule_event_id,
				ka.team_member_customer as ans_team_member_customer,
				e."month" as months,
				e."year" as years,
				row_number() over (partition by e.schedule_event_id,
				b.id,
				ka.category
			order by
				ka.date_created desc) as row_number
			from
				_kpi e
			inner join brand b on
				b.id = e.item_id
			inner join team_members tm2 on
				tm2.team = e.team_id
			left join kpi_answers ka on
				ka.schedule_event_id = e.schedule_event_id
				and ka.target_item = e.item_id
			where
				b.is_active = true
				and b.is_deleted = false
				and e.kpi_type = 'brand'
				and tm2.employee = ?`

	inputArgs = append(inputArgs, nextMonth)
	inputArgs = append(inputArgs, userID)

	query += ` and row_number = 1),
			brandAnswers as (
			select
				schedule_event_id,
				team_member_customer,
				item_id,
				item_name,
				item_description,
				product_priority,
				product_principal_name,
				brand_name,
				product_id,
				team_product_id,
				team_id,
				kpi_id,
				kpi_type,
				kpi_version_id,
				kpi_ans_id,
				answers,
				kpi_answer_category,
				target_item,
				ans_kpi_version_id,
				ans_schedule_event_id,
				ans_team_member_customer,
				months,
				years
			from
				_brandProduct
			where
				kpi_ans_id is null
				or row_number = 1),
			brandAnswersFinal as (
			select
				b.schedule_event_id,
				team_member_customer,
				item_id,
				item_name,
				item_description,
				product_priority,
				product_principal_name,
				brand_name,
				product_id,
				team_product_id,
				team_id,
				kpi_id,
				kpi_type,
				(case
					when m.brand_kpi_id is not null then kpi_version_id::text
					else kpi_version_id::text
				end) kpi_version_id,
				kpi_ans_id,
				answers,
				kpi_answer_category,
				target_item,
				ans_kpi_version_id,
				ans_schedule_event_id,
				ans_team_member_customer,
				months,
				years
			from
				brandAnswers b
			left join mergedKpis m on
				m.schedule_event_id = b.schedule_event_id
				and m.brand_kpi_id = b.kpi_id),
			finalData as (
			select
				*
			from
				brandAnswersFinal
			union all
			select
				*
			from
				prodAnswersFinal)
			select
				*
			from
				finalData
			order by
				schedule_event_id,
				kpi_type`

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := pool.Query(context.Background(), query, inputArgs...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return events, err
	}
	defer rows.Close()
	for rows.Next() {
		var event entity.KpiProductBrandAnswer
		err = rows.Scan(
			&event.ScheduleEventId,
			&event.TeamMemberCustomer,
			&event.ItemId,
			&event.ItemName,
			&event.ItemDescription,
			&event.ProductIsPriority,
			&event.ProductPrincipalName,
			&event.BrandName,
			&event.ProductId,
			&event.TeamProductID,
			&event.TeamID,
			&event.KpiId,
			&event.KpiType,
			&event.KpiVersionId,
			&event.KpiAnsId,
			&event.Answers,
			&event.Category,
			&event.TargetItem,
			&event.AnsKpiVersionId,
			&event.AnsScheduleEventId,
			&event.AnsTeamMemberCustomer,
			&event.Month,
			&event.Year,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return events, err
		}
		events = append(events, event)
	}
	return events, err
}

func GetTeams(user string, teams []string) []string {
	if pool == nil {
		pool = GetPool()
	}
	querystring := `select tm.team
	from team_members tm
	where tm.employee = $1
	and tm.approval_role = (select id from code where code.value = 'line1manager' and code.category = 'ApprovalRole')
	and tm.is_active = true and tm.is_deleted = false `
	rows, err := pool.Query(context.Background(), querystring, user)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return teams
	}
	defer rows.Close()
	for rows.Next() {
		var team string
		err := rows.Scan(
			&team,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return teams
		}
		teams = append(teams, team)
	}
	return teams

}
