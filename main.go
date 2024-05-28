package main

import (
	"ProjectModule/handlers"
	"ProjectModule/model"
	"ProjectModule/repositories"
	"ProjectModule/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

var shutdownTimer = 5 * time.Second

func main() {
	// Definišemo shutdown signal i SIGINT i SIGTERM signale
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	repoS, err := repositories.NewConfigConsulRepository()
	if err != nil {
    	fmt.Println("Error initializing Consul repository:", err)
    	os.Exit(1)
	}

	repo2, err := repositories.NewConfigGroupConsulRepository()
	if err != nil {
		fmt.Println("Error initializing ConfigGroup Consul repository:", err)
		os.Exit(1)
	}
	service := services.NewConfigService(repoS)
	groupService := services.NewConfigGroupService(repo2)
	simulateOperations(service, repoS)

	// Definišemo router i handlere
	configHandler := handlers.NewConfigHandler(service)
	configGroupHandler := handlers.NewConfigGroupHandler(groupService)

	router := mux.NewRouter()
	limiter := rate.NewLimiter(0.5, 1)
	router.HandleFunc("/configs/{name}/{version}", func(w http.ResponseWriter, r *http.Request) {
		handlers.RateLimit(limiter, http.HandlerFunc(configHandler.Get)).ServeHTTP(w, r)
	}).Methods("GET")
	router.HandleFunc("/configs", func(w http.ResponseWriter, r *http.Request) {
		handlers.RateLimit(limiter, http.HandlerFunc(configHandler.Add)).ServeHTTP(w, r)
	}).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", func(w http.ResponseWriter, r *http.Request) {
		handlers.RateLimit(limiter, http.HandlerFunc(configHandler.Delete)).ServeHTTP(w, r)
	}).Methods("DELETE")

	// ConfigGroup routes
	router.HandleFunc("/configgroups/{name}/{version}", func(w http.ResponseWriter, r *http.Request) {
		handlers.RateLimit(limiter, http.HandlerFunc(configGroupHandler.Get)).ServeHTTP(w, r)
	}).Methods("GET")
	router.HandleFunc("/configgroups", func(w http.ResponseWriter, r *http.Request) {
		handlers.RateLimit(limiter, http.HandlerFunc(configGroupHandler.Add)).ServeHTTP(w, r)
	}).Methods("POST")
	router.HandleFunc("/configgroups/{name}/{version}", func(w http.ResponseWriter, r *http.Request) {
		handlers.RateLimit(limiter, http.HandlerFunc(configGroupHandler.Delete)).ServeHTTP(w, r)
	}).Methods("DELETE")

	router.HandleFunc("/configgroups/{name}/{version}", func(w http.ResponseWriter, r *http.Request) {
		handlers.RateLimit(limiter, http.HandlerFunc(configGroupHandler.Update)).ServeHTTP(w, r)
	}).Methods("PUT")

	router.HandleFunc("/configgroups/{name}/{version}/{labels}", func(w http.ResponseWriter, r *http.Request) {
		handlers.RateLimit(limiter, http.HandlerFunc(configGroupHandler.GetConfigInListByLabels)).ServeHTTP(w, r)
	}).Methods("GET")

	router.HandleFunc("/configgroups/{name}/{version}/{labels}", func(w http.ResponseWriter, r *http.Request) {
		handlers.RateLimit(limiter, http.HandlerFunc(configGroupHandler.DeleteConfigInListByLabels)).ServeHTTP(w, r)
	}).Methods("DELETE")

	// Handler za /shutdown
	router.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Primljen zahtev za gašenje...")

		// Odgovor klijentu da je zahtev za gašenje primljen
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server će biti gašen...")

		// Čekamo određeni vremenski period pre gašenja servera
		time.AfterFunc(shutdownTimer, func() {
			// Šaljemo shutdown signal nakon isteka vremenskog perioda
			shutdownChan <- os.Interrupt
		})
	}).Methods("POST")

	//http.Handle("/", handlers.RateLimit(limiter,router.ServeHTTP))

	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}

	go func() {
		fmt.Println("Pokretanje servera na", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Greška pri pokretanju servera:", err)
			os.Exit(1)
		}
	}()

	// Čekamo shutdown signal
	<-shutdownChan
	fmt.Println("Primljen shutdown signal. Gašenje aplikacije...")

	// Countdown sa tajmerom za graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gašenje servera
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Greška pri gašenju servera:", err)
		os.Exit(1)
	}

	fmt.Println("Gašenje servera uspešno.")
	os.Exit(0) // Izlaz iz aplikacije
}

func simulateOperations(service services.ConfigService, repoS model.ConfigRepository) {
	//	(groupService services.ConfigGroupService,  repo2 model.ConfigGroupRepository)

	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5431"
	config := model.Config{
		Name:    "db_config",
		Version: 2,
		Params:  params,
	}
	config2, _ := repoS.NewConfigFromLiteral("Ime 23 params=paramsss params2=params22")
	fmt.Println("Pokretanje simulacija")
	service.Add(config2)
	service.Add(config)

	/*configData := []string{
		"GroupA",        // Naziv grupe
		"1",             // Verzija grupe
		"param1=value1", // Konfiguracioni parametar 1
		"param2=value2", // Konfiguracioni parametar 2
	}

	configGroup, _ := repo2.ParseConfigData(configData)
	fmt.Println("Naziv grupe:", configGroup.Name)
	fmt.Println("Verzija grupe:", configGroup.Version)
	fmt.Println("Konfiguracioni parametri:")
	for _, config := range configGroup.ConfigInList {
		fmt.Printf("- Naziv: %s, Vrednost: %s\n", config.Name, config.Params["value"])
	}*/

	retrievedConfig1, _ := service.Get("Ime", 23)
	retrievedConfig, _ := service.Get("db_config", 2)

	jsonData1, _ := json.Marshal(retrievedConfig1)
	jsonData, _ := json.Marshal(retrievedConfig)

	fmt.Println("Matijaaa")
	fmt.Println(string(jsonData))
	fmt.Println(string(jsonData1))
}
