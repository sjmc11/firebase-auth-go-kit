package controllers

import (
	"errors"
	"firebase-sso/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strings"
	"time"
)

type ApiController struct{}

/*************/
/**  USERS  **/
/*************/

// GetContextUser
// RETURN USER DATA OF REQUESTER
func GetContextUser(g *gin.Context) (models.User, error) {
	ctxUserData, getUserOk := g.Get("dbUserData")
	if !getUserOk {
		return models.User{}, errors.New("could not check user data")
	}

	actionUserData, getUserOk := ctxUserData.(models.User)
	if !getUserOk {
		return actionUserData, errors.New("could not cast user data")
	}

	return actionUserData, nil
}

// HandleApiUserSync
// SYNC FIREBASE USER
func (u *ApiController) HandleApiUserSync(g *gin.Context) {

	actionUserData, getUserOk := GetContextUser(g)
	if getUserOk != nil {
		g.AbortWithStatusJSON(500, gin.H{"message": getUserOk.Error()})
		return
	}

	g.JSON(201, gin.H{"user": actionUserData.FormatDates, "message": "User UUID synchronized"})
	return
}

// HandleApiUserDisable
// DISABLE USER FROM SSO REGISTRY
func (u *ApiController) HandleApiUserDisable(g *gin.Context) {

	actionUserData, getUserOk := GetContextUser(g)
	if getUserOk != nil {
		g.AbortWithStatusJSON(500, gin.H{"message": getUserOk.Error()})
		return
	}

	var disableReqField struct {
		UserID int `json:"user_id"`
	}
	bindErr := g.ShouldBindWith(&disableReqField, binding.JSON)

	if bindErr != nil {
		g.AbortWithStatusJSON(500, gin.H{"data": bindErr.Error(), "message": "Invalid data"})
		return
	}

	// Find the user being disabled
	disableUser, lookupErr := models.GetUserById(disableReqField.UserID)
	if lookupErr != nil {
		if lookupErr.Error() == "no rows in result set" {
			g.AbortWithStatusJSON(500, gin.H{"message": "Invalid user_id"})
		} else {
			g.AbortWithStatusJSON(500, gin.H{"message": lookupErr.Error()})
		}
		return
	}

	// Check if already disabled
	if disableUser.DeletedAt.Valid {
		g.AbortWithStatusJSON(500, gin.H{"message": "User is already disabled"})
		return
	}

	// Check if deleting someone else
	// Must be admin: Role 1
	if disableUser.UUID != actionUserData.UUID && actionUserData.Role != 1 {
		g.AbortWithStatusJSON(500, gin.H{"message": "Only admin can disable other users"})
		return
	}

	disableUser.PublicDeletedAt = time.Now().Unix()
	disableUser.PublicUpdatedAt = time.Now().Unix()

	disableErr := disableUser.DisableAccount()
	if disableErr != nil {
		fmt.Println(disableErr.Error())
		g.AbortWithStatusJSON(500, gin.H{"message": disableErr.Error()})
		return
	}

	g.JSON(201, gin.H{"user": disableUser.FormatDates(), "message": "User account disabled"})
	return
}

// HandleApiUserEnable
// RE-ENABLE USER
func (u *ApiController) HandleApiUserEnable(g *gin.Context) {

	var enableReqField struct {
		UserID int `json:"user_id"`
	}
	bindErr := g.ShouldBindWith(&enableReqField, binding.JSON)

	if bindErr != nil {
		g.AbortWithStatusJSON(500, gin.H{"data": bindErr.Error(), "message": "Invalid data"})
		return
	}

	// Find the user being enabled
	enabledUser, lookupErr := models.GetUserById(enableReqField.UserID)
	if lookupErr != nil {
		if lookupErr.Error() == "no rows in result set" {
			g.AbortWithStatusJSON(500, gin.H{"message": "Invalid user_id"})
		} else {
			g.AbortWithStatusJSON(500, gin.H{"message": lookupErr.Error()})
		}
		return
	}

	// Check if already disabled
	if !enabledUser.DeletedAt.Valid {
		g.AbortWithStatusJSON(500, gin.H{"message": "User is already enabled"})
		return
	}

	enabledUser.PublicUpdatedAt = time.Now().Unix()

	err := enabledUser.EnableAccount()
	if err != nil {
		g.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	g.JSON(201, gin.H{"user": enabledUser.FormatDates(), "message": "User account enabled"})
	return
}

// HandleApiCreateUser
// Create new user
func (u *ApiController) HandleApiCreateUser(g *gin.Context) {

	// No need to check requesting user. AdminOnly rule in place for this route

	var newUser models.UserCreateData
	bindErr := g.ShouldBindWith(&newUser, binding.JSON)

	if bindErr != nil {
		g.AbortWithStatusJSON(500, gin.H{"data": bindErr.Error(), "message": "Invalid data"})
		return
	}

	// Set created at unix
	newUser.CreatedAt = time.Now().Unix()
	// Set default role
	if newUser.Role == 0 {
		newUser.Role = 2
	}

	newUserID, updateErr := newUser.RegisterUser()
	if updateErr != nil {
		if strings.Contains(updateErr.Error(), "unique constraint") {
			g.AbortWithStatusJSON(500, gin.H{"message": "Email address already in use"})
		} else {
			g.AbortWithStatusJSON(500, gin.H{"message": "Something went wrong registering user", "data": updateErr.Error()})
		}
		return
	}

	g.JSON(200, gin.H{"status": "OK", "message": "User created", "user_id": newUserID})
	return
}

// HandleApiUpdateUser
// UPDATE USER INFORMATION
func (u *ApiController) HandleApiUpdateUser(g *gin.Context) {

	actionUserData, getUserOk := GetContextUser(g)
	if getUserOk != nil {
		g.AbortWithStatusJSON(500, gin.H{"message": getUserOk.Error()})
		return
	}

	// Bind post data to request
	var UserEditData models.UserEditData
	bindErr := g.ShouldBindWith(&UserEditData, binding.JSON)

	if bindErr != nil {
		g.AbortWithStatusJSON(500, gin.H{"data": bindErr.Error(), "message": "Invalid data"})
		return
	}

	// Fallback if no role is supplied
	if UserEditData.Role == 0 {
		UserEditData.Role = actionUserData.Role
	}

	// Fallback if user_id is supplied
	if UserEditData.ID == 0 {
		UserEditData.ID = actionUserData.ID
		UserEditData.UUID = actionUserData.UUID
	}

	// Check user is editing their own profile or admin.
	if (actionUserData.ID != UserEditData.ID) && actionUserData.Role != 1 {
		g.AbortWithStatusJSON(500, gin.H{"message": "Not permitted"})
		return
	}

	// Prevent user making themselves an admin
	if (actionUserData.ID == UserEditData.ID) && actionUserData.Role != 1 && UserEditData.Role == 1 {
		g.AbortWithStatusJSON(500, gin.H{"message": "You cannot make yourself an admin"})
		return
	}

	// If editing another user check they exist
	if actionUserData.ID != UserEditData.ID {
		existingUser, lookupErr := models.GetUserById(UserEditData.ID)
		if lookupErr != nil {
			if lookupErr.Error() == "no rows in result set" {
				g.AbortWithStatusJSON(500, gin.H{"message": "Invalid user_id"})
			} else {
				g.AbortWithStatusJSON(500, gin.H{"message": lookupErr.Error()})
			}
			return
		}
		// User found, apply their UUID to editData for cache clear
		UserEditData.UUID = existingUser.UUID
	}

	updateErr := UserEditData.UpdateUser()
	if updateErr != nil {
		if strings.Contains(updateErr.Error(), "unique constraint") {
			fmt.Println(updateErr.Error())
			g.AbortWithStatusJSON(500, gin.H{"message": "Email address already in use"})
		} else {
			g.AbortWithStatusJSON(500, gin.H{"message": "Something went wrong updating user", "data": updateErr.Error()})
		}
		return
	}

	g.JSON(200, gin.H{"status": "OK", "message": "User information updated"})
	return
}

// HandleApiListUsers
// UPDATE USER INFORMATION
func (u *ApiController) HandleApiListUsers(g *gin.Context) {

	userList, listErr := models.GetUsers()
	if listErr != nil {
		g.AbortWithStatusJSON(500, gin.H{"message": "Something went fetching users", "data": listErr.Error()})
		return
	}

	g.JSON(200, gin.H{"status": "OK", "users": userList})
	return
}

/****************/
/**  PROJECTS  **/
/****************/

// HandleApiCreateProject
// CREATE A PROJECT
func (u *ApiController) HandleApiCreateProject(g *gin.Context) {

	actionUserData, getUserOk := GetContextUser(g)
	if getUserOk != nil {
		g.AbortWithStatusJSON(500, gin.H{"message": getUserOk.Error()})
		return
	}

	var newProject models.Project
	bindErr := g.ShouldBindWith(&newProject, binding.JSON)

	if bindErr != nil {
		g.AbortWithStatusJSON(500, gin.H{"data": bindErr.Error(), "message": "Invalid data"})
		return
	}

	newProject.CreatedBy = actionUserData.ID
	newProject.Created = time.Now().Unix()

	_, createProjErr := newProject.Create()
	if createProjErr != nil {
		fmt.Println(createProjErr)
		g.AbortWithStatusJSON(500, gin.H{"data": createProjErr.Error(), "message": "Something went wrong creating project"})
		return
	}

	g.JSON(200, gin.H{"status": "OK", "message": "Project created"})
	return
}

// HandleApiListProjects
// FETCH A USERS PROJECTS
func (u *ApiController) HandleApiListProjects(g *gin.Context) {

	actionUserData, getUserOk := GetContextUser(g)
	if getUserOk != nil {
		g.AbortWithStatusJSON(500, gin.H{"message": getUserOk.Error()})
		return
	}

	var listProjFields struct {
		UserID int `json:"user_id"`
	}
	bindErr := g.ShouldBindWith(&listProjFields, binding.JSON)

	if bindErr != nil {
		if bindErr.Error() != "EOF" {
			g.AbortWithStatusJSON(500, gin.H{"data": bindErr.Error(), "message": "Invalid data"})
			return
		}
	}

	// Apply user ID filter if admin
	if listProjFields.UserID > 0 && actionUserData.Role == 1 {
		actionUserData.ID = listProjFields.UserID
	}

	// Get projects
	userProjects, projErr := actionUserData.GetProjects()
	if projErr != nil {
		g.AbortWithStatusJSON(500, gin.H{"message": "Something went wrong fetching projects", "data": projErr.Error()})
		return
	}

	g.JSON(200, gin.H{"status": "OK", "projects": userProjects})
	return
}

// HandleApiDeleteProject
// DELETE A PROJECT
func (u *ApiController) HandleApiDeleteProject(g *gin.Context) {

	actionUserData, getUserOk := GetContextUser(g)
	if getUserOk != nil {
		g.AbortWithStatusJSON(500, gin.H{"message": getUserOk.Error()})
		return
	}

	var DelProjFields struct {
		ProjectID int `json:"project_id" binding:"required"`
	}
	bindErr := g.ShouldBindWith(&DelProjFields, binding.JSON)

	if bindErr != nil {
		g.AbortWithStatusJSON(500, gin.H{"data": bindErr.Error(), "message": "Invalid data"})
		return
	}

	// Fetch project data by ID
	project, projErr := models.FetchProjectByID(DelProjFields.ProjectID)
	if projErr != nil {
		if projErr.Error() == "no rows in result set" {
			g.AbortWithStatusJSON(500, gin.H{"message": "Invalid project_id", "data": projErr.Error()})
		} else {
			g.AbortWithStatusJSON(500, gin.H{"message": "Could not verify project data", "data": projErr.Error()})
		}
		return
	}

	// Check user is deleting their own project or requesting user is admin
	if project.CreatedBy == actionUserData.ID || actionUserData.Role == 1 {
		err := project.Delete()
		if err != nil {
			g.AbortWithStatusJSON(500, gin.H{"message": "Something went wrong deleting project"})
			return
		}
	} else {
		g.AbortWithStatusJSON(500, gin.H{"message": "Project does not belong to you"})
		return
	}

	g.JSON(200, gin.H{"status": "Project deleted"})
	return
}

/******************/
/**  TEST ROUTE  **/
/******************/

func (u *ApiController) ApiTestRoute(g *gin.Context) {
	actionUserData, getUserOk := GetContextUser(g)
	if getUserOk != nil {
		g.AbortWithStatusJSON(500, gin.H{"message": getUserOk.Error()})
		return
	}

	g.JSON(200, gin.H{"status": "OK", "user": actionUserData})
	return
}
