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
	unMountPath string
	u4          string
	p4          string
)

// unmountCmd represents the unmount command
var unmountCmd = &cobra.Command{
	Use:   "unmount",
	Short: "Disable the kv engine at path",
	Long: `Disable a kv engine on path

	Example of the enable command(Keycloak Authentication):
		$ ./cliapp unmount --path=kv
	
	To use with Userpass Authentication:
		$ ./cliapp unmount --path=kv --user=username --pass=password

	To use a different instance:
		$ ./cliapp unmount --path=kv --user=username --pass=password --instance
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("user").Changed && cmd.Flag("pass").Changed {
			address := util.UpdateAddress(instance)
			auth.AuthenticateWithUserPass(u4, p4, address)
		} else {
			if cmd.Flag("instance").Changed {
				fmt.Println("Error: You must provide a username and password to use a different instance.")
				os.Exit(1)
			}
			address := util.UpdateAddress(false)
			auth.KeycloakAuth(address)
		}

		err := auth.Client.Sys().Unmount(unMountPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("KV secrets engine disabled at: " + unMountPath)
	},
}

func init() {
	rootCmd.AddCommand(unmountCmd)

	// path
	unmountCmd.Flags().StringVarP(&unMountPath, "path", "p", "", "path to disable engine on")

	if err := unmountCmd.MarkFlagRequired("path"); err != nil {
		fmt.Println(err)
	}

	// userpass
	unmountCmd.Flags().StringVarP(&u4, "user", "u", "", "Userpass username")
	unmountCmd.Flags().StringVarP(&p4, "pass", "a", "", "Userpass password")
	unmountCmd.MarkFlagsRequiredTogether("user", "pass")

	unmountCmd.Flags().BoolVarP(&instance, "instance", "i", false, "Use another Vault instance")
}
