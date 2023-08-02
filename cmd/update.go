/*
Copyright © 2023 Dawid Skraba <dawid.skraba@ucdconnect.ie>
*/
package cmd

import (
	"cliapp/auth"
	"cliapp/util"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	upMountPath string
	upPath      string
	upKey       string
	upValue     string
	u3          string
	p3          string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the data to the corresponding path in the key-value store",
	Long: `
	Updates the data to the corresponding path in the key-value store.
	A mount, path and a key value pair is required to update the data.
	
	Examples of the update command(Keycloak Authentication):
		$ ./cliapp update --mount=secret --path=secret/my-secret --key=val --value=foo

	To use with Userpass Authentication:
		$ ./cliapp update --mount=secret --path=secret/my-secret --key=val --value=foo --user=user --pass=pass

	To use a different instance:
		$ ./cliapp update --mount=secret --path=secret/my-secret --key=val --value=foo --user=user --pass=pass --instance
		`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("user").Changed && cmd.Flag("pass").Changed {
			address := util.UpdateAddress(instance)
			auth.AuthenticateWithUserPass(u3, p3, address)
		} else {
			if cmd.Flag("instance").Changed {
				fmt.Println("Error: You must provide a username and password to use a different instance.")
				os.Exit(1)
			}
			address := util.UpdateAddress(false)
			auth.KeycloakAuth(address)
		}

		util.ValidatePath(upPath)
		util.ValidatePath(upMountPath)

		UpsecretData := make(map[string]interface{})
		UpsecretData[upKey] = upValue

		_, err := auth.Client.KVv2(upMountPath).Patch(context.Background(), upPath, UpsecretData)
		if err != nil {
			log.Fatalf("%v", err)
		}

		fmt.Println("Secret updated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	//mount
	updateCmd.Flags().StringVarP(&upMountPath, "mount", "m", "", "The mount path to update secret")

	if err := updateCmd.MarkFlagRequired("mount"); err != nil {
		fmt.Println(err)
	}

	//path
	updateCmd.Flags().StringVarP(&upPath, "path", "p", "", "path to the secret")

	if err := updateCmd.MarkFlagRequired("path"); err != nil {
		fmt.Println(err)
	}

	//key
	updateCmd.Flags().StringVarP(&upKey, "key", "k", "", "key of the secret")

	if err := updateCmd.MarkFlagRequired("key"); err != nil {
		fmt.Println(err)
	}

	//value
	updateCmd.Flags().StringVarP(&upValue, "value", "v", "", "value of the secret")

	if err := updateCmd.MarkFlagRequired("value"); err != nil {
		fmt.Println(err)
	}

	// userpass
	updateCmd.Flags().StringVarP(&u3, "user", "u", "", "Userpass username")
	updateCmd.Flags().StringVarP(&p3, "pass", "a", "", "Userpass password")
	updateCmd.MarkFlagsRequiredTogether("user", "pass")

	updateCmd.Flags().BoolVarP(&instance, "instance", "i", false, "Use another Vault instance")

}
