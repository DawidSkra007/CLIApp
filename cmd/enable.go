/*
Copyright © 2023 Dawid Skraba <dawid.skraba@ucdconnect.ie>
*/
package cmd

import (
	"cliapp/auth"
	"cliapp/util"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	enablePath string
	u7         string
	p7         string
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a kv engine v2 on path",
	Long: `Enable a kv engine v2 on path

	Example of the enable command(Keycloak Authentication):
		$ ./cliapp enable --path=kv
	
	To use Userpass Authentication:
		$ ./cliapp enable --path=kv --user=username --pass=password

	To use a different instance:
		$ ./cliapp enable --path=kv --user=username --pass=password --instance
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("user").Changed && cmd.Flag("pass").Changed {
			address := util.UpdateAddress(instance)
			auth.AuthenticateWithUserPass(u7, p7, address)
		} else {
			if cmd.Flag("instance").Changed {
				fmt.Println("Error: You must provide a username and password to use a different instance.")
				os.Exit(1)
			}
			address := util.UpdateAddress(false)
			auth.KeycloakAuth(address)
		}

		_, err := auth.Client.Logical().Write("sys/mounts/"+enablePath, map[string]interface{}{
			"type": "kv",
			"options": map[string]interface{}{
				"version": "2",
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("KV secrets engine enabled at: " + enablePath)

		// update the policy files
		sys := auth.Client.Sys()
		policies, err := sys.ListPolicies()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, policy := range policies {
			existingPolicy, err := ReadPolicy(policy)
			if err != nil {
				fmt.Println(err)
				return
			}
			// update the policy
			// if policy contains 'admin' in the string then write the policy with the new rules
			if !strings.Contains(policy, "admin") { // if policy is not admin
				updatedPolicy := existingPolicy + `
				# New rules for path
				path "` + enablePath + `/data/*" {
					capabilities = ["create", "read", "update", "delete", "list", "patch"] 
				}

				path "` + enablePath + `/undelete/*" {
					capabilities = ["update"]
				}
				`
				UpdatePolicy(policy, updatedPolicy)
			} else { // if policy is admin
				updatedPolicy := existingPolicy +
					`
				# New rules for path
				path "` + enablePath + `/data/*" {
					capabilities = ["create", "read", "update", "delete", "list", "patch"] 
				}

				path "` + enablePath + `/metadata/*" {
				capabilities = ["list", "read", "delete"]
				}

				path "` + enablePath + `/undelete/*" {
					capabilities = ["update"]
				}

				path "` + enablePath + `/destroy/*" {
					capabilities = ["update"]
				}
				`
				UpdatePolicy(policy, updatedPolicy)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)

	// path
	enableCmd.Flags().StringVarP(&enablePath, "path", "p", "", "path to enable engine on")

	if err := enableCmd.MarkFlagRequired("path"); err != nil {
		fmt.Println(err)
	}

	// userpass
	enableCmd.Flags().StringVarP(&u7, "user", "u", "", "Userpass username")
	enableCmd.Flags().StringVarP(&p7, "pass", "a", "", "Userpass password")
	enableCmd.MarkFlagsRequiredTogether("user", "pass")

	enableCmd.Flags().BoolVarP(&instance, "instance", "i", false, "Use another Vault instance")
}

func ReadPolicy(policyName string) (string, error) {
	sys := auth.Client.Sys()
	policy, err := sys.GetPolicy(policyName)
	if err != nil {
		return "", fmt.Errorf("unable to read policy '%s': %v", policyName, err)
	}
	return policy, nil
}
