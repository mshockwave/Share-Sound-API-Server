package common

import(
	"log"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
	"github.com/gorilla/sessions"
	"github.com/wendal/errors"
)

func init(){
	initConfig()

	initLoggers()

	initSession()
}

//Configuration
func setDefaultValues(){
	Config.SetDefault("log.filePath", "")
	Config.SetDefault("log.enableStdOut", false)
	Config.SetDefault("log.enableStdErr", false)

	Config.SetDefault("server.address", "")
	Config.SetDefault("server.port", 8888)
}
func initConfig(){
	Config = viper.New()
	Config.SetConfigName(CONFIG_FILE_NAME)
	Config.AddConfigPath(".")
	Config.AddConfigPath("..")

	setDefaultValues()

	if e := Config.ReadInConfig(); e != nil {
		log.Fatalln("Error reading config file: " + e.Error())
		panic(e)
	}
}

//Loggers
func initLoggers() {
	var writer io.Writer = ioutil.Discard
	var errWriter io.Writer = ioutil.Discard

	if Config.GetBool("log.enableStdOut") {
		writer = io.MultiWriter(writer, os.Stdout)
	}
	if Config.GetBool("log.enableStdErr") {
		errWriter = io.MultiWriter(errWriter, os.Stderr)
	}

	logFilePath := Config.GetString("log.filePath")
	if len(logFilePath) > 0 {
		if file, err := os.Open(logFilePath); err == nil {
			writer = io.MultiWriter(writer, file)
			errWriter = io.MultiWriter(errWriter, file)
		}
	}

	LogV = log.New(writer, "[VERBOSE]:", log.Ldate | log.Ltime | log.Lshortfile)
	LogD = log.New(writer, "[DEBUG]:", log.Ldate | log.Ltime | log.Lshortfile)
	LogE = log.New(errWriter, "[ERROR]:", log.Ldate | log.Ltime | log.Lshortfile)
	LogW = log.New(errWriter, "[WARNING]:", log.Ldate | log.Ltime | log.Lshortfile)

	//fmt.Printf("Log enable stdout: %v\n", Config.GetBool("log.enableStdOut"))
	//fmt.Printf("Log enable stderr: %v\n", Config.GetBool("log.enableStdErr"))
}

//Sessions
func initSession(){
	if !Config.IsSet("session.privateKey") {
		panic(errors.New("No private key for session storage"))
	}

	SessionStorage = sessions.NewCookieStore([]byte(Config.GetString("session.privateKey")))
	SessionStorage.MaxAge(86400 * 3) //3 days
}
