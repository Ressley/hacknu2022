package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Ressley/hacknu/internal/app/apiserver/middleware"
	"github.com/Ressley/hacknu/internal/app/apiserver/models"
	"github.com/Ressley/hacknu/internal/app/apiserver/services"
)

func CreateApartment(response http.ResponseWriter, request *http.Request) {
	err := middleware.Authentication(response, request)
	if err != nil {
		return
	}
	response.Header().Set("Content-Type", "application/json")

	query := request.URL.Query()
	building_id := query.Get("building_id")

	building, err := services.GetBuildingByID(&building_id)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"Error":"building with ` + building_id + ` id does not exist"}`))
		return
	}

	var apartment models.Apartment
	json.NewDecoder(request.Body).Decode(&apartment)

	if err := ValidateStruct(apartment); err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		json, _ := json.Marshal(err)
		response.Write([]byte(`{"Error" : ` + string(json) + `}`))
		return
	}
	services.CreateApartment(&apartment, &building.ID)
	services.AppendApartment(&apartment.ID, &building)
	json, _ := json.Marshal(apartment)
	response.Write([]byte(`{"data":` + string(json) + `}`))
}

func DeleteApartment(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("id")
	apartment, err := services.GetApartmentByID(&id)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"Error":"apartment with ` + id + ` id does not exist"}`))
		return
	}

	err = services.RemoveApartment(&apartment)

	err = services.DeleteApartment(&(apartment.ID))
}
