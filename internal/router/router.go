package router

import (
	"log"
	"testProject/internal/api"

	"github.com/gorilla/mux"
)

const URL = "/api/v1"

func Create(documentApi *api.DocumentApi) *mux.Router {

	router := mux.NewRouter()
	log.Println("Create router")

	document := router.PathPrefix(URL + "/document").Subrouter()
	document.HandleFunc("/all", documentApi.FindAll()).Methods("GET")
	document.HandleFunc("/{id}", documentApi.FindByID()).Methods("GET")
	document.HandleFunc("/delete/{id}", documentApi.Delete()).Methods("DELETE")
	document.HandleFunc("/create", documentApi.Create()).Methods("POST")
	document.HandleFunc("/update", documentApi.Update()).Methods("PUT")

	return router
}
