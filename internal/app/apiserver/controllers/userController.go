package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ressley/hacknu/internal/app/apiserver/middleware"
	"github.com/Ressley/hacknu/internal/app/apiserver/models"
	"github.com/Ressley/hacknu/internal/app/apiserver/services"
	"go.mongodb.org/mongo-driver/bson"
)

func Follow(response http.ResponseWriter, request *http.Request) {
	err := middleware.Authentication(response, request)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"Error":"` + err.Error() + `"}`))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	//var event models.Event
	var user models.User
	var account models.Account
	authHeader, _ := middleware.FromAuthHeader(request)

	json.NewDecoder(request.Body).Decode(&user)

	query := request.URL.Query()
	userID := query.Get("userid")

	err = accountCollection.FindOne(ctx, bson.M{"token": authHeader}).Decode(&account)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"message":"` + " user not found" + `"}`))
		return
	}

	err = services.Follow(&userID, &account.User_id)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"Error":"` + err.Error() + `"}`))
		return
	}

}

func Unfollow(response http.ResponseWriter, request *http.Request) {

	err := middleware.Authentication(response, request)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	//var event models.Event
	var user models.User
	var account models.Account
	authHeader, _ := middleware.FromAuthHeader(request)

	json.NewDecoder(request.Body).Decode(&user)

	query := request.URL.Query()
	userID := query.Get("userid")

	err = accountCollection.FindOne(ctx, bson.M{"token": authHeader}).Decode(&account)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"message":"` + " user not found" + `"}`))
		return
	}

	err = services.Unfollow(&userID, &account.User_id)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"Error":"` + err.Error() + `"}`))
		return
	}
}

func GetUser(response http.ResponseWriter, request *http.Request) {

	err := middleware.Authentication(response, request)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	//var event models.Event
	var account models.Account
	authHeader, _ := middleware.FromAuthHeader(request)

	err = accountCollection.FindOne(ctx, bson.M{"token": authHeader}).Decode(&account)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"message":"` + " user not found" + `"}`))
		return
	}

	user, err := services.GetUserOneByUserID(&account.User_id)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"message":"` + " user not found" + `"}`))
		return
	}
	json.NewEncoder(response).Encode(user)
}
