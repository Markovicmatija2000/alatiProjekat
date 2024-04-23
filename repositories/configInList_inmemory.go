package repositories

import (
	"ProjectModule/model"
	"errors"
	"strings"
)

type ConfigInListInMemRepository struct {
	configInLists map[string]model.ConfigInList
}

func NewConfigInListInMemRepository() model.ConfigInListRepository {
	return &ConfigInListInMemRepository{
		configInLists: make(map[string]model.ConfigInList),
	}
}

func (c *ConfigInListInMemRepository) Add(config model.ConfigInList) {
	c.configInLists[config.Name] = config
}

func (c *ConfigInListInMemRepository) Get(name string) (model.ConfigInList, error) {
	config, ok := c.configInLists[name]
	if !ok {
		return model.ConfigInList{}, errors.New("config not found")
	}
	return config, nil
}

func (c *ConfigInListInMemRepository) Delete(name string) error {
	_, ok := c.configInLists[name]
	if !ok {
		return errors.New("config not found")
	}
	delete(c.configInLists, name)
	return nil
}

func (c *ConfigInListInMemRepository) NewConfigFromLiteral(literal string) (model.ConfigInList, error) {
	parts := strings.Fields(literal)
	if len(parts) < 1 {
		return model.ConfigInList{}, errors.New("invalid literal format")
	}

	params := make(map[string]string)
	for i := 1; i < len(parts); i++ {
		param := strings.Split(parts[i], "=")
		if len(param) != 2 {
			return model.ConfigInList{}, errors.New("invalid parameter format")
		}
		params[param[0]] = param[1]
	}

	return model.ConfigInList{
		Name:   parts[0],
		Params: params,
	}, nil
}
