package mapper

import (
	"encoding/json"
	"strings"

	"github.com/eztrade/kpi/graph/entity"
	"github.com/eztrade/kpi/graph/model"
)

func MapPicturesEntityToModel(inputEntity entity.PictureData) ([]string, string, string, string, error) {
	var Response []string
	KpiAnswerEntity := []model.KPIAnswerStruct{}
	if inputEntity.Url.Valid {
		err := json.Unmarshal([]byte(inputEntity.Url.String), &KpiAnswerEntity)
		if err != nil {
			return Response, "", "", "", err
		}
	}
	var competitorName string
	var surveyName string
	var promotionName string
	for key, eachDesign := range KpiAnswerEntity {
		if strings.ToLower(inputEntity.Type.String) == "competitor" && key == 0 {
			competitorName = eachDesign.Value[0]
		}
		if strings.ToLower(inputEntity.Type.String) == "survey" && key == 0 {
			surveyName = eachDesign.Value[0]
		}
		if strings.ToLower(inputEntity.Type.String) == "promotion" && key == 0 {
			promotionName = eachDesign.Value[0]
		}
		for _, item := range eachDesign.Value {
			if strings.Contains(strings.ToLower(item), ".png") || strings.Contains(strings.ToLower(item), ".jpg") || strings.Contains(strings.ToLower(item), ".jpeg") {
				if strings.ToLower(inputEntity.Type.String) == "product" {
					return eachDesign.Value, "", "", "", nil
				}
				if strings.ToLower(inputEntity.Type.String) == "competitor" {
					return eachDesign.Value, competitorName, "", "", nil
				}
				if strings.ToLower(inputEntity.Type.String) == "survey" {
					return eachDesign.Value, "", surveyName, "", nil
				}
				if strings.ToLower(inputEntity.Type.String) == "promotion" {
					return eachDesign.Value, "", "", promotionName, nil
				}
				if strings.ToLower(inputEntity.Type.String) == "brand" {
					return eachDesign.Value, "", "", "", nil
				}
			}
		}
	}

	return Response, "", "", "", nil
}
