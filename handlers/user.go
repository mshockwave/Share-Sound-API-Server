package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func handleRegister(resp http.ResponseWriter, req *http.Request){

}

func handleLogin(resp http.ResponseWriter, req *http.Request){

}

func ConfigureUserHandlers(router *mux.Router){
	router.HandleFunc("/register", handleRegister)
	router.HandleFunc("/login", AuthVerifierWrapper(handleLogin))
}