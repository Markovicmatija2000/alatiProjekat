package repositories

import (
	"ProjectModule/model"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/consul/api"
)

type ConfigInListConsulRepository struct {
	client *api.Client
}

func NewConfigInListConsulRepository() (model.ConfigInListRepository, error) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", os.Getenv("DB"), os.Getenv("DBPORT"))

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigInListConsulRepository{client: client}, nil
}

func (r *ConfigInListConsulRepository) Add(config model.ConfigInList) {
	kv := r.client.KV()
	p, err := json.Marshal(config)
	if err != nil {
		fmt.Println("Error marshalling config:", err)
		return
	}

	kvp := &api.KVPair{
		Key:   fmt.Sprintf("configinlist/%s", config.Name),
		Value: p,
	}

	_, err = kv.Put(kvp, nil)
	if err != nil {
		fmt.Println("Error putting config in Consul:", err)
	}
}

func (r *ConfigInListConsulRepository) Get(name string) (model.ConfigInList, error) {
	kv := r.client.KV()
	key := fmt.Sprintf("configinlist/%s", name)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return model.ConfigInList{}, err
	}
	if pair == nil {
		return model.ConfigInList{}, fmt.Errorf("config in list not found")
	}

	var config model.ConfigInList
	err = json.Unmarshal(pair.Value, &config)
	return config, err
}

func (r *ConfigInListConsulRepository) Delete(name string) error {
	kv := r.client.KV()
	key := fmt.Sprintf("configinlist/%s", name)
	_, err := kv.Delete(key, nil)
	return err
}

func (r *ConfigInListConsulRepository) NewConfigFromLiteral(literal string) (model.ConfigInList, error) {
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
		Labels: make(map[string]string),
	}, nil
}
