package main

import (
	"fmt"
	"github.com/KarthickSCode/customerService/controllers"
	"github.com/KarthickSCode/customerService/repository"
	"github.com/KarthickSCode/customerService/utils"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

var customerDao *repository.CustomerDao
var idDao *repository.IdGenDao

func main() {

	cfg := utils.DefaultConfig()

	router := mux.NewRouter().StrictSlash(true)

	initializeRouter(cfg, router)

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
}

func initializeRouter(cfg *utils.Config, router *mux.Router) {
	initMongoRepository(cfg)

	initRootRoute(router)

	initCustomerRoute(cfg, router)
}

func initCustomerRoute(cfg *utils.Config, router *mux.Router) {
	customerController := controllers.NewCustomerController(cfg, customerDao, idDao)
	router.Methods(http.MethodGet).Path("/customer/{id}").Name("GetCustomer").HandlerFunc(customerController.GetCustomer)
	router.Methods(http.MethodPost).Path("/customer").Name("SaveCustomer").HandlerFunc(customerController.SaveCustomer)

	fs := http.FileServer(http.Dir("./swaggerui/"))
	router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))
}

func initRootRoute(router *mux.Router) *mux.Route {
	return router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Welcome to the ERPLY customer API")
	})
}

func initMongoRepository(cfg *utils.Config) {
	repo, err := repository.Setup(cfg.MongoAddressURI)
	if err != nil {
		panic(err)
	}
	customerDao = repo.GetCustomerDao(cfg.DbName, cfg.DbCollection.Customer)
	idDao = repo.GetIdGenDao(cfg.DbName, cfg.DbCollection.IdGen)
}
