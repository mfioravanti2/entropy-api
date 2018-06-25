package request

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

