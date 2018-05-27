package model

type Country struct {
	Name string `json:"name"`
	File string `json:"file"`
}

type Countries []Country

