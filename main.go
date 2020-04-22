package main

import (
	"errors"
	"fmt"
	"github.com/KarthickSCode/customerService/controllers"
	_ "github.com/KarthickSCode/customerService/docs"
	"github.com/KarthickSCode/customerService/repository"
	"github.com/KarthickSCode/customerService/utils"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
)

var customerDao *repository.CustomerDao
var idDao *repository.IdGenDao

// @title ERPLY Customer API
// @version 1.0
// @description ERPLY Customer API management.
// @termsOfService http://swagger.io/terms/

// @contact.name Karthick Sivapragasam
// @contact.url https://github.com/KarthickSCode/customerService
// @contact.email karthicksivapragasam23@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationUrl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationUrl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information
func main() {

	cfg := utils.DefaultConfig()

	router := gin.New()

	initializeRouter(cfg, router)

	router.Run()
}

func initializeRouter(cfg *utils.Config, router *gin.Engine) {
	initMongoRepository(cfg)

	initRootRoute(router)

	initCustomerRoute(cfg, router)
}

func initRootRoute(router *gin.Engine) gin.IRoutes {
	return router.GET("/", func(context *gin.Context) {
		context.JSON(200, "Welcome to ERPLY API")
	})
}

func initCustomerRoute(cfg *utils.Config, router *gin.Engine) {
	customerController := controllers.NewCustomerController(cfg, customerDao, idDao)

	route := router.Group("/")
	{
		customers := route.Group("/customer")
		{
			customers.Use(auth())
			customers.GET(":id", customerController.GetCustomer)
			customers.POST("", customerController.SaveCustomer)
		}

	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

func initMongoRepository(cfg *utils.Config) {
	repo, err := repository.Setup(cfg.MongoAddressURI)
	if err != nil {
		panic(err)
	}
	customerDao = repo.GetCustomerDao(cfg.DbName, cfg.DbCollection.Customer)
	idDao = repo.GetIdGenDao(cfg.DbName, cfg.DbCollection.IdGen)
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.GetHeader("Authorization")) == 0 {
			fmt.Println("Header is missing")
			utils.NewError(c, http.StatusUnauthorized, errors.New("Authorization is required Header"))
			c.Abort()
		} else if c.GetHeader("Authorization") != "123-456-789" {
			utils.NewError(c, http.StatusUnauthorized, errors.New("Invalid Token"))
			c.Abort()
		}
		c.Next()
	}
}
