package controller

import (
	_ "SB/docs"
	"SB/internal/configs"
	"SB/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// @title Bank API
// @version 1.0
// @description This is a sample API for banking with Gin and Swagger.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token for user authentication

func RunServer() error {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", Ping)

	authG := router.Group("/auth")
	{
		authG.POST("/sign-up", createUserHandler)
		authG.POST("/sign-in", authenticateHandler)
	}

	userG := router.Group("/users", checkUserAuthentication)
	{
		userG.PATCH("", updateUserHandler)
		userG.DELETE("/:id", deleteUserHandler)
		userG.GET("/:id", getUserByIDHandler)
		userG.GET("/inactive", getInActiveUsersHandler)
		userG.POST("/restore", restoreUserHandler)
		userG.GET("/find", findUserByNameHandler)
	}

	accountG := router.Group("/accounts", checkUserAuthentication)
	{
		accountG.POST("", createAccountHandler)
		accountG.PATCH("", updateAccountHandler)
		accountG.DELETE("/:id", deleteAccountHandler)
		accountG.GET("/:id", getAccountByIDHandler)
		accountG.GET("/users/:id", getAccountByUserIDHandler)
		accountG.GET("/inactive", getInActiveAccountsHandler)
		accountG.GET("/currency", getAccountByCurrency)
		accountG.GET("/:id/balance", getAccountBalanceHandler)
	}

	transferG := router.Group("/transfers", checkUserAuthentication)
	{
		transferG.POST("", createTransferHandler)
	}

	if err := router.Run(configs.AppSettings.AppParams.PortRun); err != nil {
		logger.Error.Printf("[controller] RunServer():  Error during running HTTP server: %s", err.Error())
		return err
	}

	return nil
}

func Ping(c *gin.Context) {
	// @Summary Ping the server
	// @Description Check if the server is up and running
	// @Tags general
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router / [get]
	c.JSON(http.StatusOK, gin.H{
		"message": "Server is up and running",
	})
}
