/*
Copyright © 2023 Dawid Skraba <dawid.skraba@ucdconnect.ie>
*/
package cmd

import (
	"cliapp/auth"
	"cliapp/util"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	u6       string
	p6       string
	instance bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available mount paths",
	Long: `
	List available mount paths that are available in the Vault server.
	These are the paths that you can use to write secrets to and perform other operations on.

	Example of the list command(Keycloak Authentication):
		$ ./cliapp list

	To use Userpass Authentication:
		$ ./cliapp list --user=username --pass=password

	To use a different instance:
		$ ./cliapp list --user=username --pass=password --instance
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("user").Changed && cmd.Flag("pass").Changed {
			address := util.UpdateAddress(instance)
			auth.AuthenticateWithUserPass(u6, p6, address)
		} else {
			if cmd.Flag("instance").Changed {
				fmt.Println("Error: You must provide a username and password to use a different instance.")
				os.Exit(1)
			}
			address := util.UpdateAddress(false)
			auth.KeycloakAuth(address)
		}

		mounts, err := auth.Client.Sys().ListMounts()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Available mount paths:")
		for path := range mounts {
			fmt.Println(path)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// userpass
	listCmd.Flags().StringVarP(&u6, "user", "u", "", "Userpass username")
	listCmd.Flags().StringVarP(&p6, "pass", "a", "", "Userpass password")
	listCmd.MarkFlagsRequiredTogether("user", "pass")

	listCmd.Flags().BoolVarP(&instance, "instance", "i", false, "Use another Vault instance")

}
