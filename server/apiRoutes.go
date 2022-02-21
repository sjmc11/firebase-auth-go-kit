package server

import (
	"firebase-sso/controllers"
	"firebase-sso/server/middleware"
	"github.com/gin-gonic/gin"
)

func LoadApiRoutes(router *gin.Engine) {

	api := router.Group("/api/v1", middleware.FbClient.FirebaseAuthCheck())

	apiController := controllers.ApiController{}

	// Firebase User sync
	api.POST("/user/sync", middleware.FbClient.AppAuthCheck(false))

	// Create user
	// Accepts "UserCreateData" struct
	api.POST("/user/create", middleware.FbClient.AppAuthCheck(true), apiController.HandleApiCreateUser)

	// Update user information
	// Accepts "UserEditData" struct
	api.POST("/user/update", middleware.FbClient.AppAuthCheck(false), apiController.HandleApiUpdateUser)

	// Disable user account
	// Requires "user_id"
	api.POST("/user/disable", middleware.FbClient.AppAuthCheck(false), apiController.HandleApiUserDisable)

	// Re-enable user account
	// Requires "user_id"
	api.POST("/user/enable", middleware.FbClient.AppAuthCheck(true), apiController.HandleApiUserEnable)

	// List users
	// No params
	api.GET("/user/list", middleware.FbClient.AppAuthCheck(false), apiController.HandleApiListUsers)

	// Create project
	// Accepts "Project" struct
	api.POST("/projects/create", middleware.FbClient.AppAuthCheck(false), apiController.HandleApiCreateProject)

	// List user projects
	// Accepts "user_id" (optional)
	api.GET("/projects/list", middleware.FbClient.AppAuthCheck(false), apiController.HandleApiListProjects)

	// Delete project
	// Requires "project_id"
	api.POST("/projects/delete", middleware.FbClient.AppAuthCheck(false), apiController.HandleApiDeleteProject)

	// Testing route
	api.POST("/test", middleware.FbClient.AppAuthCheck(true), apiController.ApiTestRoute)
}
