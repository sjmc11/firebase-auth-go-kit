package cmd

import (
	"errors"
	"firebase-sso/core"
	"firebase-sso/models"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// createuserCmd represents the createuser command
var createuserCmd = &cobra.Command{
	Use:   "user:create",
	Short: "Create an admin user",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		// ******************* //
		//     Create user     //
		// ******************* //

		// USER EMAIL
		emailPrompt := promptui.Prompt{
			Label: "User email",
			Validate: func(input string) error {
				if len(input) < 4 {
					return errors.New("email should have at least 4 characters")
				}
				return nil
			},
		}

		userEmail, emailErr := emailPrompt.Run()
		if emailErr != nil {
			fmt.Println(emailErr.Error())
			os.Exit(1)
		}

		// ADMIN USER FIRST NAME
		fNamePrompt := promptui.Prompt{
			Label: "User first name",
			Validate: func(input string) error {
				if len(input) < 2 {
					return errors.New("first name should have at least 2 characters")
				}
				return nil
			},
		}

		userFname, fNameErr := fNamePrompt.Run()
		if fNameErr != nil {
			fmt.Println(fNameErr.Error())
			os.Exit(1)
		}

		// ADMIN USER LAST NAME
		lNamePrompt := promptui.Prompt{
			Label: "User last name",
			Validate: func(input string) error {
				if len(input) < 2 {
					return errors.New("last name should have at least 2 characters")
				}
				return nil
			},
		}

		userLname, lNameErr := lNamePrompt.Run()
		if lNameErr != nil {
			fmt.Println(lNameErr.Error())
			os.Exit(1)
		}

		var newUser = models.UserCreateData{
			FirstName: userFname,
			LastName:  userLname,
			Email:     userEmail,
			Role:      1, // 1 = admin, 2 = default user
			CreatedAt: time.Now().Unix(),
		}

		core.InitBackgroundServices()

		newUserID, err := newUser.RegisterUser()
		if err != nil {
			fmt.Println(err.Error())
			log.Fatal("Could not create user account")
		}

		fmt.Println(aurora.BgGreen("Created admin user: " + newUser.Email))

		// ****************************** //
		//   Create example project? Y/N  //
		// ****************************** //

		doProjectPrompt := promptui.Prompt{
			Label: "Create an example project? Y to continue. Anything else to skip.",
			Validate: func(input string) error {
				//if strings.ToLower(strings.TrimSpace(input)) == "y" {
				//	return errors.New("last name should have at least 2 characters")
				//}
				return nil
			},
		}
		doExampleProject, doProjErr := doProjectPrompt.Run()
		if doProjErr != nil {
			fmt.Println(doProjErr.Error())
			os.Exit(1)
		}
		if strings.ToLower(strings.TrimSpace(doExampleProject)) == "y" {
			// Create an example project

			// Project title
			projTitlePrompt := promptui.Prompt{
				Label: "Project title",
				Validate: func(input string) error {
					if len(input) < 2 {
						return errors.New("project title should have at least 3 characters")
					}
					return nil
				},
			}

			projTitle, projTitleErr := projTitlePrompt.Run()
			if projTitleErr != nil {
				fmt.Println(projTitleErr.Error())
				os.Exit(1)
			}

			var newProject = models.Project{
				Title:     projTitle,
				Slug:      slug.Make(projTitle),
				Created:   time.Now().Unix(),
				CreatedBy: newUserID,
			}

			newProjectID, err := newProject.Create()
			if err != nil {
				return
			}

			fmt.Println(aurora.BgGreen("Created new project ID: " + strconv.Itoa(newProjectID)))
		}

	},
}

func init() {
	rootCmd.AddCommand(createuserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createuserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createuserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
