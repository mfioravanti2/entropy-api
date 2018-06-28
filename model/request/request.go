package request

import "github.com/mfioravanti2/entropy-api/model"

type Attribute struct {
	Mnemonic string `json:"mnemonic"`
	Format   string `json:"format"`
	Tag 	 string	`json:"tag"`
}

type Attributes []Attribute

type Person struct {
	Nationality string `json:"nationality"`
	PersonID    string `json:"person_id"`
	Attributes  Attributes `json:"attributes"`
}

type People []Person

type Request struct {
	Locale string `json:"locale"`
	People People `json:"people"`
}

func (a *Attribute) Validate() (bool, error) {
	if ok, err := model.ValidateAttributeMnemonic( a.Mnemonic ); !ok {
		return false, err
	}

	if ok, err := model.ValidateFormat( a.Format ); !ok {
		return false, err
	}

	if ok, err := model.ValidateTag( a.Tag ); !ok {
		return false, err
	}

	return true, nil
}

func (p *Person) Validate() (bool, error) {
	if ok, err := model.ValidateCountryCode( p.Nationality ); !ok {
		return false, err
	}

	for _, a := range p.Attributes {
		if ok, err := a.Validate(); !ok {
			return false, err
		}
	}

	return true, nil
}

func (r *Request) Validate() (bool, error) {
	if ok, err := model.ValidateCountryCode( r.Locale ); !ok {
		return false, err
	}

	for _, p := range r.People {
		if ok, err := p.Validate(); !ok {
			return false, err
		}
	}

	return true, nil
}

