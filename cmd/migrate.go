/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"firebase-sso/core"
	"firebase-sso/models"
	"fmt"
	"github.com/logrusorgru/aurora"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Execute system tables migration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		core.InitBackgroundServices()

		// Create system schema
		dbErr := models.CreateSchema("system")
		if dbErr != nil {
			if dbErr.Error() == "schema already created" {
				fmt.Println(aurora.Green(dbErr.Error()))
			} else {
				fmt.Println(aurora.Red(dbErr.Error()))
			}
		}

		// Create users table
		dbErr = models.CreateUsersTable()
		if dbErr != nil {
			fmt.Println(aurora.Red(dbErr.Error()))
		}

		// Example table for demo purposes
		// Customize and copy this function to create tables for your app
		dbErr = models.CreateUserProjectsTable()
		if dbErr != nil {
			fmt.Println(aurora.Red(dbErr.Error()))
		}

		// Create a testing user
		var newUser = models.UserCreateData{
			FirstName: "Test",
			LastName:  "User",
			Email:     "development",
			Role:      1, // Admin
			CreatedAt: time.Now().Unix(),
		}
		newUserID, err := newUser.RegisterUser()
		if err != nil {
			if err.Error() == "email already in use" {
				fmt.Println("Test user exists already")
			} else {
				log.Fatal("Could not create user account")
			}
		} else {
			fmt.Println(aurora.BgGreen("Created test user ID: " + strconv.Itoa(newUserID) + " development token now supported."))
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
