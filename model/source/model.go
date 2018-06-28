package source

import (
	"errors"
	"fmt"
	"time"
	"sort"
)

type Heuristic struct {
	Id     string   `json:"id"`
	Notes  string   `json:"notes"`
	Match  []string `json:"match"`
	Insert []string `json:"insert"`
	Remove []string `json:"remove"`
}

type Heuristics []Heuristic

func ( h Heuristics ) ToStrings() ([]string, error) {
	var ids []string

	for _, heuristic := range h {
		ids = append(ids, heuristic.Id)
	}

	sort.Strings(ids)

	return ids, nil
}

type Format struct {
	Format string     `json:"format"`
	Score  float64    `json:"score"`
}

type Formats []Format

type Source struct {
	Title string     `json:"title"`
	Org   string     `json:"organization"`
	Date  time.Time  `json:"date"`
	URI   string     `json:"url"`
}

type Sources []Source

type Attribute struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Mnemonic string  `json:"mnemonic"`
	Notes    string  `json:"notes"`
	Sources  Sources `json:"sources"`
	Formats  Formats `json:"formats"`
}

type Attributes []Attribute

func ( a Attributes ) ToStrings() ([]string, error) {
	var mnemonics []string

	for _, attribute := range a {
		mnemonics = append(mnemonics, attribute.Mnemonic)
	}

	sort.Strings(mnemonics)

	return mnemonics, nil
}

type Model struct {
	Locale       string     `json:"locale"`
	Threshold    float64    `json:"threshold"`
	K            int        `json:"k"`
	ModelVersion string     `json:"version"`
	ModelDate    time.Time  `json:"timestamp"`
	Heuristics   Heuristics `json:"heuristics"`
	Attributes   Attributes `json:"attributes"`
}

type Threshold struct {
	Locale     string  `json:"locale"`
	Threshold  float64 `json:"threshold"`
	K          int     `json:"k"`
}

func ( m *Model ) Score( n string, t string ) (float64,error) {
	var a Attribute
	var f Format

	for _, a = range m.Attributes {
		if a.Mnemonic == n {
			for _, f = range a.Formats {
				if f.Format == t {
					return f.Score, nil
				}
			}

			s := fmt.Sprintf("attribute (%s) with format (%s) not found in country (%s)", n, t, m.Locale)
			return 0.0, errors.New(s)
		}
	}

	s := fmt.Sprintf("attribute (%s) not found in country (%s)", n, m.Locale)
	return 0.0, errors.New(s)
}

