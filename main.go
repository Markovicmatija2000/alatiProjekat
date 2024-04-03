package main

import (
	"ProjectModule/repositories"
	"ProjectModule/services"
)

func main() {
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)
	service.Hello()
}
