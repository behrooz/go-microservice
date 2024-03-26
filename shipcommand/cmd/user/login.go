/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	username string
	password string
	host     string
)

type loginobj struct {
	Username string
	Password string
}

func login(username, password, host string) {
	url := "http://" + host + "/login"

	login := loginobj{
		Username: username,
		Password: password,
	}

	jsonBody, err := json.Marshal(login)
	if err != nil {
		fmt.Println("Error marshalling json", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading reponse body", err)
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("login successful", string(body))
		saveToken(string(body))
	} else {
		fmt.Println("login failed with status code:", resp.StatusCode)
	}

}

func saveToken(token string) {
	homedir, _ := os.UserHomeDir()
	shipDir := filepath.Join(homedir, ".ship")
	os.MkdirAll(shipDir, 0755)

	content := []byte(token)
	err := os.WriteFile(shipDir+"/"+"token", content, 0644)
	if err != nil {
		fmt.Println("Error with saving token", err)
	}
}

// loginCmd represents the login command
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to server with username and apssword",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		login(username, password, host)
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	LoginCmd.Flags().StringVarP(&username, "username", "u", "<username>", "<username>")
	LoginCmd.Flags().StringVarP(&password, "password", "p", "<password>", "<password>")
	LoginCmd.Flags().StringVarP(&host, "host", "s", "<host>", "<host>")

	if err := LoginCmd.MarkFlagRequired("username"); err != nil {
		fmt.Println(err.Error())
	}

	if err := LoginCmd.MarkFlagRequired("password"); err != nil {
		fmt.Println(err.Error())
	}

	if err := LoginCmd.MarkFlagRequired("host"); err != nil {
		fmt.Println(err.Error())
	}

}
