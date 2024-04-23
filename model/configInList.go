package model

type ConfigInList struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
	//	Labels map[string]string `json:"labels"`
}

type ConfigInListRepository interface {
	Add(configInList ConfigInList)
	Get(name string) (ConfigInList, error)
	Delete(name string) error
	NewConfigFromLiteral(literal string) (ConfigInList, error)
}
