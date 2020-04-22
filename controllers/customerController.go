package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/KarthickSCode/customerService/repository"
	"github.com/KarthickSCode/customerService/utils"
	"github.com/erply/api-go-wrapper/pkg/api"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type CustomerController struct {
	config      *utils.Config
	customerDao *repository.CustomerDao
	idGenDao    *repository.IdGenDao
}

type saveCustomerResponse struct {
	Id            int `json:"customerID"`
	AlreadyExists int `json:"alreadyExists"`
}

func NewCustomerController(config *utils.Config, customerDao *repository.CustomerDao,
	idGenDao *repository.IdGenDao) *CustomerController {
	return &CustomerController{
		config:      config,
		customerDao: customerDao,
		idGenDao:    idGenDao,
	}
}

func (uc *CustomerController) parseRequest(r *http.Request) (api.Customer, error) {
	dec := json.NewDecoder(r.Body)
	var customer api.Customer
	err := dec.Decode(&customer)

	if err != nil {
		fmt.Println("Decode Error:", err.Error())
		return api.Customer{}, err
	}

	return customer, nil
}

func (controller *CustomerController) SaveCustomer(w http.ResponseWriter, r *http.Request) {

	var alreadyExists = 0
	customer, err := controller.parseRequest(r)

	if checkErrorExist(err, w) {
		return
	}

	if customer.CustomerID == 0 {
		nextId, err := controller.idGenDao.NextId(controller.config.IdGenerationKey.CustomerKey)
		if checkErrorExist(err, w) {
			return
		}
		customer.CustomerID = nextId
		controller.customerDao.SaveCustomer(&customer)
	} else {
		alreadyExists = 1
		controller.customerDao.UpdateCustomer(&customer)
	}

	var customerId = customer.CustomerID

	var response = saveCustomerResponse{
		Id:            customerId,
		AlreadyExists: alreadyExists,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (controller *CustomerController) GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if checkErrorExist(err, w) {
		return
	}
	obj, err := controller.customerDao.FindCustomer(id)
	if checkErrorExist(err, w) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(obj)
}

func checkErrorExist(err error, w http.ResponseWriter) bool {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.ErrorResponse{Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return true
	}
	return false
}
