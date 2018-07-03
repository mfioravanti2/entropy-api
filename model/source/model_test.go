package source

import (
	"testing"
	"fmt"
)

type scoreTest struct {
	Mnemonic string
	Format string
	Score float64
}

func TestModel_Score(t *testing.T) {
	a := Attributes{
		{ Mnemonic: "personal.number", Name: "Personal Number", Id: "9382fc71-3248-4a14-b9fe-2a9ed4b1dce2", Formats: Formats{
			{ Format: "mean", Score: 1.0 },
			{ Format: "naive", Score: 2.0 },
			{ Format: "rare", Score: 4.0 },
		}},
		{ Mnemonic: "work.number", Name: "Work Number", Id: "dfa98b32-e09e-4062-b384-75a8252e7e1f", Formats: Formats{
			{ Format: "mean", Score: 10.0 },
			{ Format: "naive", Score: 20.0 },
			{ Format: "rare", Score: 40.0 },
		}},
	}

	m := Model{ Locale: "US", Threshold: 4.0, Attributes: a }

	valid := []scoreTest{
		{ Mnemonic: "personal.number", Format: "naive", Score: 2.0 },
		{ Mnemonic: "personal.number", Format: "mean", Score: 1.0 },
		{ Mnemonic: "personal.number", Format: "rare", Score: 4.0 },
		{ Mnemonic: "work.number", Format: "naive", Score: 20.0 },
		{ Mnemonic: "work.number", Format: "mean", Score: 10.0 },
		{ Mnemonic: "work.number", Format: "rare", Score: 40.0 },
	}

	for _, item := range valid {
		if h, err := m.Score( item.Mnemonic, item.Format ); h != item.Score || err != nil {
			s := fmt.Sprintf("expected score{mnemonic:'%s', format:'%s')  = %f", item.Mnemonic, item.Format, h )
			t.Error( s )
		}
	}
}
