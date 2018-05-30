package response

import "time"

type Attribute struct {
	Name   string  `json:"name"`
	Format string  `json:"format"`
	Score  float64 `json:"score"`
}

type Attributes []Attribute

type Response struct {
	Pii          bool      `json:"pii"`
	Locale       string    `json:"locale"`
	Score        float64   `json:"score"`
	ModelVersion string    `json:"model_version"`
	RunDate      time.Time `json:"run_date"`
	Attributes   Attributes `json:"attributes"`
}