package request

import (
	"testing"
	"fmt"
)

func TestRequest_Validate(t *testing.T) {
	valid := []Request{
		{Locale: "us", People: People{
			{Nationality: "us", PersonID: "0", Attributes: Attributes{
				{Mnemonic: "phone.nanpa.full", Format: "rare", Tag: "work"},
				{Mnemonic: "ssn", Format: "mean", Tag: ""},
			}},
		}},
		{Locale: "us", People: People{
			{Nationality: "us", PersonID: "0", Attributes: Attributes{
				{Mnemonic: "phone.nanpa.full", Format: "rare", Tag: "work"},
				{Mnemonic: "ssn", Format: "mean", Tag: ""},
			}},
			{Nationality: "us", PersonID: "1", Attributes: Attributes{
				{Mnemonic: "phone.nanpa.full", Format: "rare", Tag: "work"},
				{Mnemonic: "ssn", Format: "mean", Tag: ""},
			}},
		}},
	}

	for _, r := range valid {
		if ok, err := r.Validate(); !ok || err != nil {
			s := fmt.Sprintf("expected request{l:'%s')  = ok", r.Locale )
			t.Error( s )
		}
	}

	invalid := []Request{
		{Locale: "usa", People: People{
			{Nationality: "us", PersonID: "0", Attributes: Attributes{
				{Mnemonic: "phone.nanpa.full", Format: "rare", Tag: "work"},
				{Mnemonic: "ssn", Format: "mean", Tag: ""},
			}},
		}},
		{Locale: "us", People: People{
			{Nationality: "usa", PersonID: "0", Attributes: Attributes{
				{Mnemonic: "phone.nanpa.full", Format: "rare", Tag: "work"},
				{Mnemonic: "ssn", Format: "mean", Tag: ""},
			}},
			{Nationality: "us", PersonID: "1", Attributes: Attributes{
				{Mnemonic: "phone.nanpa.full", Format: "rare", Tag: "work"},
				{Mnemonic: "ssn", Format: "mean", Tag: "..."},
			}},
		}},
	}

	for _, r := range invalid {
		if ok, err := r.Validate(); ok || err == nil {
			s := fmt.Sprintf("expected request{l:'%s') != ok", r.Locale )
			t.Error( s )
		}
	}
}

func TestPerson_Validate(t *testing.T) {
	valid := People{
		{Nationality: "us", PersonID: "0", Attributes: Attributes{
			{Mnemonic: "phone.nanpa.full", Format: "rare", Tag: "work"},
			{Mnemonic: "ssn", Format: "mean", Tag: ""},
		}},
		{Nationality: "us", PersonID: "1", Attributes: Attributes{
			{Mnemonic: "phone.nanpa.full", Format: "rare", Tag: "work"},
			{Mnemonic: "ssn", Format: "mean", Tag: ""},
		}},
	}

	for _, p := range valid {
		if ok, err := p.Validate(); !ok || err != nil {
			s := fmt.Sprintf("expected person{n:'%s', i:'%s')  = ok", p.Nationality, p.PersonID )
			t.Error( s )
		}
	}

	invalid := People{
		{Nationality: "usa", PersonID: "0", Attributes: Attributes{
			{Mnemonic: "phone.nanpa.full", Format: "rare", Tag: "work"},
			{Mnemonic: "ssn", Format: "mean", Tag: ""},
		}},
		{Nationality: "us", PersonID: "1", Attributes: Attributes{
			{Mnemonic: "phone.nanpa.full", Format: "none", Tag: "work"},
			{Mnemonic: "ssn", Format: "mean", Tag: ""},
		}},
		{Nationality: "us", PersonID: "2", Attributes: Attributes{
			{Mnemonic: "phone.nanpa.full", Format: "mean", Tag: "work"},
			{Mnemonic: "ssn", Format: "mean", Tag: "..."},
		}},
	}

	for _, p := range invalid {
		if ok, err := p.Validate(); ok || err == nil {
			s := fmt.Sprintf("expected person{n:'%s', i:'%s') != ok", p.Nationality, p.PersonID )
			t.Error( s )
		}
	}

}

func TestAttribute_Validate(t *testing.T) {
	valid := Attributes{
		{Mnemonic: "phone.nanpa.full", Format: "rare", Tag: "work"},
		{Mnemonic: "ssn", Format: "mean", Tag: ""},
	}

	for _, a := range valid {
		if ok, err := a.Validate(); !ok || err != nil {
			s := fmt.Sprintf("expected attribute{m:'%s', f:'%s', t:'%s')  = ok", a.Mnemonic, a.Format, a.Tag )
			t.Error( s )
		}
	}

	invalid := Attributes{
		{Mnemonic: "phone.nanpa.full", Format: "none", Tag: "work"},
		{Mnemonic: "ssn", Format: "mean", Tag: "..."},
		{Mnemonic: "ssn++", Format: "mean", Tag: "..."},
	}

	for _, a := range invalid {
		if ok, err := a.Validate(); ok || err == nil {
			s := fmt.Sprintf("expected attribute{m:'%s', f:'%s', t:'%s') != ok", a.Mnemonic, a.Format, a.Tag )
			t.Error( s )
		}
	}
}
