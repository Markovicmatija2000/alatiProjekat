package repositories

import (
	"ProjectModule/model"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

type ConfigGroupConsulRepository struct {
	client *api.Client
}

func NewConfigGroupConsulRepository() (model.ConfigGroupRepository, error) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", os.Getenv("DB"), os.Getenv("DBPORT"))

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigGroupConsulRepository{client: client}, nil
}

func (r *ConfigGroupConsulRepository) AddGroup(configGroup model.ConfigGroup) error {
    kv := r.client.KV()

    // Marshal the configGroup into JSON
    p, err := json.Marshal(configGroup)
    if err != nil {
        // Handle the error, such as logging or other error handling mechanism
        return err
    }

    // Create a KVPair
    kvp := &api.KVPair{
        Key:   fmt.Sprintf("configgroup/%s/%d", configGroup.Name, configGroup.Version),
        Value: p,
    }

    // Put the KVPair into Consul KV store
    _, err = kv.Put(kvp, nil)
    if err != nil {
        // Handle the error, such as logging or other error handling mechanism
        return err
    }

    return nil
}




func (r *ConfigGroupConsulRepository) GetGroup(name string, version int) (model.ConfigGroup, error) {
	kv := r.client.KV()
	key := fmt.Sprintf("configgroup/%s/%d", name, version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	if pair == nil {
		return model.ConfigGroup{}, fmt.Errorf("config group not found")
	}

	var configGroup model.ConfigGroup
	err = json.Unmarshal(pair.Value, &configGroup)
	return configGroup, err
}

func (r *ConfigGroupConsulRepository) DeleteGroup(name string, version int) error {
	kv := r.client.KV()
	key := fmt.Sprintf("configgroup/%s/%d", name, version)
	_, err := kv.Delete(key, nil)
	return err
}

func (r *ConfigGroupConsulRepository) GetConfigInListByLabels(name string, version int, labels []model.ConfigInList) ([]model.ConfigInList, error) {

	configGroup, err := r.GetGroup(name, version)
	if err != nil {
		return nil, err
	}

	var matchingConfigInList []model.ConfigInList


LabelLoop:
	for _, configInList := range configGroup.ConfigInList {

		if len(configInList.Labels) != len(labels) {
			continue
		}


		for _, providedLabel := range labels {
			found := false
			for labelName, labelValue := range configInList.Labels {
				if providedLabel.Name == labelName && providedLabel.Params["value"] == labelValue {
					found = true
					break
				}
			}
			if !found {
				continue LabelLoop
			}
		}


		matchingConfigInList = append(matchingConfigInList, configInList)
	}

	if len(matchingConfigInList) == 0 {
		return nil, fmt.Errorf("configInList with specified labels not found for %s version %d", name, version)
	}

	return matchingConfigInList, nil
}


func (r *ConfigGroupConsulRepository) DeleteConfigInListByLabels(name string, version int, labels []model.ConfigInList) error {
	// Get the configGroup from the repository
	configGroup, err := r.GetGroup(name, version)
	if err != nil {
		return err
	}

	// Iterate over the configInList and remove those that match the provided labels
	var updatedConfigInList []model.ConfigInList
LabelLoop:
	for _, configInList := range configGroup.ConfigInList {

		// Check if the number of labels in the configInList matches the number of provided labels
		if len(configInList.Labels) != len(labels) {
			updatedConfigInList = append(updatedConfigInList, configInList)
			continue
		}

		// Iterate over the provided labels and compare them with the labels in the configInList
		for _, providedLabel := range labels {
			found := false
			for labelName, labelValue := range configInList.Labels {
				if providedLabel.Name == labelName && providedLabel.Params["value"] == labelValue {
					found = true
					break
				}
			}
			if !found {
				updatedConfigInList = append(updatedConfigInList, configInList)
				continue LabelLoop
			}
		}
	}

	// Update the configGroup in the Consul repository
	configGroup.ConfigInList = updatedConfigInList
	err = r.AddGroup(configGroup)
	if err != nil {
		return err
	}

	return nil
}

