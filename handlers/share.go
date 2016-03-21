package handlers
import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2/bson"
	"github.com/mshockwave/share-sound-api-server/common"
)

func GetSessionUserId(req *http.Request) (bson.ObjectId, error){
	if v, err := common.GetSessionValue(req, common.USER_ID_SESSION_KEY); err != nil || v == nil{
		return bson.ObjectId(""), errors.New("Invalid session id format")
	}else{
		if str, found := v.(string); found {
			if bson.IsObjectIdHex(str) {
				return bson.ObjectIdHex(str), nil
			}else{
				return bson.ObjectId(""), errors.New("Invalid session id format")
			}
		}else{
			return bson.ObjectId(""), errors.New("Invalid session id format")
		}
	}
}

func AuthVerifierWrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request){
		if _, err := GetSessionUserId(req); err != nil {
			r := common.SimpleResult{
				Message: "Error",
				Description: "Please Login First",
			}
			common.ResponseStatusAsJson(resp, 403, &r)
			return
		}

		handler(resp, req)
	}
}
