package services

import (
	"ProjectModule/model"
)

type ConfigInListService struct {
	repo model.ConfigInListRepository
}

func NewConfigInListService(repo model.ConfigInListRepository) ConfigInListService {
	return ConfigInListService{
		repo: repo,
	}
}

func (s ConfigInListService) Add(config model.ConfigInList) {
	s.repo.Add(config)
}

func (s ConfigInListService) Get(name string) (model.ConfigInList, error) {
	return s.repo.Get(name)
}

func (s ConfigInListService) Delete(name string) error {
	return s.repo.Delete(name)
}
