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
	"mime/multipart"
	"github.com/mshockwave/share-sound-api-server/storage"
	"strings"
	"io"
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
	}else if(req.Method == "POST" || req.Method == "PUT"){
		editProfile(resp, req)
	}else{
		//Not supported yet
		common.ResponseStatusAsJson(resp, 404, nil)
	}
}

func editProfile(resp http.ResponseWriter, req *http.Request){

	var userId string
	if id, err := GetSessionUserId(req); err != nil {
		common.LogE.Printf("Error fetching user sessoin id: %s\n", err)
		common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
			Message: "Error",
		})
		return
	}else{
		userId = id
	}

	if err := req.ParseForm();err != nil {
		common.LogE.Printf("Edit profile failed: %s\n", err.Error())
		common.ResponseStatusAsJson(resp, 400, &common.SimpleResult{
			Message: "Error",
			Description: "Wrong Edit Format",
		})
		return
	}

	var newPwd string = req.FormValue("password")
	var thumbnail multipart.File = nil
	var thumbnailHeader *multipart.FileHeader = nil

	thumbnail, thumbnailHeader, _ = req.FormFile("thumbnail")

	//Change Thumbnail
	var thumbObjName string = ""
	if thumbnail != nil {
		if client, err := storage.GetNewStorageClient(); err == nil{
			defer client.Close()

			bucket := client.GetDefaultBucket()
			thumbObjName = common.STORAGE_THUMBNAIL_FOLDER
			if thumbnailHeader == nil {
				//No file extension
				thumbObjName = common.PathJoin(thumbObjName, common.GetDefaultSecureHash())
			}else{
				//Keep file extension
				segs := strings.Split(thumbnailHeader.Filename, ".")
				if len(segs) > 1 {
					thumbObjName = common.PathJoin(thumbObjName,
									common.GetDefaultSecureHash() + "." + segs[len(segs) - 1])
				}else{
					thumbObjName = common.PathJoin(thumbObjName, common.GetDefaultSecureHash())
				}
			}

			obj := bucket.Object(thumbObjName)
			w := obj.NewWriter(client.Ctx)
			if _, e := io.Copy(w, thumbnail); e != nil {
				common.LogE.Printf("Error writing thumbnail file: %s\n", e)
				thumbObjName = ""
			}
			w.Close()
		}
	}

	//Update db entries
	if client, err := db.GetNewDataStoreClient(); err == nil {
		user := dbSchema.User{}
		_, e := client.RunInTransaction(func(tx *datastore.Transaction) error{
			key := client.NewKey(dbSchema.USER_PROFILE_KIND, userId, 0, UserProfileRootKey)

			if err := tx.Get(key, &user); err != nil {
				return err
			}

			//Update password
			if len(newPwd) > 0{
				if resBytes, e := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost); e == nil{
					user.Auth.PasswordBcryptHash = string(resBytes)
					//user.Auth.PasswordBcryptCost = bcrypt.DefaultCost
				}
			}

			//Update thumbnail
			if len(thumbObjName) > 0 {
				user.Thumbnail = thumbObjName
			}

			if _, err := tx.Put(key, &user); err != nil {
				return err
			}

			return nil
		})

		if e != nil {
			common.LogE.Printf("Update user profile of %s failed: %s\n", userId, e)
			common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
				Message: "Error",
			})

		}else{
			resProfile := &resultUserProfile{}
			resProfile.FromDbSchema(&user)
			common.ResponseOkAsJson(resp, resProfile)
		}
	}else{
		common.LogE.Printf("Get datastore failed: %s\n", err)
		common.ResponseStatusAsJson(resp, 500, &common.SimpleResult{
			Message: "Error",
		})
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