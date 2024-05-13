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

func (c *ConfigGroupInMemRepository) AddConfigToGroup(group model.ConfigGroup, config model.ConfigInList) error {
	group.ConfigInList = append(group.ConfigInList, config)
	return nil
}

func (c *ConfigGroupInMemRepository) RemoveConfigFromGroup(group model.ConfigGroup, index int) error {
	if index < 0 || index >= len(group.ConfigInList) {
		return errors.New("index out of range")
	}
	group.ConfigInList = append(group.ConfigInList[:index], group.ConfigInList[index+1:]...)
	return nil
}

/*func (c *ConfigInMemRepository) NewConfigGroupFromLiteral(literal string) (model.ConfigGroup, error) {
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

func (r *ConfigGroupInMemRepository) ParseConfigData(data []string) (model.ConfigGroup, error) {
	var configGroup model.ConfigGroup

	configGroup.Name = data[0]
	version, err := strconv.Atoi(data[1])
	if err != nil {
		return model.ConfigGroup{}, errors.New("invalid version format")
	}
	configGroup.Version = version

	for i := 2; i < len(data); i++ {
		configInList := strings.Split(data[i], "=")
		if len(configInList) != 2 {
			return model.ConfigGroup{}, errors.New("invalid parameter format")
		}

		config := model.ConfigInList{
			Name: configInList[0],
			Params: map[string]string{
				"value": configInList[1], // Assuming your parameter is named "value"
			},
		}
		configGroup.ConfigInList = append(configGroup.ConfigInList, config)
	}

	return configGroup, nil
}*/
