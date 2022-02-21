package middleware

import (
	"context"
	"errors"
	"firebase-sso/helpers/env"
	"firebase-sso/models"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/logrusorgru/aurora"
	"google.golang.org/api/option"
	"strings"
)

// Authentication layer of middleware
// First, requests go through the firebase auth check
// Second, requests go through the App auth check
// Users must exist in the database already in order to authenticate
// If not they are removed from firebase Auth
// Finally, requests are passed to their controller logic to be handled

type FireBase struct {
	Client *auth.Client
}

var FbClient FireBase

func (fbClient *FireBase) InitFirebaseClient() (*auth.Client, error) {
	if env.Get("FIREBASE_KEY") == "" {
		return nil, errors.New("define FIREBASE_KEY json file location in .env")
	}
	// Setup firebase
	opt := option.WithCredentialsFile(env.Get("FIREBASE_KEY"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	FbClient.Client = authClient

	return authClient, nil
}

func (fbClient *FireBase) FirebaseAuthCheck() gin.HandlerFunc {
	return func(gContext *gin.Context) {
		// Get auth token from headers
		authToken := gContext.Request.Header.Get("Authorization")
		authToken = strings.TrimSpace(strings.Replace(authToken, "Bearer", "", 1))

		if authToken == "" || authToken == "null" {
			gContext.AbortWithStatusJSON(403, gin.H{"message": "No auth token detected"})
			return
		}

		// Allow "development" token for testing
		if authToken == "development" {
			gContext.Set("fbUserData", &auth.Token{
				AuthTime: 0,
				Issuer:   "",
				Audience: "",
				Expires:  0,
				IssuedAt: 0,
				Subject:  "",
				UID:      "development",
				Firebase: auth.FirebaseInfo{},
				Claims: map[string]interface{}{
					"email": "development",
				},
			})
			gContext.Next()
			return
		}

		// Verify with firebase
		fbUserData, err := fbClient.Client.VerifyIDToken(context.Background(), authToken)
		if err != nil {
			gContext.AbortWithStatusJSON(403, gin.H{"message": "OAuth token invalid", "error": err.Error()})
			return
		}

		// Store firebase user data in request context
		gContext.Set("fbUserData", fbUserData)

		gContext.Next()
		return
	}
}

func (fbClient *FireBase) AppAuthCheck(adminOnly bool) gin.HandlerFunc {
	return func(gContext *gin.Context) {
		// Pickup firebase user data from context
		ctxUserData, getFbDataOk := gContext.Get("fbUserData")
		if !getFbDataOk {
			gContext.AbortWithStatusJSON(403, gin.H{"message": "Could not fetch firebase user data"})
			return
		}

		// Cast firebase context to struct
		fbUserData, getUserOk := ctxUserData.(*auth.Token)
		if !getUserOk {
			gContext.AbortWithStatusJSON(500, gin.H{"message": "Could not cast firebase user data"})
			return
		}

		// Get user email
		fbUserEmail, emailOk := fbUserData.Claims["email"]
		if !emailOk {
			gContext.AbortWithStatusJSON(500, gin.H{"message": "Could not get email from firebase"})
			return
		}
		// Get custom user data & check exists
		dbUserData, dbErr := models.GetUserByEmail(fbUserEmail.(string))
		if dbErr != nil {
			if dbErr.Error() == "no rows in result set" {

				// User not registered in our app back-end
				// Undo the registration from firebase auth
				err := fbClient.Client.DeleteUser(context.Background(), fbUserData.UID)
				if err != nil {
					fmt.Println(aurora.Red("Error deleting user from firebase."))
					//fmt.Println(err.Error())
				}
				fmt.Println(aurora.Yellow("API request from unregistered user: " + fbUserEmail.(string) + " : " + gContext.FullPath()))
				gContext.AbortWithStatusJSON(403, gin.H{"message": "User not registered"})
			} else {
				gContext.AbortWithStatusJSON(500, gin.H{"message": "Re-sync lacuna account to authorize", "data": dbErr.Error()})
			}
			return
		}

		// Set UUID if empty
		if len(dbUserData.UUID) < 3 {
			dbUserData.UUID = fbUserData.UID
			uuidErr := dbUserData.SetUUID()
			if uuidErr != nil {
				gContext.AbortWithStatusJSON(500, gin.H{"message": "something went wrong authorizing user"})
				return
			}
		}

		// Check account active
		if dbUserData.DeletedAt.Valid {
			gContext.AbortWithStatusJSON(403, gin.H{"message": "Account disabled"})
			return
		}

		// Check admin only
		if adminOnly && dbUserData.Role != 1 {
			gContext.AbortWithStatusJSON(403, gin.H{"message": "Admin only"})
			return
		}

		// Handle user sync inside of middleware
		if gContext.FullPath() == "/api/v1/user/sync" {
			gContext.JSON(201, gin.H{"user": dbUserData, "message": "User UUID synchronized"})
			return
		}

		gContext.Set("fbUserData", nil)
		gContext.Set("dbUserData", dbUserData)

		fmt.Println(aurora.BgBlack("API request from: " + dbUserData.Email + " / " + gContext.FullPath()))

		gContext.Next()
	}
}
