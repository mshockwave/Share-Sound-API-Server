package main

import(
	"github.com/gorilla/mux"
	"github.com/mshockwave/share-sound-api-server/handlers"
)

func SetUpRouters() *mux.Router{
	router := mux.NewRouter()

	handlers.ConfigureUserHandlers(router.PathPrefix("/user").Subrouter())

	return router
}
