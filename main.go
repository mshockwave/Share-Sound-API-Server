package main

import(
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/mshockwave/share-sound-api-server/handlers"
	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/context"
	"github.com/mshockwave/share-sound-api-server/common"

	_ "github.com/mshockwave/share-sound-api-server/storage"
	_ "github.com/mshockwave/share-sound-api-server/datastore"
)

func setUpRouters() *mux.Router{
	router := mux.NewRouter()

	handlers.ConfigureUserHandlers(router.PathPrefix("/user").Subrouter())

	return router
}

func setUpCORS() goHandlers.CORSOption{
	origins := make([]string, 1)
	origins[0] = "*"

	return goHandlers.AllowedOrigins(origins)
}

func main() {

	router := setUpRouters()
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
