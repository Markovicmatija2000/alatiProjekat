package handlers

import (
	"ProjectModule/model"
	"ProjectModule/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ConfigHandler struct {
	service services.ConfigService
}

func NewConfigHandler(service services.ConfigService) ConfigHandler {
	return ConfigHandler{
		service: service,
	}
}

// GET /configs/{name}/{version}
func (c ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	// dobavi naziv i verziju
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// pozovi servis metodu
	config, err := c.service.Get(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// vrati odgovor
	resp, err := json.Marshal(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content−Type", "application/json")
	w.Write(resp)
}

func (c ConfigHandler) Add(w http.ResponseWriter, r *http.Request) {
	var config model.Config

	// Decode the JSON payload into the Config struct
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the Add method of the service to add the new configuration
	c.service.Add(config)

	// Respond with success status
	w.WriteHeader(http.StatusCreated)
}

// DELETE /configs/{name}/{version}
func (c ConfigHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version format", http.StatusBadRequest)
		return
	}

	// Check if the configuration exists
	_, err = c.service.Get(name, version)
	if err != nil {
		http.Error(w, "Config not found", http.StatusNotFound)
		return
	}

	// Delete the configuration
	err = c.service.Delete(name, version)
	if err != nil {
		http.Error(w, "Failed to delete config", http.StatusInternalServerError)
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}
