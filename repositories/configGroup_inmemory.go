package repositories

import (
	"ProjectModule/model"
	"errors"
	"fmt"
	"reflect"
)

type ConfigGroupInMemRepository struct {
	configGroups map[string]model.ConfigGroup
}

func NewConfigGroupInMemRepository() *ConfigGroupInMemRepository {
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

func (c *ConfigGroupInMemRepository) DeleteConfigInListByLabels(name string, version int, labels []model.ConfigInList) error {
	// Get the configGroup from the repository
	configGroup, err := c.GetGroup(name, version)
	if err != nil {
		return err
	}

	// Iterate over the configInList and remove those that match the provided labels
	var updatedConfigInList []model.ConfigInList
	for _, configInList := range configGroup.ConfigInList {
		if !reflect.DeepEqual(configInList.Labels, labels) {
			updatedConfigInList = append(updatedConfigInList, configInList)
		}
	}

	// Update the configGroup in the repository
	configGroup.ConfigInList = updatedConfigInList
	c.AddGroup(configGroup)

	return nil
}

func (c *ConfigGroupInMemRepository) GetConfigInListByLabels(name string, version int, labels []model.ConfigInList) ([]model.ConfigInList, error) {
	// Get the configGroup from the repository
	configGroup, err := c.GetGroup(name, version)
	if err != nil {
		return nil, err
	}

	// Print labels for debugging
	fmt.Println("Provided labels:", labels)

	// Initialize a slice to store matching configInList
	var matchingConfigInList []model.ConfigInList

	// Iterate over the configInList and add those that match the provided labels
LabelLoop:
	for _, configInList := range configGroup.ConfigInList {

		fmt.Println("ConfigInList labels:", configInList.Labels)

		// Check if the number of labels in the configInList matches the number of provided labels
		if len(configInList.Labels) != len(labels) {
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
				continue LabelLoop
			}
		}

		// If all labels match, add the configInList to matchingConfigInList
		matchingConfigInList = append(matchingConfigInList, configInList)
	}

	if len(matchingConfigInList) == 0 {
		return nil, fmt.Errorf("configInList with specified labels not found for %s version %d", name, version)
	}

	return matchingConfigInList, nil
}
