package model

type ConfigGroup struct {
	Name       string   `json:"name"`
	Version    int      `json:"version"`
	ConfigList []Config `json:"configList"`
}

type ConfigGroupRepository interface {
	AddGroup(configGroup ConfigGroup)
	GetGroup(name string, version int) (ConfigGroup, error)
	DeleteGroup(name string, version int) error
}