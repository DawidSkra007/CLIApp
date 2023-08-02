/*
Copyright © 2023 Dawid Skraba <dawid.skraba@ucdconnect.ie>
*/
package cmd

import (
	"cliapp/auth"
	"cliapp/util"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	username  string
	password  string
	policyAdd string
	u10       string
	p10       string
)

// createUserCmd represents the createUser command
var createUserCmd = &cobra.Command{
	Use:   "createUser",
	Short: "Create a new user with a given username, password and policy",
	Long: `
	Create a new user with a given username, password and policy

	Example of the createUser command(Keycloak Authentication):
		$ ./cliapp createUser --username=User2 --password=pass --policy=user-policy

	To use Userpass Authentication:
		$ ./cliapp createUser --username=User2 --password=pass --policy=user-policy --user=username --pass=password
	
	To use a different instance:
		$ ./cliapp createUser --username=User2 --password=pass --policy=user-policy --user=username --pass=password --instance
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("user").Changed && cmd.Flag("pass").Changed {
			address := util.UpdateAddress(instance)
			auth.AuthenticateWithUserPass(u10, p10, address)
		} else {
			if cmd.Flag("instance").Changed {
				fmt.Println("Error: You must provide a username and password to use a different instance.")
				os.Exit(1)
			}
			address := util.UpdateAddress(false)
			auth.KeycloakAuth(address)
		}

		username = strings.ToLower(username)
		err := AddUserWithPolicy(username, password, policyAdd)
		if err != nil {
			log.Fatalf("unable to add user: %v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(createUserCmd)

	//username
	createUserCmd.Flags().StringVarP(&username, "username", "u", "", "Name of the user")

	if err := createUserCmd.MarkFlagRequired("username"); err != nil {
		fmt.Println(err)
	}

	//password
	createUserCmd.Flags().StringVarP(&password, "password", "w", "", "Password of the user")

	if err := createUserCmd.MarkFlagRequired("password"); err != nil {
		fmt.Println(err)
	}

	//policy
	createUserCmd.Flags().StringVarP(&policyAdd, "policy", "p", "", "policy for user")

	if err := createUserCmd.MarkFlagRequired("policy"); err != nil {
		fmt.Println(err)
	}

	// userpass
	createUserCmd.Flags().StringVarP(&u10, "user", "a", "", "Userpass username")
	createUserCmd.Flags().StringVarP(&p10, "pass", "s", "", "Userpass password")
	createUserCmd.MarkFlagsRequiredTogether("user", "pass")

	createUserCmd.Flags().BoolVarP(&instance, "instance", "i", false, "Use another Vault instance")

}

func AddUserWithPolicy(username, password, policy string) error {
	userExists, err := UserExists(username)
	if err != nil {
		return err
	}
	if userExists {
		return fmt.Errorf("user '%s' already exists", username)
	}

	policyExists, err := PolicyExists(policy)
	if err != nil {
		return err
	}
	if !policyExists {
		return fmt.Errorf("policy '%s' does not exist", policy)
	}

	data := map[string]interface{}{
		"password": password,
		"policies": policy,
	}

	_, err = auth.Client.Logical().Write("auth/userpass/users/"+username, data)
	if err != nil {
		return err
	}

	log.Printf("User '%s' added with policy '%s'", username, policy)
	return nil
}

func UserExists(username string) (bool, error) {
	secret, err := auth.Client.Logical().List("auth/userpass/users")
	if err != nil {
		return false, err
	}

	if secret != nil && secret.Data["keys"] != nil {
		users := secret.Data["keys"].([]interface{})
		for _, user := range users {
			if user == username {
				return true, nil
			}
		}
	}
	return false, nil
}

func PolicyExists(policyName string) (bool, error) {
	sys := auth.Client.Sys()
	policies, err := sys.ListPolicies()
	if err != nil {
		return false, err
	}

	for _, policy := range policies {
		if policy == policyName {
			return true, nil
		}
	}
	return false, nil
}
