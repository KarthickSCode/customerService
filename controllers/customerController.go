package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/KarthickSCode/customerService/repository"
	"github.com/KarthickSCode/customerService/utils"
	"github.com/erply/api-go-wrapper/pkg/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CustomerController struct {
	config      *utils.Config
	customerDao *repository.CustomerDao
	idGenDao    *repository.IdGenDao
}

type SaveCustomerResponse struct {
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

// SaveCustomer godoc
// @Summary Adding a new Customer
// @Description Add customer details by json
// @Tags customer
// @Accept  json
// @Produce  json
// @Param customer body api.Customer true "Add Customer"
// @Success 200 {object} controllers.SaveCustomerResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /customer [post]
func (controller *CustomerController) SaveCustomer(ctx *gin.Context) {

	var alreadyExists = 0

	var customer api.Customer
	err := ctx.ShouldBindJSON(&customer)

	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if customer.CustomerID == 0 {
		nextId, err := controller.idGenDao.NextId(controller.config.IdGenerationKey.CustomerKey)
		if err != nil {
			utils.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		customer.CustomerID = nextId
		controller.customerDao.SaveCustomer(&customer)
	} else {
		alreadyExists = 1
		controller.customerDao.UpdateCustomer(&customer)
	}

	var customerId = customer.CustomerID

	var response = SaveCustomerResponse{
		Id:            customerId,
		AlreadyExists: alreadyExists,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetCustomer godoc
// @Summary Get the customer details
// @Description get customer by ID
// @Tags customer
// @Accept  json
// @Produce  json
// @Param id path int true "Customer ID"
// @Success 200 {object} api.Customer
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /customer/{id} [get]
func (controller *CustomerController) GetCustomer(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	obj, err := controller.customerDao.FindCustomer(id)
	if err != nil {
		utils.NewError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, obj)
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
