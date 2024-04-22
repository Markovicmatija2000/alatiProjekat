package services

import "ProjectModule/model"

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
