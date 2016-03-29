package handlers

import (
	"net/http"
	"mime"
	"encoding/base64"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/mshockwave/share-sound-api-server/common"
)

func handleUpload(resp http.ResponseWriter, req *http.Request){
	if !(req.Method == "POST" || req.Method == "PUT") {
		//Not supported
		common.ResponseStatusAsJson(resp, 404, nil)
		return
	}

	contentType, exist := req.Header["Content-Type"]
	if !exist {
		common.ResponseStatusAsJson(resp, 400, &common.SimpleResult{
			Message: "Error",
			Description: "Must specify Content-Type",
		})
		return
	}

	mediaType,_, err := mime.ParseMediaType(contentType)
	if err != nil {
		common.ResponseStatusAsJson(resp, 400, &common.SimpleResult{
			Message: "Error",
			Description: "Unrecognized Content-Type: " + contentType,
		})
		return
	}
	if !(mediaType == "application/octet-stream" || mediaType == "application/base64") {
		common.ResponseStatusAsJson(resp, 400, &common.SimpleResult{
			Message: "Error",
			Description: "Only application/octet-stream and application/base64 are supported now",
		})
		return
	}


	bodyRaw, _ := ioutil.ReadAll(req.Body)
	//defer req.Body.Close()

	var bodyContent []byte

	//Decode base64
	if(mediaType == "application/base64") {
		decoder := base64.StdEncoding

		if _, err := decoder.Decode(bodyContent, bodyRaw); err != nil{
			common.ResponseStatusAsJson(resp, 400, &common.SimpleResult{
				Message: "Error",
				Description: "Wrong base64 encoding",
			})
			return
		}
	}else{
		bodyContent = bodyRaw
	}


}

func ConfigureStoryHandler(router *mux.Router){
	router.HandleFunc("/upload", AuthVerifierWrapper(handleUpload))
}
