package repositories

import (
	"ProjectModule/model"
	"errors"
	"fmt"
)

type ConfigInMemRepository struct {
	configs map[string]model.Config
}

func NewConfigInMemRepository() model.ConfigRepository {
	return ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}

func (c ConfigInMemRepository) Add(config model.Config) {
	key := fmt.Sprintf("%s/%d", config.Name, config.Version)
	c.configs[key] = config
}

func (c ConfigInMemRepository) Get(name string, version int) (model.Config, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.Config{}, errors.New("config not found")
	}
	return config, nil
}
