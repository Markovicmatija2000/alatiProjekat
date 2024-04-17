package main

import (
	"ProjectModule/model"
	"ProjectModule/repositories"
	"ProjectModule/services"
	"encoding/json"
	"fmt"
)

func main() {
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)
	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5432"
	config := model.Config{
		Name:    "db_config",
		Version: 2,
		Params:  params,
	}
	config2, error := repo.NewConfigFromLiteral("Ime 23 params=paramsss params2=params22")
	service.Add(config2)
	service.Add(config)
	fmt.Print(error)
	fmt.Print(error)
	

	fmt.Println("marko")


	retrievedConfig1, err := service.Get("Ime", 23)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

	fmt.Println("matijaaa")


	retrievedConfig, err := service.Get("db_config", 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }




	jsonData1, err := json.Marshal(retrievedConfig1)
    if err != nil {
        fmt.Println("Error marshaling JSON:", err)
        return
    }


    jsonData, err := json.Marshal(retrievedConfig)
    if err != nil {
        fmt.Println("Error marshaling JSON:", err)
        return
    }

    fmt.Println(string(jsonData))

	
    fmt.Println(string(jsonData1))
}
