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

func CreateCommunity(response http.ResponseWriter, request *http.Request) {
	err := middleware.Authentication(response, request)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}
	response.Header().Set("Content-Type", "application/json")
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)

	var community models.Community
	//var user models.User
	var account models.Account
	authHeader, _ := middleware.FromAuthHeader(request)

	json.NewDecoder(request.Body).Decode(&community)

	//json.Unmarshal([]byte((request.FormValue("json"))), &community)
	//request.ParseMultipartForm(10 << 20)

	err = accountCollection.FindOne(ctx, bson.M{"token": authHeader}).Decode(&account)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`{"message":"` + " user not found" + `"}`))
		return
	}

	/*file, handler, err := request.FormFile("photo")
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error Retrieving the File `))
		response.Write([]byte(`Error ` + err.Error()))
		return
	}
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
	}*/
	_, err = services.GetUserOneByUserID(&account.User_id)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}

	//	community.Photo = &fileid
	//community.City = user.City
	//community.Participants = append(community.Participants, user.User_id)
	//community.Admin = &user.User_id

	err = services.CreateCommunity(&community)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`you allready have this community`))
		return
	}
}
