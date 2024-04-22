package handlers

import (
	"ProjectModule/model"
	"ProjectModule/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ConfigGroupHandler struct {
	service services.ConfigGroupService
}

func NewConfigGroupHandler(service services.ConfigGroupService) ConfigGroupHandler {
	return ConfigGroupHandler{
		service: service,
	}
}

// GET /configgroups/{name}/{version}
func (c ConfigGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Get the name and version from URL parameters
	name := mux.Vars(r)["name"]
	versionStr := mux.Vars(r)["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version format", http.StatusBadRequest)
		return
	}

	// Call the service method to get the config group
	configGroup, err := c.service.GetGroup(name, version)
	if err != nil {
		http.Error(w, "Config group not found", http.StatusNotFound)
		return
	}

	// Marshal the config group to JSON
	resp, err := json.Marshal(configGroup)
	if err != nil {
		http.Error(w, "Failed to marshal config group", http.StatusInternalServerError)
		return
	}

	// Set response headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// POST /configgroups
func (c ConfigGroupHandler) Add(w http.ResponseWriter, r *http.Request) {
	var configGroup model.ConfigGroup

	// Decode the JSON payload into the ConfigGroup struct
	err := json.NewDecoder(r.Body).Decode(&configGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the AddGroup method of the service to add the new config group
	c.service.AddGroup(configGroup)

	// Respond with success status
	w.WriteHeader(http.StatusCreated)
}

// DELETE /configgroups/{name}/{version}
func (c ConfigGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version format", http.StatusBadRequest)
		return
	}

	// Check if the config group exists
	_, err = c.service.GetGroup(name, version)
	if err != nil {
		http.Error(w, "Config group not found", http.StatusNotFound)
		return
	}

	// Call the DeleteGroup method of the service to delete the config group
	err = c.service.DeleteGroup(name, version)
	if err != nil {
		http.Error(w, "Failed to delete config group", http.StatusInternalServerError)
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}
