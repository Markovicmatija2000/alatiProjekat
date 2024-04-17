package repositories

import (
	"ProjectModule/model"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ConfigInMemRepository struct {
	configs map[string]model.Config
}

func NewConfigInMemRepository() model.ConfigRepository {
	return &ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}

func (c *ConfigInMemRepository) Add(config model.Config) {
	key := fmt.Sprintf("%s/%d", config.Name, config.Version)
	c.configs[key] = config
}

func (c *ConfigInMemRepository) Get(name string, version int) (model.Config, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.Config{}, errors.New("config not found")
	}
	return config, nil
}

func (c *ConfigInMemRepository) Delete(name string, version int) error {
	key := fmt.Sprintf("%s/%d", name, version)
	_, ok := c.configs[key]
	if !ok {
		return errors.New("config not found")
	}
	delete(c.configs, key)
	return nil
}

func (c *ConfigInMemRepository) NewConfigFromLiteral(literal string) (model.Config, error) {
	parts := strings.Fields(literal)
	if len(parts) < 3 {
		return model.Config{}, errors.New("invalid literal format")
	}

	version, err := strconv.Atoi(parts[1])
	if err != nil {
		return model.Config{}, errors.New("invalid version format")
	}

	params := make(map[string]string)
	for i := 2; i < len(parts); i++ {
		param := strings.Split(parts[i], "=")
		if len(param) != 2 {
			return model.Config{}, errors.New("invalid parameter format")
		}
		params[param[0]] = param[1]
	}

	return model.Config{
		Name:    parts[0],
		Version: version,
		Params:  params,
	}, nil
}