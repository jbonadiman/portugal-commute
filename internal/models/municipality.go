package models

type Municipality struct {
	Name          string   `json:"nome"`
	Neighborhoods []string `json:"freguesias"`
}
