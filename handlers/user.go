package handlers

import (
	"time"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mshockwave/share-sound-api-server/common"
	dbSchema "github.com/mshockwave/share-sound-api-server/datastore/schema"
	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
	db "github.com/mshockwave/share-sound-api-server/datastore"
	"google.golang.org/cloud/datastore"
)

type registerForm struct {
	UserName	string `schema:"username"`
	Email		string `schema:"email"`
	Password	string `schema:"password"`
}
func (this registerForm) CopyToSchema() dbSchema.User {
	pwdHash, e := bcrypt.GenerateFromPassword([]byte(this.Password),bcrypt.DefaultCost)
	if e != nil {
		return dbSchema.User{}
	}

	return dbSchema.User{
		Username: this.UserName,
		Email: this.Email,

		Auth: dbSchema.UserAuth{
			PasswordBcryptHash: string(pwdHash),
			PasswordBcryptCost: bcrypt.DefaultCost,
		},

		CreatedTimeStamp: time.Now(),

		Thumbnail: common.DEFAULT_USER_THUMBNAIL_PATH,
	}
}

type resultUserProfile struct {
	Username    string
	Email       string
	CreatedDate time.Time
	Thumbnail   string "" //TODO
}
func (this *resultUserProfile) FromDbSchema(dbForm *dbSchema.User){
	this.Username = dbForm.Username
	this.Email = dbForm.Email
	this.CreatedDate = dbForm.CreatedTimeStamp
	this.Thumbnail = dbForm.Thumbnail
}

func handleRegister(resp http.ResponseWriter, req *http.Request){
	var err error

	err = req.ParseForm()
	if err != nil {
		common.ResponseStatusAsJson(resp, 400, &common.SimpleResult{
			Message: "Error",
			Description: "Wrong Form Format",
		})
		return
	}

	formDecoder := schema.NewDecoder()

	regForm := registerForm{}
	err = formDecoder.Decode(&regForm, req.PostForm)
	if err != nil {
		common.ResponseStatusAsJson(resp, 400, &common.SimpleResult{
			Message: "Error",
			Description: "Wrong Form Format",
		})
		return
	}

	regDbForm := regForm.CopyToSchema()

	if client, e := db.GetNewDataStoreClient(); e == nil {

		//Check existence
		query := datastore.NewQuery(dbSchema.USER_PROFILE_KIND).
					Ancestor(UserProfileRootKey).
					Filter("Email =", regDbForm.Email)
		it := client.Run(query)
		u := dbSchema.User{}
		if _, exist := it.Next(&u); exist == nil {
			//User Exist
			common.ResponseStatusAsJson(resp, 400, &common.SimpleResult{
				Message: "Error",
				Description: "User Existed",
			})
			return
		}

		key := client.NewKey(dbSchema.USER_PROFILE_KIND, regDbForm.Email, 0, UserProfileRootKey)
		_, err := client.Client.Put(client.Ctx, key, &regDbForm)
		if err != nil {
			common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
				Message: "Error",
				Description: "Insert DB Error",
			})
			return
		}

		//Setup session
		if e := common.SetSessionValue(req, resp, common.USER_ID_SESSION_KEY, key.Name()); e != nil{
			common.LogE.Println("Error setting session")
		}

		result := &resultUserProfile{}
		result.FromDbSchema(&regDbForm)
		common.ResponseOkAsJson(resp, result)
	}else{
		common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
			Message: "Error",
			Description: "Internal DB Error",
		})
	}
}

type loginForm struct {
	Email		string `schema:"email"`
	Password	string `schema:"password"`
}
func handleLogin(resp http.ResponseWriter, req *http.Request){
	var err error

	err = req.ParseForm()
	if err != nil {
		common.LogE.Printf("Login failed: %s\n", err.Error())
		common.ResponseStatusAsJson(resp, 400, &common.SimpleResult{
			Message: "Error",
			Description: "Wrong Login Format",
		})
		return
	}

	formDecoder := schema.NewDecoder()

	lgForm := loginForm{}
	err = formDecoder.Decode(&lgForm, req.PostForm)
	if err != nil {
		common.LogE.Printf("Login failed: %s\n", err.Error())
		common.ResponseStatusAsJson(resp, 400, &common.SimpleResult{
			Message: "Error",
			Description: "Wrong Login Format",
		})
		return
	}

	if client, e := db.GetNewDataStoreClient(); e == nil {
		//Query
		q := datastore.NewQuery(dbSchema.USER_PROFILE_KIND).
				Ancestor(UserProfileRootKey).
				Filter("Email =", lgForm.Email)
		it := client.Run(q)
		result := dbSchema.User{}
		var resultKey *datastore.Key
		if k, exist := it.Next(&result); exist != nil {
			//Not exist
			common.ResponseStatusAsJson(resp, 403, &common.SimpleResult{
				Message: "Error",
				Description: "Login Failed",
			})
			return
		}else{
			resultKey = k
		}

		//Check password
		correct := bcrypt.CompareHashAndPassword([]byte(result.Auth.PasswordBcryptHash), []byte(lgForm.Password))
		if correct != nil {
			common.ResponseStatusAsJson(resp, 403, &common.SimpleResult{
				Message: "Error",
				Description: "Login Failed",
			})
			return
		}

		//Success
		if e := common.SetSessionValue(req, resp, common.USER_ID_SESSION_KEY, resultKey.Name()); e != nil {
			common.LogE.Println("Error setting session")
		}
		resProfile := &resultUserProfile{}
		resProfile.FromDbSchema(&result)
		common.ResponseOkAsJson(resp, resProfile)
	}else{
		common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
			Message: "Error",
			Description: "Internal DB Error",
		})
	}

}

func handleCheckLogin(resp http.ResponseWriter, req *http.Request) {
	common.ResponseOkAsJson(resp, &common.SimpleResult{
		Message: "Is Login",
		Description: "True",
	})
}

func handleProfile(resp http.ResponseWriter, req *http.Request){
	if(req.Method == "GET"){
		id, err := GetSessionUserId(req)
		if err != nil {
			common.LogE.Printf("Error fetching user sessoin id: %s\n", err)
			common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
				Message: "Error",
			})
		}else{
			query := datastore.NewQuery(dbSchema.USER_PROFILE_KIND).
						Ancestor(UserProfileRootKey).
						Filter("Email =", id)
			if client, e := db.GetNewDataStoreClient(); e == nil {
				it := client.Run(query)
				u := dbSchema.User{}
				if _, exist := it.Next(&u); exist == nil {
					resProfile := &resultUserProfile{}
					resProfile.FromDbSchema(&u)
					common.ResponseOkAsJson(resp, resProfile)
				}else{
					common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
						Message: "Error",
						Description: "User Not Found",
					})
				}
			}else{
				common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
					Message: "Error",
					Description: "Internal DB Error",
				})
			}
		}
	}else{
		//Not supported yet
		common.ResponseStatusAsJson(resp, 404, nil)
	}
}

func handleLogout(resp http.ResponseWriter, req *http.Request){
	if e := common.SetSessionValue(req, resp, common.USER_ID_SESSION_KEY, nil); e != nil {
		common.LogE.Println("Error setting session")
	}
	common.ResponseOkAsJson(resp, &common.SimpleResult{
		Message: "Logout",
		Description: "Success",
	})
}

func ConfigureUserHandlers(router *mux.Router){
	router.HandleFunc("/register", handleRegister)
	router.HandleFunc("/login", handleLogin)
	router.HandleFunc("/login/status", AuthVerifierWrapper(handleCheckLogin))
	router.HandleFunc("/profile", AuthVerifierWrapper(handleProfile))
	router.HandleFunc("/logout", AuthVerifierWrapper(handleLogout))
}