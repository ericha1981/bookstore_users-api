package app

import (
	"github.com/ericha1981/bookstore_users-api/logger"
	"github.com/gin-gonic/gin"
)

/*
	Entry point!!
	Redirect requests coming in here to controllers. First layer of our application
*/

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	logger.Info("about to start the application...")
	router.Run(":8080")
}