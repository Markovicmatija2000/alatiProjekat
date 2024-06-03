package model

type ConfigGroup struct {
	Name         string         `json:"name"`
	Version      int            `json:"version"`
	ConfigInList []ConfigInList `json:"configInList"`
}

type ConfigGroupRepository interface {
	AddGroup(configGroup ConfigGroup) error
	GetGroup(name string, version int) (ConfigGroup, error)
	DeleteGroup(name string, version int) error

	GetConfigInListByLabels(name string, version int, labels []ConfigInList) ([]ConfigInList, error)
	DeleteConfigInListByLabels(name string, version int, labels []ConfigInList) error

	//ParseConfigData(data []string) (ConfigGroup, error)
}
