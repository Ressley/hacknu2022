package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Ressley/hacknu/internal/app/apiserver/helpers"
	"github.com/Ressley/hacknu/internal/app/apiserver/models"
	"github.com/Ressley/hacknu/internal/app/apiserver/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var client, err = helpers.GetMongoClient()
var accountCollection *mongo.Collection = client.Database(helpers.DB).Collection(helpers.ACCOUNTS)

func GetAccountCollection() *mongo.Collection {
	return accountCollection
}

// HashPassword ...
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

// VerifyPassword ...
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or password is incorrect")
		check = false
	}

	return check, msg
}

// CreateUser is the api used to tget a single user
func SignUp(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	nullFirstName := ""
	nullLastName := ""
	nullLogin := ""
	nullPassword := ""

	account := models.Account{
		First_name: &nullFirstName,
		Last_name:  &nullLastName,
		Login:      &nullLogin,
		Password:   &nullPassword,
	}

	json.Unmarshal([]byte((request.FormValue("json"))), &account)
	request.ParseMultipartForm(10 << 20)

	json.NewDecoder(request.Body).Decode(&account)

	if *account.Last_name == "" || *account.First_name == "" || *account.Login == "" || *account.Password == "" {
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte(`{"message":"` + "FirstName, LastName, Login and Password cannot be null" + `"}`))
		return
	}

	var fileid string

	file, handler, err := request.FormFile("photo")
	if err == nil {

		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte(`Error ` + err.Error()))
			return
		}

		fileid, err = services.UploadFile(handler.Filename, fileBytes)
		if err != nil {
			response.WriteHeader(http.StatusMethodNotAllowed)
			response.Write([]byte(`Error ` + err.Error()))
			return
		}
	}

	*account.Password = HashPassword(*account.Password)
	account.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	account.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	account.ID = primitive.NewObjectID()
	account.User_id = account.ID.Hex()

	token, refreshToken, _ := helpers.GenerateAllTokens(*account.Login, *account.First_name, *account.Last_name, account.User_id)
	account.Token = &token
	account.Refresh_token = &refreshToken

	err = services.CreateUser(&account, &fileid)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	_, err = accountCollection.InsertOne(ctx, account)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	response.Write([]byte(`{"token" : "` + token + `"}`))
}

func Login(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	var account models.Account
	var dbAccount models.Account
	var uid string
	json.NewDecoder(request.Body).Decode(&account)

	numberRegex := regexp.MustCompile(`^87[0-7][0-8].`)

	if numberRegex.Match([]byte(*account.Login)) && len(*account.Login) == 11 {
		user, err := services.GetUserOneByNumber(account.Login)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		} else {
			uid = user.User_id
		}
	}
	id, _ := primitive.ObjectIDFromHex(uid)
	err := accountCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&dbAccount)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	userPass := []byte(*account.Password)
	dbPass := []byte(*dbAccount.Password)
	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)
	if passErr != nil {
		log.Println(passErr)
		response.Write([]byte(`{"response":"Wrong Password"}`))
		return
	}
	jwtToken, refreshToken, _ := helpers.GenerateAllTokens(*dbAccount.Login, *dbAccount.First_name, *dbAccount.Last_name, dbAccount.User_id)

	helpers.UpdateAllTokens(jwtToken, refreshToken, dbAccount.User_id)
	response.Write([]byte(`{"token" : "` + jwtToken + `"}`))
}
