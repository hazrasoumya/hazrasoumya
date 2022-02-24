package model

type KpiTypeAndCategory struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Value       string `json:"value"`
	Description string `json:"description"`
}
