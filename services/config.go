package services

import (
	"ProjectModule/model"
	"fmt"
)

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(repo model.ConfigRepository) ConfigService {
	return ConfigService{
		repo: repo,
	}
}

func (s ConfigService) Hello() {
	fmt.Println("Pozdrav :)")
}

// todo: implementiraj metode za dodavanje, brisanje, dobavljanje itd.
