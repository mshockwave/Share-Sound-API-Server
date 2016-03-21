package main

import(
	"net/http"
	"fmt"

	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/context"
	"github.com/mshockwave/share-sound-api-server/common"
)

func setUpCORS() goHandlers.CORSOption{
	origins := make([]string, 1)
	origins[0] = "*"

	return goHandlers.AllowedOrigins(origins)
}

func main() {

	router := SetUpRouters()
	allowOrigins := setUpCORS()

	http.Handle("/", router)


	addrStr := fmt.Sprintf("%s:%d",
		common.Config.GetString("server.address"),
		common.Config.GetInt("server.port"))
	common.LogV.Printf("Listen address: %s\n", addrStr)
	common.LogE.Fatal(http.ListenAndServe(
		addrStr,
		context.ClearHandler(goHandlers.CORS(allowOrigins)(http.DefaultServeMux)),
	))

}
