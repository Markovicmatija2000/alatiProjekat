package repositories

import (
	"ProjectModule/model"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
)

type ConfigConsulRepository struct {
	client *api.Client
}

func NewConfigConsulRepository() (model.ConfigRepository, error) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", os.Getenv("DB"), os.Getenv("DBPORT"))

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigConsulRepository{client: client}, nil
}

func (r *ConfigConsulRepository) Add(config model.Config) error {
	kv := r.client.KV()
	p, err := json.Marshal(config)
	if err != nil {
		return err
	}

	kvp := &api.KVPair{
		Key:   fmt.Sprintf("config/%s/%d", config.Name, config.Version),
		Value: p,
	}

	_, err = kv.Put(kvp, nil)
	return err
}

func (r *ConfigConsulRepository) Get(name string, version int) (model.Config, error) {
	kv := r.client.KV()
	key := fmt.Sprintf("config/%s/%d", name, version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return model.Config{}, err
	}
	if pair == nil {
		return model.Config{}, fmt.Errorf("config not found")
	}

	var config model.Config
	err = json.Unmarshal(pair.Value, &config)
	return config, err
}

func (r *ConfigConsulRepository) Delete(name string, version int) error {
	kv := r.client.KV()
	key := fmt.Sprintf("config/%s/%d", name, version)
	_, err := kv.Delete(key, nil)
	return err
}

func (r *ConfigConsulRepository) NewConfigFromLiteral(literal string) (model.Config, error) {
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
