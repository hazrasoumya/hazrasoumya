package postgres

import (
	"context"
	"sort"
	"time"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	"github.com/eztrade/kpi/graph/util"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
)

func InsertKpiAnswer(entity *entity.KpiAnswer, response *model.KpiResponse, loggedInUserEntity *entity.LoggedInUser) {
	if pool == nil {
		pool = GetPool()
	}
	tx, err := pool.Begin(context.Background())
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		response.Message = "Failed to begin transaction"
		response.ErrorCode = 999
		response.Error = true
	}
	defer tx.Rollback(context.Background())
	timenow := util.GetCurrentTime()
	querystring := "INSERT INTO kpi_answers (answers, created_by, date_created, kpi_version_id, category, team_member_customer, target_item, schedule_event_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8)"
	commandTag, err := tx.Exec(context.Background(), querystring, entity.Answers, loggedInUserEntity.ID, timenow, entity.KpiVersionId, entity.Category, entity.TeamMemberCustomerID, entity.TargetItem, entity.ScheduleEvent)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		response.Message = "Failed to create kpi answers"
		response.ErrorCode = 110
		response.Error = true
		return
	}
	if commandTag.RowsAffected() != 1 {
		response.Message = "Invalid Kpi answers data"
		response.ErrorCode = 111
		response.Error = true
		return
	} else {
		response.Message = "Kpi answers successfully inserted"
		response.Error = false
	}
	txErr := tx.Commit(context.Background())
	if txErr != nil {
		response.Message = "Failed to commit Kpi answers data"
		response.ErrorCode = 112
		response.Error = true
	}
}

func GetKpiAnswerDetails(inputModel *model.GetKpiAnswerRequest, loggedInUserEntity *entity.LoggedInUser, isFirstHit bool) ([]entity.KpiAnswerData, bool, error) {
	if pool == nil {
		pool = GetPool()
	}

	answerMap := make(map[string]*entity.KpiAnswerData)
	kpiAnswers := []entity.KpiAnswerData{}

	var rows pgx.Rows
	var err error

	if isFirstHit {
		query := `with _events as(SELECT ka.id, ka.answers, k.id as kpi_id, ka.kpi_version_id, ka.schedule_event_id, ka.category, ka.date_created, ka.team_member_customer, ka.target_item, 
			row_number() over (partition by ka.team_member_customer, ka.schedule_event_id , ka.category order by ka.date_created desc) as row_number, kv.design 
			FROM kpi_answers ka 
			INNER JOIN kpi_versions kv ON kv.id = ka.kpi_version_id 
			INNER JOIN kpi k on k.id = kv.kpi_id
			INNER join "user" u on u.id = ka.created_by
			INNER join sales_organisation so on so.id = u.sales_organisation 
			WHERE k.is_active = true AND k.is_deleted = false 
			AND kv.is_active = true AND kv.is_deleted = false
			AND u.is_active = true AND u.is_deleted = false 
			AND so.is_active =true AND so.is_deleted = false
			and u.sales_organisation =?`

		var inputArgs []interface{}
		inputArgs = append(inputArgs, loggedInUserEntity.SalesOrganisaton)
		if inputModel != nil {
			if inputModel.KpiID != nil && *inputModel.KpiID != "" {
				query = query + ` and k.id = ?`
				inputArgs = append(inputArgs, *inputModel.KpiID)
			}
			if inputModel.KpiVersionID != nil && *inputModel.KpiVersionID != "" {
				query = query + ` and ka.kpi_version_id = ?`
				inputArgs = append(inputArgs, *inputModel.KpiVersionID)
			}
			if inputModel.Category != nil && *inputModel.Category != 0 {
				query = query + ` and ka.category = ?`
				inputArgs = append(inputArgs, *inputModel.Category)
			}
			if inputModel.TeamMemberCustomerID != nil && *inputModel.TeamMemberCustomerID != "" {
				query = query + ` and ka.team_member_customer = ?`
				inputArgs = append(inputArgs, *inputModel.TeamMemberCustomerID)
			}
			if inputModel.ScheduleEvent != nil && *inputModel.ScheduleEvent != "" {
				query = query + ` and ka.schedule_event_id = ?`
				inputArgs = append(inputArgs, *inputModel.ScheduleEvent)
			}
			if inputModel.TargetItem != nil && *inputModel.TargetItem != "" {
				query = query + ` and ka.target_item = ?`
				inputArgs = append(inputArgs, *inputModel.TargetItem)
			}
		}

		query += `)
			select id, answers, kpi_id, kpi_version_id, schedule_event_id, category, date_created, team_member_customer, target_item, 
			       (case when position('Proposed Stock' in design::text) > 0 then true else false end) as is_proposed_stock 
			from _events
			where row_number = 1`

		query = sqlx.Rebind(sqlx.DOLLAR, query)

		rows, err = pool.Query(context.Background(), query, inputArgs...)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return []entity.KpiAnswerData{}, false, err
		}
	} else {
		query := `with _events as(
			SELECT 
				ka.id, 
				ka.answers, 
				k.id as kpi_id, 
				ka.kpi_version_id, 
				ka.schedule_event_id, 
				ka.category, 
				ka.date_created, 
				ka.team_member_customer, 
				ka.target_item, 
				se.event_date, 
				row_number() over (partition by ka.team_member_customer, ka.schedule_event_id , ka.category order by ka.date_created desc) as row_number, 
  				kv.design 
			FROM kpi_answers ka 
				INNER JOIN kpi_versions kv ON kv.id = ka.kpi_version_id 
				INNER JOIN kpi k on k.id = kv.kpi_id
				inner join schedule_event se on se.id = ka.schedule_event_id
				INNER join "user" u on u.id = ka.created_by
				INNER join sales_organisation so on so.id = u.sales_organisation 
			WHERE u.is_active = true AND u.is_deleted = false 
				AND so.is_active =true AND so.is_deleted = false
				and u.sales_organisation = ? 
				and ka.category = ? 
				and ka.team_member_customer = ? 
				and ka.target_item = ? 
		)
		select id, answers, kpi_id, kpi_version_id, schedule_event_id, category, date_created, team_member_customer, target_item, 
		       (case when position('Proposed Stock' in design::text) > 0 then true else false end) as is_proposed_stock 
		from _events
		where row_number = 1
		order by event_date desc 
		limit 1`

		query = sqlx.Rebind(sqlx.DOLLAR, query)

		var inputArgs []interface{}
		inputArgs = append(inputArgs, loggedInUserEntity.SalesOrganisaton, *inputModel.Category, *inputModel.TeamMemberCustomerID, *inputModel.TargetItem)

		rows, err = pool.Query(context.Background(), query, inputArgs...)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return []entity.KpiAnswerData{}, false, err
		}
	}

	var isProposedStock bool

	for rows.Next() {
		var id pgtype.UUID
		var answer string
		var kpiId pgtype.UUID
		var kpiVersionId pgtype.UUID
		var teamMemberCustomerId pgtype.UUID
		var scheduleEventId pgtype.UUID
		var targetItem pgtype.UUID
		var category int64
		var dateCreated time.Time

		err := rows.Scan(&id, &answer, &kpiId, &kpiVersionId, &scheduleEventId, &category, &dateCreated, &teamMemberCustomerId, &targetItem, &isProposedStock)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
		}
		idByte := id.Get().([16]byte)
		ansIdStr := util.UUIDV4ToString(idByte)

		var kpiIdStr string
		if value, _ := kpiId.Value(); value != nil {
			kpiIdByte := kpiId.Get().([16]byte)
			kpiIdStr = util.UUIDV4ToString(kpiIdByte)
		}

		var kpiVerIdStr string
		if value, _ := kpiId.Value(); value != nil {
			kpiVerIdByte := kpiVersionId.Get().([16]byte)
			kpiVerIdStr = util.UUIDV4ToString(kpiVerIdByte)
		}

		var scheduleEventIdStr string
		if value, _ := scheduleEventId.Value(); value != nil {
			scheduleEventIdByte := scheduleEventId.Get().([16]byte)
			scheduleEventIdStr = util.UUIDV4ToString(scheduleEventIdByte)
		}

		timeStampString := util.GetTimeUnixTimeStamp(dateCreated)
		compositeKey := timeStampString + "-" + ansIdStr

		var tmcIdStr string
		if value, _ := teamMemberCustomerId.Value(); value != nil {
			tmcIdByte := teamMemberCustomerId.Get().([16]byte)
			tmcIdStr = util.UUIDV4ToString(tmcIdByte)
		}

		var targetItemStr string
		if value, _ := targetItem.Value(); value != nil {
			targetItemByte := targetItem.Get().([16]byte)
			targetItemStr = util.UUIDV4ToString(targetItemByte)
		}

		var kpiAnswer *entity.KpiAnswerData

		if _, value := answerMap[compositeKey]; !value {
			kpiAnswer = &entity.KpiAnswerData{ID: ansIdStr, Answer: answer, KpiId: kpiIdStr, KpiVersionId: kpiVerIdStr, Category: category, TeamMemberCustomerID: tmcIdStr, ScheduleEventID: scheduleEventIdStr, TargetItem: targetItemStr}
			answerMap[compositeKey] = kpiAnswer
		} else {
			kpiAnswer = answerMap[compositeKey]
		}
	}

	keys := make([]string, 0, len(answerMap))
	for k := range answerMap {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[j] < keys[i]
	})
	for _, key := range keys {
		kpiAnswers = append(kpiAnswers, *answerMap[key])
	}

	return kpiAnswers, isProposedStock, err
}

func RetrievePicturesInfo(retrievePictureinput model.PicturesInput, pictureZip model.PictureZipInput, flag bool, loggedInUserEntity *entity.LoggedInUser, getTeamIDS []uuid.UUID, monthFilter bool, startDate string, yearFilter bool, endDate string, productBrandCustomerId []string) ([]entity.PictureData, error) {
	if pool == nil {
		pool = GetPool()
	}
	var query string
	var inputArgs []interface{}
	var onlySalesRep bool
	query = `with _events as (select distinct(tm.team), ka.answers, 
		(case 
			when c.value in ('permanent') then 
			(case 
				when ka.target_item in (select id from team_products tps where tps.is_active = true and tps.is_delete = false) then 'product'
				when ka.target_item in (select id from brand where is_active = true and is_deleted = false) then 'brand'
			end)
			else c.value 
		end) as type,
		(case when (ka.target_item not in (select id from team_products tps where tps.is_active = true and tps.is_delete = false)) then 
		(select b2.brand_name from brand b2
		where ka.target_item = b2.id)
		else 
		(select p2.material_description from product p2 
		inner join team_products tp2 on tp2.material_code = p2.id 
		where tp2.id = ka.target_item
		and p2.is_active = true and p2.is_deleted = false)
		end) as name,tmc.customer,
		row_number() over (partition by ka.team_member_customer, ka.schedule_event_id, ka.target_item, ka.category order by ka.date_created desc) as row_number
		from kpi_answers ka 
		inner join code c on ka.category = c.id
		inner join kpi_versions kv on kv.id = ka.kpi_version_id 
		inner join kpi kp on kp.id = kv.kpi_id 
		inner join team_member_customer tmc on ka.team_member_customer = tmc.id
		inner join team_members tm on tm.id = tmc.team_member 
		where c.value in('permanent','competitor','survey','promotion') and c.category = 'KPICategory' `
	if !onlySalesRep {
		query = query + `and (tm.team in(?)`
		inputArgs = append(inputArgs, getTeamIDS)
		if loggedInUserEntity.AuthRole != "sfe" {
			query = query + ` or tm.employee = ?`
			inputArgs = append(inputArgs, loggedInUserEntity.ID)
		}
		query = query + `) `
	} else {
		query = query + ` and tm.employee = ?`
		inputArgs = append(inputArgs, loggedInUserEntity.ID)
	}
	query = query + ` and tmc.is_active = true and tmc.is_deleted = false 
		and tm.is_active = true and tm.is_deleted = false
		and kv.is_active =true and kv.is_deleted = false and kp.is_deleted = false`
	var rows pgx.Rows
	var err error
	var dataForUpsert []string
	for key, value := range productBrandCustomerId {
		if key == 0 {
			// dataForUpsert = append(dataForUpsert, "'")
			dataForUpsert = append(dataForUpsert, "'"+value+"'")
			// dataForUpsert = append(dataForUpsert, "'")
		} else {
			dataForUpsert = append(dataForUpsert, ",")
			dataForUpsert = append(dataForUpsert, "'"+value+"'")
		}
	}

	if &retrievePictureinput != nil && !flag {
		if retrievePictureinput.Type == "Customer" {
			if productBrandCustomerId != nil {
				query = query + ` and tmc.customer in(`
				for key, value := range productBrandCustomerId {
					if key == 0 {
						query = query + `'` + value + `'`
					} else {
						query = query + `, '` + value + `'`
					}
				}
				query = query + `)`
			}
		}
		if retrievePictureinput.Type == "Product" {
			query = query + ` and ka.target_item in (select tp.id from team_products tp
				 inner join product p on tp.material_code = p.id
				 where tp.is_active = true and tp.is_delete = false
				 and p.is_active = true and p.is_deleted = false`
			if productBrandCustomerId != nil {
				query = query + ` and p.id in (`
				for key, value := range productBrandCustomerId {
					if key == 0 {
						query = query + `'` + value + `'`
					} else {
						query = query + `, '` + value + `'`
					}
				}
				query = query + `)`
			}
			query = query + `)`
		}
		if retrievePictureinput.Type == "Brand" {
			if productBrandCustomerId != nil {
				query = query + ` and ka.target_item in (`
				for key, value := range productBrandCustomerId {
					if key == 0 {
						query = query + `'` + value + `'`
					} else {
						query = query + `, '` + value + `'`
					}
				}
				query = query + `)`
			}
		}
		if retrievePictureinput.TeamID != "" {
			query = query + ` and tm.team = ?`
			inputArgs = append(inputArgs, retrievePictureinput.TeamID)
		}

		if startDate != "" || endDate != "" {
			query = query + ` and to_char(ka.date_created,'YYYY/MM/DD') between ? and ?`
			inputArgs = append(inputArgs, startDate)
			inputArgs = append(inputArgs, endDate)
		}
	}
	query += ` group by tm.team,ka.answers, ka.target_item, ka.team_member_customer, ka.schedule_event_id, c.value, name, tmc.customer,ka.date_created,ka.category`

	query += `)
	select team, answers, type, name, customer 
	from _events
	where row_number = 1`

	if flag || !onlySalesRep {
		query, inputArgs, err = sqlx.In(query, inputArgs...)
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err = pool.Query(context.Background(), query, inputArgs...)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		return []entity.PictureData{}, err
	}
	pictures := []entity.PictureData{}
	defer rows.Close()
	for rows.Next() {
		picture := entity.PictureData{}
		err := rows.Scan(
			&picture.Team,
			&picture.Url,
			&picture.Type,
			&picture.Name,
			&picture.Customer,
		)
		if err != nil {
			logengine.GetTelemetryClient().TrackException(err.Error())
			return pictures, err
		}
		pictures = append(pictures, picture)
	}
	return pictures, err
}
