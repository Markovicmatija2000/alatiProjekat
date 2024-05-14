package services

import (
	"ProjectModule/model"
	"errors"
)

type ConfigGroupService struct {
	repo model.ConfigGroupRepository
}

func NewConfigGroupService(repo model.ConfigGroupRepository) ConfigGroupService {
	return ConfigGroupService{
		repo: repo,
	}
}

func (s ConfigGroupService) AddGroup(configGroup model.ConfigGroup) {
	s.repo.AddGroup(configGroup)
}

func (s ConfigGroupService) GetGroup(name string, version int) (model.ConfigGroup, error) {
	return s.repo.GetGroup(name, version)
}

func (s ConfigGroupService) DeleteGroup(name string, version int) error {
	return s.repo.DeleteGroup(name, version)
}

func (s ConfigGroupService) UpdateGroup(configGroup model.ConfigGroup) error {
	// Check if the config group exists
	_, err := s.repo.GetGroup(configGroup.Name, configGroup.Version)
	if err != nil {
		return errors.New("config group not found")
	}

	// Call the AddGroup method of the repository to update the config group
	s.repo.AddGroup(configGroup)

	return nil
}

func (s ConfigGroupService) DeleteConfigInListByLabels(name string, version int, labels []model.ConfigInList) error {
	return s.repo.DeleteConfigInListByLabels(name, version, labels)
}

func (s ConfigGroupService) GetConfigInListByLabels(name string, version int, labels []model.ConfigInList) ([]model.ConfigInList, error) {
	return s.repo.GetConfigInListByLabels(name, version, labels)
}
