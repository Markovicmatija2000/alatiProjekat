package handlers

import (
	"ProjectModule/model"
	"ProjectModule/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type ConfigGroupHandler struct {
	groupService        services.ConfigGroupService
	configInListService services.ConfigInListService
}

func NewConfigGroupHandler(service services.ConfigGroupService) ConfigGroupHandler {
	return ConfigGroupHandler{
		groupService: service,
	}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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
	configGroup, err := c.groupService.GetGroup(name, version)
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
	c.groupService.AddGroup(configGroup)

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
	_, err = c.groupService.GetGroup(name, version)
	if err != nil {
		http.Error(w, "Config group not found", http.StatusNotFound)
		return
	}

	// Call the DeleteGroup method of the service to delete the config group
	err = c.groupService.DeleteGroup(name, version)
	if err != nil {
		http.Error(w, "Failed to delete config group", http.StatusInternalServerError)
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}

// PUT /configgroups/{name}/{version}
func (c ConfigGroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version format", http.StatusBadRequest)
		return
	}

	// Decode the JSON payload into the ConfigInList struct
	var configInList model.ConfigInList
	err = json.NewDecoder(r.Body).Decode(&configInList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the existing config group
	existingGroup, err := c.groupService.GetGroup(name, version)
	if err != nil {
		http.Error(w, "Config group not found", http.StatusNotFound)
		return
	}

	// Add the new configInList to the existing group
	existingGroup.ConfigInList = append(existingGroup.ConfigInList, configInList)

	// Call the UpdateGroup method of the service to update the config group
	err = c.groupService.UpdateGroup(existingGroup)
	if err != nil {
		http.Error(w, "Failed to update config group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GET /configgroups/{name}/{version}/{labels}
func (c ConfigGroupHandler) GetConfigInListByLabels(w http.ResponseWriter, r *http.Request) {
	// Parse the request URL parameters
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version format", http.StatusBadRequest)
		return
	}
	labelsStr := vars["labels"]
	labelPairs := strings.Split(labelsStr, ";")
	labels := make([]model.ConfigInList, 0, len(labelPairs))
	for _, pair := range labelPairs {
		labelParts := strings.Split(pair, ":")
		if len(labelParts) != 2 {
			http.Error(w, "Invalid label format", http.StatusBadRequest)
			return
		}
		labels = append(labels, model.ConfigInList{Name: labelParts[0], Params: map[string]string{"value": labelParts[1]}})
	}

	// Calling the service method to get the config in list by labels

	configInLists, err := c.groupService.GetConfigInListByLabels(name, version, labels)
	if err != nil {
		http.Error(w, "Failed to get config in list by labels", http.StatusInternalServerError)
		return
	}

	// Marshaling the config in lists to JSON

	resp, err := json.Marshal(configInLists)
	if err != nil {
		http.Error(w, "Failed to marshal config in lists", http.StatusInternalServerError)
		return
	}

	// Seting response headers and write response

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// DELETE /configgroups/{name}/{version}/{labels}
func (c ConfigGroupHandler) DeleteConfigInListByLabels(w http.ResponseWriter, r *http.Request) {
	// Parse the request URL parameters
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version format", http.StatusBadRequest)
		return
	}
	labelsStr := vars["labels"]
	labelPairs := strings.Split(labelsStr, ";")
	labels := make([]model.ConfigInList, 0, len(labelPairs))
	for _, pair := range labelPairs {
		labelParts := strings.Split(pair, ":")
		if len(labelParts) != 2 {
			http.Error(w, "Invalid label format", http.StatusBadRequest)
			return
		}
		labels = append(labels, model.ConfigInList{Name: labelParts[0], Params: map[string]string{"value": labelParts[1]}})
	}

	// Call the service method to delete the config in list by labels
	err = c.groupService.DeleteConfigInListByLabels(name, version, labels)
	if err != nil {
		http.Error(w, "Failed to delete config in list by labels", http.StatusInternalServerError)
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}
