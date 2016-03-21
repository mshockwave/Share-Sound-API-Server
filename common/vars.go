package common

import (
	"log"

	"github.com/spf13/viper"
	"github.com/gorilla/sessions"
)

const(
	USER_AUTH_SESSION = "user-auth"
	USER_ID_SESSION_KEY = "user_id"

	MAIN_STORAGE_BUCKET = "org-sharesound-main-bucket"
)

var(
	//Configurations
	CONFIG_FILE_NAME string = "config"
	Config	*viper.Viper

	//Loggers
	LogV	*log.Logger
	LogD	*log.Logger
	LogE	*log.Logger
	LogW	*log.Logger

	//Session
	SessionStorage *sessions.CookieStore
)
