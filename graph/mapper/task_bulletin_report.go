package mapper

import (
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/eztrade/kpi/graph/azure"
	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/logengine"
	"github.com/eztrade/kpi/graph/model"
	postgres "github.com/eztrade/kpi/graph/postgres"
)

func MapValidTaskBulletinData(input []entity.ReportData, dateRange []string) []entity.ReportData {
	results := []entity.ReportData{}
	for _, value := range input {
		flag, num := DateValueInRange(value, dateRange)
		if flag {
			max := dateRange[num]
			min := dateRange[num-1]
			max1 := max[:10]
			min1 := min[:10]
			dateValue := min1 + "/" + max1
			result := entity.ReportData{}
			result.WeekDateValue = dateValue
			result.BulletinId = value.BulletinId
			result.BulletinTitle = value.BulletinTitle
			result.PrincipalName = value.PrincipalName
			result.CustomerId = value.CustomerId
			result.Attachments = value.Attachments
			result.CustomerName = value.CustomerName
			result.TeamName = value.TeamName
			result.UserName = value.UserName
			result.ActiveDirectory = value.ActiveDirectory
			result.CreationDate = value.CreationDate
			result.TargetDate = value.TargetDate
			result.Status = value.Status
			result.BulletinType = value.BulletinType
			result.Status = value.Status
			result.WeekNumber = num
			result.Remarks = value.Remarks
			result.FeedbackDate = value.FeedbackDate
			results = append(results, result)
		}
	}
	return results
}

func DateValueInRange(input entity.ReportData, dateRange []string) (bool, int) {
	var result bool
	firstDate := dateRange[0]
	lastDate := dateRange[len(dateRange)-1]
	stDate := firstDate[:10]
	endDate := lastDate[:10]
	layoutISO := "2006-01-02"
	loc, _ := time.LoadLocation("Asia/Kolkata")
	t1, _ := time.ParseInLocation(layoutISO, stDate, loc)
	t2, _ := time.ParseInLocation(layoutISO, endDate, loc)
	if strings.EqualFold("Friday", firstDate[11:]) {
		t1 = t1.Add(time.Hour * 18)
	}
	if strings.EqualFold("Friday", lastDate[11:]) {
		t2 = t2.Add(time.Hour * 18)
	} else {
		t2 = t2.Add(time.Hour * 23)
		t2 = t2.Add(time.Minute * 59)
		t2 = t2.Add(time.Second * 59)
	}
	feedBackTime := input.FeedbackDate

	for i := 0; i < len(dateRange)-1; i++ {
		if len(dateRange) == 2 {
			if t1.Before(feedBackTime) && feedBackTime.Before(t2) {
				return true, i + 1
			}
		} else {
			current := dateRange[i]
			currentDate := current[:10]
			next := dateRange[i+1]
			nextDate := next[:10]
			prev, _ := time.ParseInLocation(layoutISO, currentDate, loc)
			upcoming, _ := time.ParseInLocation(layoutISO, nextDate, loc)
			prev = prev.Add(time.Hour * 18)
			upcoming = upcoming.Add(time.Hour * 18)

			if i == 0 {
				if t1.Before(feedBackTime) && feedBackTime.Before(upcoming) {
					return true, i + 1
				}
			} else if i == len(dateRange)-2 {
				if t2.After(feedBackTime) && feedBackTime.After(prev) {
					return true, i + 1
				}
			} else {
				if upcoming.After(feedBackTime) && prev.Before(feedBackTime) {
					return true, i + 1
				}
			}
		}
	}
	return result, 0

}

func MapLatestFeedBack(input []entity.ReportData) map[entity.LatestFeedBack]entity.ReportData {
	slice := make(map[entity.LatestFeedBack]entity.ReportData)
	for _, value := range input {
		var enity entity.LatestFeedBack
		enity.BulletinId = value.BulletinId.String
		enity.WeekNumber = value.WeekNumber
		enity.Customer = value.CustomerId.String
		if v, ok := slice[enity]; !ok {
			data := entity.ReportData{}
			data.BulletinId = value.BulletinId
			data.ActiveDirectory = value.ActiveDirectory
			data.BulletinTitle = value.BulletinTitle
			data.BulletinType = value.BulletinType
			data.FeedbackDate = value.FeedbackDate
			data.PrincipalName = value.PrincipalName
			data.TeamName = value.TeamName
			data.UserName = value.UserName
			data.CustomerId = value.CustomerId
			data.CustomerName = value.CustomerName
			data.CreationDate = value.CreationDate
			data.TargetDate = value.TargetDate
			data.WeekNumber = value.WeekNumber
			data.WeekDateValue = value.WeekDateValue
			data.Status = value.Status
			data.Remarks = value.Remarks
			data.Attachments = value.Attachments
			slice[enity] = data
		} else {
			if value.FeedbackDate.After(v.FeedbackDate) {
				data := entity.ReportData{}
				data.BulletinId = value.BulletinId
				data.ActiveDirectory = value.ActiveDirectory
				data.BulletinTitle = value.BulletinTitle
				data.BulletinType = value.BulletinType
				data.FeedbackDate = value.FeedbackDate
				data.PrincipalName = value.PrincipalName
				data.TeamName = value.TeamName
				data.UserName = value.UserName
				data.CustomerName = value.CustomerName
				data.CustomerId = value.CustomerId
				data.CreationDate = value.CreationDate
				data.TargetDate = value.TargetDate
				data.WeekNumber = value.WeekNumber
				data.WeekDateValue = value.WeekDateValue
				data.Status = value.Status
				data.Remarks = value.Remarks
				data.Attachments = value.Attachments
				slice[enity] = data
			}
		}
	}
	return slice

}

func MapValueOfMapToEntity(input map[entity.LatestFeedBack]entity.ReportData, array []string) []entity.ReportData {
	output := []entity.ReportData{}
	for _, value := range input {
		for i := 1; i <= len(array)-1; i++ {
			key := entity.LatestFeedBack{}
			data := entity.ReportData{}
			key.Customer = value.CustomerId.String
			key.BulletinId = value.BulletinId.String
			key.WeekNumber = i
			if v, ok2 := input[key]; ok2 {
				data.BulletinId = v.BulletinId
				data.ActiveDirectory = v.ActiveDirectory
				data.BulletinTitle = v.BulletinTitle
				data.BulletinType = v.BulletinType
				data.FeedbackDate = v.FeedbackDate
				data.PrincipalName = v.PrincipalName
				data.TeamName = v.TeamName
				data.UserName = v.UserName
				data.CustomerId = v.CustomerId
				data.CustomerName = v.CustomerName
				data.CreationDate = v.CreationDate
				data.TargetDate = v.TargetDate
				data.WeekNumber = v.WeekNumber
				data.WeekDateValue = v.WeekDateValue
				data.Status = v.Status
				data.Remarks = v.Remarks
				data.Attachments = v.Attachments
				output = append(output, data)
			}

		}
	}
	return output
}

func MapFeedBackReportToModel(input []entity.ReportData) []*model.TaskReport {
	dataSet := []entity.UniqueTaskBulletinReport{}
	feedBackSet := []entity.UniqueFeedBack{}
	for _, value := range input {
		feedBack := entity.UniqueFeedBack{}
		data := entity.UniqueTaskBulletinReport{}
		data.BulletinId = value.BulletinId.String
		data.ActiveDirectory = value.ActiveDirectory.String
		data.BulletinTitle = value.BulletinTitle.String
		data.BulletinType = value.BulletinType.String
		data.PrincipalName = value.PrincipalName.String
		data.TeamName = value.TeamName.String
		data.UserName = value.UserName.String
		data.CustomerName = value.CustomerName.String
		data.CustomerId = value.CustomerId.String
		data.CreationDate = value.CreationDate.String
		data.TargetDate = value.TargetDate.String
		feedBack.BulletinId = value.BulletinId.String
		feedBack.WeekNumber = value.WeekNumber
		feedBack.WeekDate = value.WeekDateValue
		feedBack.Customer = value.CustomerId.String
		feedBack.Status = value.Status.String
		feedBack.Remark = value.Remarks.String
		feedBack.Attachments = value.Attachments.String
		dataSet = append(dataSet, data)
		feedBackSet = append(feedBackSet, feedBack)
	}
	response := []*model.TaskReport{}
	uniqueData := UniqueReportDataSet(dataSet)
	uniQueFeedback := UniqueFeedbacks(feedBackSet)
	for _, value := range uniqueData {
		output := model.TaskReport{}
		BulletinTitle := value.BulletinTitle
		BulletinType := value.BulletinType
		PrinicpalName := value.PrincipalName
		TeamName := value.TeamName
		UserName := value.UserName
		ActiveDirectory := value.ActiveDirectory
		CustomerName := value.CustomerName
		CreationDate := value.CreationDate
		TargetDate := value.TargetDate
		output.ActiveDirectory = &ActiveDirectory
		output.PrincipalName = &PrinicpalName
		output.BulletinTitle = &BulletinTitle
		output.BulletinType = &BulletinType
		output.TeamName = &TeamName
		output.UserName = &UserName
		output.CustomerName = &CustomerName
		output.CreationDate = &CreationDate
		output.TargetDate = &TargetDate
		WeeklyData := []*model.WeeklyFeedback{}
		for _, inner := range uniQueFeedback {
			if inner.BulletinId == value.BulletinId && inner.Customer == value.CustomerId {
				weekData := model.WeeklyFeedback{}
				status := inner.Status
				remrks := inner.Remark
				week := inner.WeekNumber
				weekDate := inner.WeekDate
				weekData.WeekDateValue = &weekDate
				weekData.Status = &status
				weekData.Remarks = &remrks
				weekData.WeekNumber = &week
				attachments, _ := postgres.GetAttachmentsByIDs(inner.Attachments)
				weekData.Attachments = attachments
				WeeklyData = append(WeeklyData, &weekData)
			}
		}
		output.WeeklyFeedback = append(output.WeeklyFeedback, WeeklyData...)
		response = append(response, &output)
	}
	return response
}

func UniqueReportDataSet(input []entity.UniqueTaskBulletinReport) []entity.UniqueTaskBulletinReport {
	keys := make(map[entity.UniqueTaskBulletinReport]bool)
	uniqueList := []entity.UniqueTaskBulletinReport{}
	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func UniqueFeedbacks(input []entity.UniqueFeedBack) []entity.UniqueFeedBack {
	keys := make(map[entity.UniqueFeedBack]bool)
	uniqueList := []entity.UniqueFeedBack{}
	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func MapValueOfMapToEntityForReport(input map[entity.LatestFeedBack]entity.ReportData, array []string) []entity.ReportData {
	output := []entity.ReportData{}
	for _, value := range input {
		for i := 1; i <= len(array)-1; i++ {
			key := entity.LatestFeedBack{}
			data := entity.ReportData{}
			key.Customer = value.CustomerId.String
			key.BulletinId = value.BulletinId.String
			key.WeekNumber = i
			if _, ok := input[key]; !ok {
				data.BulletinId = value.BulletinId
				data.ActiveDirectory = value.ActiveDirectory
				data.BulletinTitle = value.BulletinTitle
				data.BulletinType = value.BulletinType
				data.FeedbackDate = value.FeedbackDate
				data.PrincipalName = value.PrincipalName
				data.TeamName = value.TeamName
				data.UserName = value.UserName
				data.CustomerId = value.CustomerId
				data.CustomerName = value.CustomerName
				data.CreationDate = value.CreationDate
				data.TargetDate = value.TargetDate
				data.WeekNumber = i

			} else {
				if v, ok2 := input[key]; ok2 {
					data.BulletinId = v.BulletinId
					data.ActiveDirectory = v.ActiveDirectory
					data.BulletinTitle = v.BulletinTitle
					data.BulletinType = v.BulletinType
					data.FeedbackDate = v.FeedbackDate
					data.PrincipalName = v.PrincipalName
					data.TeamName = v.TeamName
					data.UserName = v.UserName
					data.CustomerId = v.CustomerId
					data.CustomerName = v.CustomerName
					data.CreationDate = v.CreationDate
					data.TargetDate = v.TargetDate
					data.WeekNumber = v.WeekNumber
					data.Status = v.Status
					data.Remarks = v.Remarks
					data.Attachments = v.Attachments
				}
			}
			output = append(output, data)
		}
	}
	return output
}

func GrnaretReortForTaskBulletin(input []*model.TaskReport, array []string) (model.TaskReportOutput, error) {
	response := model.TaskReportOutput{}
	f := excelize.NewFile()
	style1, _ := f.NewStyle(
		`{
			"font":
			{
				"bold":true,
				"family":"Verdana",
				"color":"#800000",
				"size":12
			},
			"alignment":
			{
				"horizontal":"center",
				"vertical":"center"
			},
			"border":
			[
				{"type":"left","color":"000000","style":1},
				{"type":"top","color":"000000","style":1},
				{"type":"right","color":"000000","style":1},
				{"type":"bottom","color":"000000","style":1}
			]
		}`)

	style2, _ := f.NewStyle(
		`{
				"font":
				{
					"bold":true,
					"family":"Verdana",
					"color":"#000000",
					"size":8
				},
				"alignment":
				{
					"horizontal":"center",
					"vertical":"center"
				},
				"border":
				[
					{"type":"left","color":"000000","style":1},
					{"type":"top","color":"000000","style":1},
					{"type":"right","color":"000000","style":1},
					{"type":"bottom","color":"000000","style":1}
				]
			}`)

	style3, _ := f.NewStyle(
		`{
				"font":
				{
					"family":"Verdana",
					"color":"#000000",
					"size":8
				},
				"alignment":
				{
					"horizontal":"center",
					"vertical":"center"
				},
				"border":
				[
					{"type":"left","color":"000000","style":1},
					{"type":"top","color":"000000","style":1},
					{"type":"right","color":"000000","style":1},
					{"type":"bottom","color":"000000","style":1}
				]
			}`)

	style4, _ := f.NewStyle(
		`{
		"font":
		{
			"bold":false,
			"family":"Verdana",
			"color":"#000000",
			"size":8
		},
		"alignment":
		{
			"horizontal":"left",
			"vertical":"center",
			"shrink_to_fit":false,
			"wrap_text":false
		},
		"border":
		[
			{"type":"left","color":"000000","style":1},
			{"type":"top","color":"000000","style":1},
			{"type":"bottom","color":"000000","style":1},
			{"type":"right","color":"000000","style":1}
		]
	}`)

	cat := []string{"Task Bulletin Report"}
	for _, v := range cat {
		f.SetSheetName("Sheet1", v)

		f.MergeCell(v, "E1", "I1")
		f.SetCellStyle(v, "E1", "I1", style1)
		f.SetCellValue(v, "E1", "Task Bulletin Weekly Report")

		f.MergeCell(v, "A6", "I6")
		f.SetCellStyle(v, "A6", "I6", style1)

		f.MergeCell(v, "E2", "I2")
		f.SetCellStyle(v, "E2", "I2", style2)
		f.SetCellValue(v, "E2", "Every week considered upto Friday 6.00 PM")

		f.SetCellStyle(v, "A3", "A3", style4)
		f.SetCellValue(v, "A3", "Report Start Date:")

		f.SetCellStyle(v, "A4", "A4", style4)
		f.SetCellValue(v, "A4", "Report End Date:")

		f.MergeCell(v, "B3", "C3")
		f.SetCellStyle(v, "B3", "C3", style4)
		f.SetCellValue(v, "C3", array[0])

		f.MergeCell(v, "B4", "C4")
		f.SetCellStyle(v, "B4", "C4", style4)
		f.SetCellValue(v, "C4", array[len(array)-1])

		f.SetCellStyle(v, "A7", "A7", style2)
		f.SetCellValue(v, "A7", "Task Bulletin Title")

		f.SetCellStyle(v, "B7", "B7", style2)
		f.SetCellValue(v, "B7", "Task Bulletin Type")

		f.SetCellStyle(v, "C7", "C7", style2)
		f.SetCellValue(v, "C7", "Principal Name")

		f.SetCellStyle(v, "D7", "D7", style2)
		f.SetCellValue(v, "D7", "Customer Name")

		f.SetCellStyle(v, "E7", "E7", style2)
		f.SetCellValue(v, "E7", "Team Name")

		f.SetCellStyle(v, "F7", "F7", style2)
		f.SetCellValue(v, "F7", "Employee Name")

		f.SetCellStyle(v, "G7", "G7", style2)
		f.SetCellValue(v, "G7", "Employee AD Name")

		f.SetCellStyle(v, "H7", "H7", style2)
		f.SetCellValue(v, "H7", "Creation Date")

		f.SetCellStyle(v, "I7", "I7", style2)
		f.SetCellValue(v, "I7", "Target Date")

		st := 10
		for i := 1; i <= len(array)-1; i++ {
			bColumn := getColumnName(st) + strconv.Itoa(6)
			stColumn := getColumnName(st+3) + strconv.Itoa(6)
			str1 := array[i-1]
			str2 := array[i]
			str11 := str1[5:10]
			str11 = str11 + "/"
			str11 = str11 + str2[5:10]
			f.MergeCell(v, bColumn, stColumn)
			f.SetCellStyle(v, bColumn, stColumn, style2)
			f.SetCellValue(v, bColumn, str11)
			st += 4
		}

		bAndORowNum := 7
		l := 2
		i := 1

		for k := 0; k < len(array)-1; k++ {
			bColumn := getColumnName(i+9) + strconv.FormatInt(int64(bAndORowNum), 10)
			oColumn := getColumnName(l+9) + strconv.FormatInt(int64(bAndORowNum), 10)
			oColumn2 := getColumnName(l+11) + strconv.FormatInt(int64(bAndORowNum), 10)

			f.SetCellStyle(v, bColumn, bColumn, style2)
			f.SetCellValue(v, bColumn, "Status")

			f.MergeCell(v, oColumn, oColumn2)
			f.SetCellStyle(v, oColumn, oColumn2, style2)
			f.SetCellValue(v, oColumn, "Remarks")
			l += 4
			i += 4
		}

		for j, value := range input {
			f.SetCellStyle(v, "A"+strconv.Itoa(j+8), "A"+strconv.Itoa(j+8), style3)
			f.SetCellValue(v, "A"+strconv.Itoa(j+8), *value.BulletinTitle)

			f.SetCellStyle(v, "B"+strconv.Itoa(j+8), "B"+strconv.Itoa(j+8), style3)
			f.SetCellValue(v, "B"+strconv.Itoa(j+8), *value.BulletinType)

			f.SetCellStyle(v, "C"+strconv.Itoa(j+8), "C"+strconv.Itoa(j+8), style3)
			f.SetCellValue(v, "C"+strconv.Itoa(j+8), *value.PrincipalName)

			f.SetCellStyle(v, "D"+strconv.Itoa(j+8), "D"+strconv.Itoa(j+8), style3)
			f.SetCellValue(v, "D"+strconv.Itoa(j+8), *value.CustomerName)

			f.SetCellStyle(v, "E"+strconv.Itoa(j+8), "E"+strconv.Itoa(j+8), style3)
			f.SetCellValue(v, "E"+strconv.Itoa(j+8), *value.TeamName)

			f.SetCellStyle(v, "F"+strconv.Itoa(j+8), "F"+strconv.Itoa(j+8), style3)
			f.SetCellValue(v, "F"+strconv.Itoa(j+8), *value.UserName)

			f.SetCellStyle(v, "G"+strconv.Itoa(j+8), "G"+strconv.Itoa(j+8), style3)
			f.SetCellValue(v, "G"+strconv.Itoa(j+8), *value.ActiveDirectory)

			f.SetCellStyle(v, "H"+strconv.Itoa(j+8), "H"+strconv.Itoa(j+8), style3)
			f.SetCellValue(v, "H"+strconv.Itoa(j+8), *value.CreationDate)

			f.SetCellStyle(v, "I"+strconv.Itoa(j+8), "I"+strconv.Itoa(j+8), style3)
			f.SetCellValue(v, "I"+strconv.Itoa(j+8), *value.TargetDate)
			r := 9
			s := 10

			for _, inner := range value.WeeklyFeedback {
				bColumn := getColumnName(*inner.WeekNumber+r) + strconv.Itoa(j+8)
				oColumn := getColumnName(*inner.WeekNumber+s) + strconv.Itoa(j+8)
				oColumn2 := getColumnName(*inner.WeekNumber+s+2) + strconv.Itoa(j+8)

				f.SetCellStyle(v, bColumn, bColumn, style2)
				f.SetCellValue(v, bColumn, *inner.Status)

				f.MergeCell(v, oColumn, oColumn2)
				f.SetCellStyle(v, oColumn, oColumn2, style2)
				f.SetCellValue(v, oColumn, *inner.Remarks)
				r += 3
				s += 3
			}
		}
	}

	fileName := "TaskBulletin_" + time.Now().Format("20060102150405") + ".xlsx"
	// for saving file locally
	// if err := f.SaveAs(fileName); err != nil {
	// 	println(err.Error())
	// }
	// response.URL = ""

	blobURL, err := azure.UploadBytesToBlob(getBytesFromExcel(f), fileName)
	if err != nil {
		return response, err
	}
	response.URL = blobURL
	return response, nil

}

func getColumnName(columnNumber int) string {
	dividend := columnNumber
	columnName := ""
	var modulo int
	for dividend > 0 {
		modulo = (dividend - 1) % 26
		columnName = string(65+modulo) + columnName
		dividend = (int)((dividend - modulo) / 26)
	}
	return columnName
}

func getBytesFromExcel(f *excelize.File) []byte {
	buffer, err := f.WriteToBuffer()
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	myresult := buffer.Bytes()
	return myresult
}
