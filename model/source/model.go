package source

import (
	"errors"
	"fmt"
	"time"
)

type Heuristic struct {
	Match  []string `json:"match"`
	Insert []string `json:"insert"`
	Remove []string `json:"remove"`
}

type Heuristics []Heuristic

type Format struct {
	Format string `json:"format"`
	Score  float64    `json:"score"`
}

type Formats []Format

type Attribute struct {
	Name    string `json:"name"`
	Formats Formats `json:"formats"`
}

type Attributes []Attribute

type Model struct {
	Locale     string  `json:"locale"`
	Threshold  float64 `json:"threshold"`
	K          int     `json:"k"`
	ModelVersion string `json:"version"`
	ModelDate time.Time `json:"timestamp"`
	Heuristics Heuristics `json:"heuristics"`
	Attributes Attributes `json:"attributes"`
}

type Threshold struct {
	Locale     string  `json:"locale"`
	Threshold  float64 `json:"threshold"`
	K          int     `json:"k"`
}

func GetScore( m Model, n string, t string ) (float64,error) {
	var a Attribute
	var f Format

	for _, a = range m.Attributes {
		if a.Name == n {
			for _, f = range a.Formats {
				if f.Format == t {
					return f.Score, nil
				}
			}
		}
	}

	s := fmt.Sprintf("attribute (%s/%s) not found in country (%s)", n, t, m.Locale)
	return 0.0, errors.New(s)
}

