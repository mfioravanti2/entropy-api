package response

import "time"

type Attribute struct {
	Name   string  `json:"name"`
	Locale string    `json:"locale,omitempty"`
	Format string  `json:"format"`
	Score  float64 `json:"score"`
}

type Attributes []Attribute

type Data struct {
	Pii          bool      `json:"pii"`
	Locale       string    `json:"locale"`
	Score        float64   `json:"score"`
	ApiVersion   string    `json:"model_version"`
	RunDate      time.Time `json:"run_date"`
	Attributes   Attributes `json:"attributes"`
}

type Errors struct {
	Messages []string `json:"messages,omitempty"`
}

type Response struct {
	Data   *Data `json:"data,omitempty"`
	Errors *Errors `json:"errors,omitempty"`
}