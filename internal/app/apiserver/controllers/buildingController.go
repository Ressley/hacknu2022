package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Ressley/hacknu/internal/app/apiserver/middleware"
	"github.com/Ressley/hacknu/internal/app/apiserver/models"
	"github.com/Ressley/hacknu/internal/app/apiserver/services"
)

func CreateBuilding(response http.ResponseWriter, request *http.Request) {
	err := middleware.Authentication(response, request)
	if err != nil {
		return
	}
	response.Header().Set("Content-Type", "application/json")

	var building models.Building
	json.NewDecoder(request.Body).Decode(&building)

	if err := ValidateStruct(building); err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		json, _ := json.Marshal(err)
		response.Write([]byte(`{"Error" : ` + string(json) + `}`))
		return
	}

	err = services.CreateBuilding(&building)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"Error" : "building with this name allready exists"}`))
		return
	}
	json.NewEncoder(response).Encode(building)
	return
}

func GetBuilding(response http.ResponseWriter, request *http.Request) {
	// err := middleware.Authentication(response, request)
	// if err != nil {
	// 	return
	// }
	response.Header().Set("Content-Type", "application/json")
	var building models.Building

	query := request.URL.Query()
	building_id := query.Get("building_id")

	building, err = services.GetBuildingByID(&building_id)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"Error":"building with ` + building_id + ` id does not exist"}`))
		return
	}
	json.NewEncoder(response).Encode(building)
}

func ListBuildings(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var buildings []models.BuildingMeta

	buildings, err = services.ListBuildings()
	json.NewEncoder(response).Encode(buildings)
}
