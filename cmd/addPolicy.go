/*
Copyright © 2023 Dawid Skraba <dawid.skraba@ucdconnect.ie>
*/
package cmd

import (
	"cliapp/auth"
	"cliapp/util"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	policy string
	u11    string
	p11    string
)

// addPolicyCmd represents the addPolicy command
var addPolicyCmd = &cobra.Command{
	Use:   "addPolicy",
	Short: "Add a policy to the Vault server",
	Long: `
	Add a policy to the Vault server. The policy should be specified as a file.

	Examples of the addPolicy command(Keycloak Authentication):
		$ ./cliapp addPolicy --policy=@user-policy.hcl

	To use Userpass Authentication:
		$ ./cliapp addPolicy --policy=@user-policy.hcl --user=username --pass=password
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("user").Changed && cmd.Flag("pass").Changed {
			address := util.UpdateAddress(instance)
			auth.AuthenticateWithUserPass(u11, p11, address)
		} else {
			if cmd.Flag("instance").Changed {
				fmt.Println("Error: You must provide a username and password to use a different instance.")
				os.Exit(1)
			}
			address := util.UpdateAddress(false)
			auth.KeycloakAuth(address)
		}

		if strings.Contains(policy, "@") {
			policy = strings.TrimPrefix(policy, "@")
			// delete last 4 characters from policy name
			policyCHeck := policy[:len(policy)-4]
			check, err := PolicyExists(policyCHeck) // check if policy already exists
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if check == true { //exit if policy already exists
				fmt.Println("Policy already exists")
				os.Exit(1)
			} else { //add policy if it does not exist
				if err := WritePolicy(policy, "./policies/"+policy); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		} else {
			fmt.Println("Policy file not specified Correctly")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addPolicyCmd)

	//policy
	addPolicyCmd.Flags().StringVarP(&policy, "policy", "p", "", "policy file to be added")

	if err := addPolicyCmd.MarkFlagRequired("policy"); err != nil {
		fmt.Println(err)
	}

	// userpass
	addPolicyCmd.Flags().StringVarP(&u11, "user", "u", "", "Userpass username")
	addPolicyCmd.Flags().StringVarP(&p11, "pass", "a", "", "Userpass password")
	addPolicyCmd.MarkFlagsRequiredTogether("user", "pass")

	addPolicyCmd.Flags().BoolVarP(&instance, "instance", "i", false, "Use another Vault instance")

}

func WritePolicy(policyName, policyFile string) error {
	hclContent, err := ioutil.ReadFile(policyFile)
	if err != nil {
		return fmt.Errorf("unable to read policy file: %v", err)
	}

	policyName = policyName[:len(policyName)-4]

	sys := auth.Client.Sys()
	if err := sys.PutPolicy(policyName, string(hclContent)); err != nil {
		return fmt.Errorf("unable to write policy: %v", err)
	}

	fmt.Printf("Policy '%s' written successfully.\n", policyName)
	return nil
}

func UpdatePolicy(policyName, content string) error {
	sys := auth.Client.Sys()
	if err := sys.PutPolicy(policyName, content); err != nil {
		return fmt.Errorf("unable to write policy: %v", err)
	}

	fmt.Printf("Policy '%s' written successfully.\n", policyName)
	return nil
}
