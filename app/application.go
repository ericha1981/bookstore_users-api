package app

import "github.com/gin-gonic/gin"

/*
	Entry point!!
	Redirect requests coming in here to controllers. First layer of our application
*/

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	router.Run(":8080")
}