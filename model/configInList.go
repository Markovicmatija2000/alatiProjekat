package model

type ConfigInList struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
//	Labels map[string]string `json:"labels"`
}
