/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"firebase-sso/core"
	"firebase-sso/server"
	"firebase-sso/server/middleware"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Lacuna SSO API",
	Long:  `Run the firebase integrated single sign on API to authenticate requests for Lacuna tools`,
	Run: func(cmd *cobra.Command, args []string) {
		core.InitBackgroundServices()
		// FIREBASE
		_, fbErr := middleware.FbClient.InitFirebaseClient()
		if fbErr != nil {
			fmt.Println(aurora.Red("Something went wrong initializing Firebase client"))
			fmt.Println(fbErr.Error())
			return
		}
		server.StartWebServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
