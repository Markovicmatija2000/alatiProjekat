package repositories

import (
	"ProjectModule/model"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func (c *ConfigInMemRepository) NewConfigGroupFromLiteral(literal string) (model.ConfigGroup, error) {
	parts := strings.Fields(literal)
	if len(parts) < 3 {
		return model.ConfigGroup{}, errors.New("invalid literal format")
	}

	version, err := strconv.Atoi(parts[1])
	if err != nil {
		return model.ConfigGroup{}, errors.New("invalid version format")
	}

	configInLists := make(map[string]string)
	for i := 2; i < len(parts); i++ {
		configInList := strings.Split(parts[i], "=")
		if len(configInList) != 2 {
			return model.ConfigGroup{}, errors.New("invalid parameter format")
		}
		configInLists[configInList[0]] = configInList[1]
	}

	return model.ConfigGroup{
		Name:         parts[0],
		Version:      version,
		ConfigInList: configInLists,
	}, nil
}
