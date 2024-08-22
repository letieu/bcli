package cmd

import (
	"github.com/letieu/bcli/api"
	"github.com/letieu/bcli/view"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication",
	Long:  `Authentication`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login",
	Long:  `Login to blueprint`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Logging in...")

		file, _ := cmd.Flags().GetString("file")

        var username, password string
        var err error

		if file != "" {
            username, password, err = getUserPassFromFile(file)
		} else {
			username, _ = cmd.Flags().GetString("username")
			password, _ = cmd.Flags().GetString("password")
		}

		err = api.Login(strings.TrimSpace(username), strings.TrimSpace(password))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Login successful")
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout",
	Long:  `Logout from blueprint`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Logging out...")

		err := api.Logout()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Logout successful")
	},
}

var infoCmd = &cobra.Command{
    Use: "info",
    Short: "Info",
    Long: "Get blueprint user info",
    Run: func(cmd *cobra.Command, args []string) {
        userInfo, err := api.GetUserInfo()

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        view.PrintUser(&userInfo)
    },
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(infoCmd)

	loginCmd.Flags().StringP("username", "u", "", "Username")
	loginCmd.Flags().StringP("password", "p", "", "Password")
	loginCmd.Flags().StringP("file", "f", "", "File containing username and password")
    loginCmd.MarkFlagFilename("file")

	loginCmd.MarkFlagsRequiredTogether("username", "password")
	loginCmd.MarkFlagsOneRequired("username", "file")
}

func getUserPassFromFile(file string) (string, string, error) {
	f, err := os.Open(file)

	if err != nil {
		return "", "", err
	}

	defer f.Close()
	var username, password string

	fmt.Fscanf(f, "%s\n%s", &username, &password)
	return username, password, nil
}
