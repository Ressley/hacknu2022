package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Ressley/hacknu/internal/app/apiserver/middleware"
	"github.com/Ressley/hacknu/internal/app/apiserver/models"
	"github.com/Ressley/hacknu/internal/app/apiserver/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateEvent(response http.ResponseWriter, request *http.Request) {
	err := middleware.Authentication(response, request)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}
	response.Header().Set("Content-Type", "application/json")
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	var event models.Event
	var user models.User
	var account models.Account
	authHeader, _ := middleware.FromAuthHeader(request)

	json.Unmarshal([]byte((request.FormValue("json"))), &event)
	request.ParseMultipartForm(10 << 20)
	err = accountCollection.FindOne(ctx, bson.M{"token": authHeader}).Decode(&account)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"message":"` + " user not found" + `"}`))
		return
	}

	file, handler, err := request.FormFile("photo")
	if err == nil {

		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte(`Error ` + err.Error()))
			return
		}

		fileid, err := services.UploadFile(handler.Filename, fileBytes)
		if err != nil {
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte(`Error ` + err.Error()))
			return
		}
		event.Photo = &fileid
	}
	user, err = services.GetUserOneByUserID(&account.User_id)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}

	communityType, err := services.GetCommunityByName(event.TypeName)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + "Community Does not exist"))
		return
	}

	//event.City = user.City
	event.Participants = append(event.Participants, user.User_id)
	event.Type = &communityType
	event.Admin = &user.User_id
	event.ID = primitive.NewObjectID()

	eventID, err := services.CreateEvent(&event)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`you allready have this community`))
		return
	}
	user.Events = append(user.Events, eventID.InsertedID.(primitive.ObjectID).Hex())
	err = services.UpdateUserOne(&user)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`you allready have this community`))
		return
	}
}

func DeleteEvent(response http.ResponseWriter, request *http.Request) {
	err := middleware.Authentication(response, request)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}
	response.Header().Set("Content-Type", "application/json")
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	var account models.Account
	authHeader, _ := middleware.FromAuthHeader(request)

	err = accountCollection.FindOne(ctx, bson.M{"token": authHeader}).Decode(&account)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"message":"` + " user not found" + `"}`))
		return
	}
	query := request.URL.Query()
	eventID := query.Get("eventid")

	event, err := services.GetEventByID(&eventID)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"Error":"event with ` + eventID + ` id does not exist"}`))
		return
	}

	if event.Admin != &account.User_id {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"Error":"you are not admin"}`))
		return
	}

	err = services.DeleteEvent(event.ID.Hex())
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"Error":"event with ` + eventID + ` id does not exist"}`))
		return
	}
}

func GetEvent(response http.ResponseWriter, request *http.Request) {
	err := middleware.Authentication(response, request)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	var account models.Account
	authHeader, _ := middleware.FromAuthHeader(request)

	err = accountCollection.FindOne(ctx, bson.M{"token": authHeader}).Decode(&account)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"message":"` + " user not found" + `"}`))
		return
	}

	query := request.URL.Query()
	eventID := query.Get("eventid")
	var community models.Community
	if eventID != "" {
		event, err := services.GetEventByID(&eventID)
		if err != nil {
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte(`{"Error":"event with ` + eventID + ` id does not exist"}`))
			return
		}
		json.NewEncoder(response).Encode(event)
		return
	}
	json.NewDecoder(request.Body).Decode(&community)
	if *community.Name == "/All" {
		event, err := services.GetEventAll()
		if err != nil {
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte(`{"Error":"event with ` + eventID + ` id does not exist"}`))
			return
		}
		json.NewEncoder(response).Encode(event)
		return
	}
	event, err := services.GetCommunityByFilter(community.Name, &account)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"Error":"event with ` + eventID + ` id does not exist"}`))
		return
	}
	json.NewEncoder(response).Encode(event)
	return
}
