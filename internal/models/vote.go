package models

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int32  `json:"voice"`
}
