package controllers

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Ressley/hacknu/internal/app/apiserver/services"
)

func UploadFile(response http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(10 << 20)
	file, handler, err := request.FormFile("file")
	log.Print(file)
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
	id, err := services.UploadFile(handler.Filename, fileBytes)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}
	response.Write([]byte(`{"fileid" : "` + id + `"}`))
}

func DownloadFile(response http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()
	id := query.Get("fileid")

	bytes, err := services.DownloadFile(id)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}
	response.Write(bytes)
}
