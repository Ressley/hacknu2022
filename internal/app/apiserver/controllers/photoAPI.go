package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/Ressley/hacknu/internal/app/apiserver/services"
)

func UploadPhoto(response http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(10 << 20)
	file, handler, err := request.FormFile("photo")
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
	id, err := services.UploadPhoto(handler.Filename, fileBytes)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}
	response.Write([]byte(`{"fileid" : "` + id + `"}`))
}

func DownloadPhoto(response http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()
	id := query.Get("fileid")

	bytes, err := services.DownloadPhoto(id)
	if err != nil {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(`Error ` + err.Error()))
		return
	}
	response.Write(bytes)
}
