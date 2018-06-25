package response

import "time"

type Attribute struct {
	Mnemonic string  `json:"mnemonic"`
	Tag 	 string	 `json:"tag"`
	Locale   string  `json:"locale,omitempty"`
	Format   string  `json:"format"`
	Score    float64 `json:"score"`
}

type Attributes []Attribute

type Heuristics []string

type Person struct {
	Id		    string      `json:"id"`
	Nationality string		`json:"nationality"`
	Score 		float64		`json:"score"`
	Attributes  Attributes  `json:"attributes,omitempty"`
	Heuristics	*Heuristics  `json:"heuristics,omitempty"`
}

type People []Person

type Data struct {
	Pii          bool       `json:"pii"`
	Locale       string     `json:"locale"`
	Score        float64    `json:"score"`
	ApiVersion   string     `json:"api_version"`
	RunDate      time.Time  `json:"run_date"`
	People   	 People 	`json:"people,omitempty"`
}

type Errors struct {
	Messages []string `json:"messages,omitempty"`
}

type Response struct {
	Data   *Data   `json:"data,omitempty"`
	Errors *Errors `json:"errors,omitempty"`
}