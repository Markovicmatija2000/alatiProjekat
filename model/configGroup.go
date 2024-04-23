package model

type ConfigGroup struct {
	Name         string         `json:"name"`
	Version      int            `json:"version"`
	ConfigInList []ConfigInList `json:"configInList"`
}

type ConfigGroupRepository interface {
	AddGroup(configGroup ConfigGroup)
	GetGroup(name string, version int) (ConfigGroup, error)
	DeleteGroup(name string, version int) error

	AddConfigToGroup(group ConfigGroup, config ConfigInList) error
	RemoveConfigFromGroup(group ConfigGroup, key int) error

	//ParseConfigData(data []string) (ConfigGroup, error)
}
