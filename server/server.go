package server

import (
	"firebase-sso/helpers/env"
	"firebase-sso/server/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/logrusorgru/aurora"
	"log"
)

var Router *gin.Engine

func StartWebServer() {
	fmt.Println(aurora.Bold("*** STARTING API WEB SERVER ***").BgGreen())

	Router = gin.Default()

	Router.Use(middleware.CORSMiddleware())

	// Return 200 for health checks on base URL
	Router.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{"message": "OK"})
	})

	LoadApiRoutes(Router)

	serveErr := Router.Run(":" + env.Get("PORT"))
	if serveErr != nil {
		log.Fatal(serveErr)
	}
}
