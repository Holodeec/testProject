package runner

import (
	"log"
	"net/http"
	"testProject/internal/api"
	"testProject/internal/config"
	"testProject/internal/repository"
	"testProject/internal/router"
	"testProject/internal/service"

	"github.com/gorilla/mux"
)

func Run() {
	configuration := config.Get()
	connect, err := config.Connect(configuration)
	if err != nil {
		log.Fatal(err)
	}
	defer connect.Close()

	documentRepository := repository.NewDocumentRepository(connect.DB)
	documentService := service.NewDocumentService(documentRepository, configuration.Cache.Ttl)
	documentApi := api.NewDocumentApi(documentService)

	serveMux := router.Create(documentApi)
	err = start(configuration, serveMux)
	if err != nil {
		log.Fatal(err)
	}
}

func start(conf *config.Configuration, router *mux.Router) error {
	log.Println("Server started on port " + conf.Server.Port)
	return http.ListenAndServe(conf.Server.GetAddr(), router)
}
