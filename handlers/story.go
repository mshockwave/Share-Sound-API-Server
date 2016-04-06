package handlers

import (
	"net/http"
	"mime"
	"encoding/base64"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/mshockwave/share-sound-api-server/common"
	dbSchema "github.com/mshockwave/share-sound-api-server/datastore/schema"
	db "github.com/mshockwave/share-sound-api-server/datastore"
	"github.com/mshockwave/share-sound-api-server/handlers/protos"
	"github.com/mshockwave/share-sound-api-server/storage"
	"github.com/golang/protobuf/proto"
)

func handleUpload(resp http.ResponseWriter, req *http.Request){
	if !(req.Method == "POST" || req.Method == "PUT") {
		//Not supported
		common.ResponseStatusAsJson(resp, 404, nil)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if len(contentType) <= 0 {
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

	uid,_ := GetSessionUserId(req)
	story := protos.Story{}
	if err := proto.Unmarshal(bodyContent, &story); err != nil {
		common.ResponseStatusAsJson(resp, 400, &common.SimpleResult{
			Message: "Error",
			Description: "Wrong binary layout format",
		})
		return
	}

	dbStory := dbSchema.StoryMeta{
		Id: dbSchema.HashId(common.GetDefaultSecureHash()),
		UploaderEmail: uid,

		Title: story.Title,
		Description: story.Description,
	}

	if client, err := storage.GetNewStorageClient(); err == nil {
		var audioAttachments []dbSchema.AudioAttachmentMeta
		for _,item := range story.AudioAttachments {
			audio := dbSchema.AudioAttachmentMeta{}
			if err := (&audio).FromProtoBuf(item, client); err == nil {
				audioAttachments = append(audioAttachments, audio)
			}
		}
		dbStory.AudioAttachments = audioAttachments

		var imgAttachments []dbSchema.ImageAttachmentMeta
		for _,item := range story.ImageAttachments {
			img := dbSchema.ImageAttachmentMeta{}
			if err := (&img).FromProtoBuf(item, client); err == nil {
				imgAttachments = append(imgAttachments, img)
			}
		}
		dbStory.ImageAttachments = imgAttachments
	}else{
		common.LogE.Println("Get storage client error: " + err.Error())
		common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
			Message: "Error",
		})
		return
	}

	if client, err := db.GetNewDataStoreClient(); err == nil {
		key := client.NewKey(dbSchema.STORY_KIND, dbStory.Id.String(), 0, StoryRootKey)
		_, e := client.Client.Put(client.Ctx, key, &dbStory)
		if e != nil {
			common.LogE.Println("Insert story failed: " + e.Error())
			common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
				Message: "Error",
				Description: "Upload story failed",
			})
			return
		}

		common.ResponseOkAsJson(resp, &common.SimpleResult{
			Message: "OK",
		})
	}else{
		common.LogE.Println("Get datastore client error: " + err.Error())
		common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
			Message: "Error",
		})
		return
	}
}

func ConfigureStoryHandler(router *mux.Router){
	router.HandleFunc("/upload", AuthVerifierWrapper(handleUpload))
}
