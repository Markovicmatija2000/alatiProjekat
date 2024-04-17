package repositories

import (
	"ProjectModule/model"
	"errors"
	"fmt"
)

type ConfigGroupInMemRepository struct {
	configGroups map[string]model.ConfigGroup
}

func NewConfigGroupInMemRepository() model.ConfigGroupRepository {
	return &ConfigGroupInMemRepository{
		configGroups: make(map[string]model.ConfigGroup),
	}
}

func (c *ConfigGroupInMemRepository) AddGroup(configGroup model.ConfigGroup) {
	key := fmt.Sprintf("%s/%d", configGroup.Name, configGroup.Version)
	c.configGroups[key] = configGroup
}

func (c *ConfigGroupInMemRepository) GetGroup(name string, version int) (model.ConfigGroup, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	configGroup, ok := c.configGroups[key]
	if !ok {
		return model.ConfigGroup{}, errors.New("config group not found")
	}
	return configGroup, nil
}

func (c *ConfigGroupInMemRepository) DeleteGroup(name string, version int) error {
	key := fmt.Sprintf("%s/%d", name, version)
	_, ok := c.configGroups[key]
	if !ok {
		return errors.New("config group not found")
	}
	delete(c.configGroups, key)
	return nil
}
