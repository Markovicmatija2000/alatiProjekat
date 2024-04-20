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
	repo2 := repositories.NewConfigGroupInMemRepository()
	service := services.NewConfigService(repo)
	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5431"
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

	configData := []string{
		"GroupA",        // Name of the group
		"1",             // Version of the group
		"param1=value1", // Configuration parameter 1
		"param2=value2", // Configuration parameter 2
		// Add more configuration parameters as needed
	}

	configGroup, err := repo2.ParseConfigData(configData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Group Name:", configGroup.Name)
	fmt.Println("Group Version:", configGroup.Version)
	fmt.Println("Configuration Parameters:")
	for _, config := range configGroup.ConfigInList {
		fmt.Printf("- Name: %s, Value: %s\n", config.Name, config.Params["value"])
	}

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
