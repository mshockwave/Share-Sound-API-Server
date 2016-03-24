package handlers
import (
	"errors"
	"net/http"

	"github.com/mshockwave/share-sound-api-server/common"
	"github.com/mshockwave/share-sound-api-server/datastore/schema"
	db "github.com/mshockwave/share-sound-api-server/datastore"
	"google.golang.org/cloud/datastore"
)

var(
	UserProfileRootKey *datastore.Key
)

func init(){
	//Module initializer

	UserProfileRootKey = datastore.NewKey(db.GetContext(), schema.USER_PROFILE_KIND, "admin@sharesound.org", 0, nil)
}

func GetSessionUserId(req *http.Request) (string, error){
	if v, err := common.GetSessionValue(req, common.USER_ID_SESSION_KEY); err != nil || v == nil{
		return "", errors.New("Invalid session id format")
	}else{
		if str, found := v.(string); found {
			return str, nil
		}else{
			return "", errors.New("Invalid session id format")
		}
	}
}

func AuthVerifierWrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request){
		if _, err := GetSessionUserId(req); err != nil {
			//common.LogE.Printf("Error session value: %s\n", err.Error())
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
