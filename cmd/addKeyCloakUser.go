/*
Copyright © 2023 Dawid Skraba <dawid.skraba@ucdconnect.ie>
*/
package cmd

import (
	"bytes"
	"cliapp/auth"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const (
	keycloakAdminURL = "http://localhost:8080/auth/admin/realms/my_realm"
	keycloakAuthURL  = "http://localhost:8080/auth/realms/master/protocol/openid-connect/token"
	client_secret    = "admin"
)

var (
	keyUser         string
	adminPassword   string
	adminUsername   string
	newUserUsername string
	newUserPassword string
)

// addKeyCloakUserCmd represents the addKeyCloakUser command
var addKeyCloakUserCmd = &cobra.Command{
	Use:   "addKeyCloakUser",
	Short: "Add a user to KeyCloak server",
	Long: `
	You can add a user to KeyCloak server using this command. To perform this command you
	neeed to authenticate yourself as an admin of the keycloak server. You can do this by providing your admin username
	and password. You also need to provide the username and password of the new user you want to add.

	Example of the addKeyCloakUser command:
		$ ./cliapp addKeyCloakUser -u=admin -p=password -s=greg -a=gregpass

		$ ./cliapp addKeyCloakUser -u=admin -p=password -s=greg -a=@password.txt

	`,
	Run: func(cmd *cobra.Command, args []string) {
		email := newUserUsername + "@gmail.com"
		user := KeycloakUser{
			Username: newUserUsername,
			Email:    email,
			Enabled:  true,
		}

		token, err := auth.GetAdminToken(adminUsername, adminPassword, "admin-cli", client_secret, keycloakAuthURL)
		if err != nil {
			log.Fatalf("Error getting Keycloak token: %v", err)
		}

		userID, err := createKeycloakUser(token, user)
		if err != nil {
			log.Fatalf("Error creating Keycloak user: %v", err)
		}

		password := KeycloakPassword{
			Value:     newUserPassword,
			Temporary: false,
			Type:      "password",
		}

		if strings.Contains(newUserPassword, "@") { // JSON file as input
			file := strings.Split(newUserPassword, "@")[1]
			passwordByte, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatalf("Error reading password from file: %v", err)
			}
			pass := string(passwordByte)
			password = KeycloakPassword{
				Value:     pass,
				Temporary: false,
				Type:      "password",
			}
		}

		err = setKeycloakUserPassword(token, userID, password)
		if err != nil {
			log.Fatalf("Error setting Keycloak user password: %v", err)
		}

		groupID, err := getKeycloakGroupIDByName(token, "vault-client")
		if err != nil {
			log.Fatalf("Error fetching Keycloak group ID: %v", err)
		}

		err = addUserToKeycloakGroup(token, userID, groupID)
		if err != nil {
			log.Fatalf("Error adding user to Keycloak group: %v", err)
		}

		fmt.Println("User created successfully")
	},
}

func init() {
	rootCmd.AddCommand(addKeyCloakUserCmd)

	// ADMIN parameters
	addKeyCloakUserCmd.Flags().StringVarP(&adminUsername, "adminUsername", "u", "", "Keycloak Admin Username")
	if err := addKeyCloakUserCmd.MarkFlagRequired("adminUsername"); err != nil {
		fmt.Println(err)
	}
	addKeyCloakUserCmd.Flags().StringVarP(&adminPassword, "adminPassword", "p", "", "Keycloak Admin Password")
	if err := addKeyCloakUserCmd.MarkFlagRequired("adminPassword"); err != nil {
		fmt.Println(err)
	}

	// NEW USER parameters
	addKeyCloakUserCmd.Flags().StringVarP(&newUserUsername, "newUserUsername", "s", "", "New User Username")
	if err := addKeyCloakUserCmd.MarkFlagRequired("newUserUsername"); err != nil {
		fmt.Println(err)
	}
	addKeyCloakUserCmd.Flags().StringVarP(&newUserPassword, "newUserPassword", "a", "", "New User Password")
	if err := addKeyCloakUserCmd.MarkFlagRequired("newUserPassword"); err != nil {
		fmt.Println(err)
	}

}

type KeycloakUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Enabled  bool   `json:"enabled"`
}

type KeycloakPassword struct {
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
	Type      string `json:"type"`
}

type KeycloakGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func createKeycloakUser(token *oauth2.Token, user KeycloakUser) (string, error) {
	client := &http.Client{}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user JSON: %w", err)
	}

	req, err := http.NewRequest("POST", keycloakAdminURL+"/users", bytes.NewBuffer(userJSON))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create user, status: %d, response: %s", resp.StatusCode, string(body))
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("location header not found in the response")
	}

	userID := location[strings.LastIndex(location, "/")+1:]
	return userID, nil
}

func setKeycloakUserPassword(token *oauth2.Token, userID string, password KeycloakPassword) error {
	client := &http.Client{}

	passwordJSON, err := json.Marshal(password)
	if err != nil {
		return fmt.Errorf("failed to marshal password JSON: %w", err)
	}

	req, err := http.NewRequest("PUT", keycloakAdminURL+"/users/"+userID+"/reset-password", bytes.NewBuffer(passwordJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to set password, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

func getKeycloakGroupIDByName(token *oauth2.Token, groupName string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", keycloakAdminURL+"/groups?search="+groupName, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch groups: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to fetch groups, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var groups []KeycloakGroup
	if err := json.Unmarshal(body, &groups); err != nil {
		return "", fmt.Errorf("failed to unmarshal groups JSON: %w", err)
	}

	if len(groups) == 0 {
		return "", fmt.Errorf("group not found: %s", groupName)
	}

	return groups[0].ID, nil
}

func addUserToKeycloakGroup(token *oauth2.Token, userID, groupID string) error {
	client := &http.Client{}

	req, err := http.NewRequest("PUT", keycloakAdminURL+"/users/"+userID+"/groups/"+groupID, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to add user to group, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}
