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

	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

var (
	getMount   string
	getPath    string
	getVersion int
	u1         string
	p1         string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the value from vaults key-value store at a given key name",
	Long: ` 
	Retrieves the value from Vault's key-value store at the given key name. If no
	key exists with that name, an error is returned. If a key exists with that
	name but has no data, nothing is returned. It is also possible to view a value 
	with a different version than the current using the "--version" flag 
	
	Examples of the write command(Keycloak Authentication):
		$ ./cliapp get --mount=secret --path=secret/my-secret

		$ ./cliapp get --mount=secret --path=secret/my-secret --version=2
	
	To use with Userpass Authentication:
		$ ./cliapp get --mount=secret --path=secret/my-secret --user=username --pass=password

	To use a different instance:
		$ ./cliapp get --mount=secret --path=secret/my-secret --user=username --pass=password --instance
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if getVersion < 0 {
			fmt.Println("Error: Version of secret cannot be less than one")
			os.Exit(1)
		}

		if cmd.Flag("user").Changed && cmd.Flag("pass").Changed {
			address := util.UpdateAddress(instance)
			auth.AuthenticateWithUserPass(u1, p1, address)
		} else {
			if cmd.Flag("instance").Changed {
				fmt.Println("Error: You must provide a username and password to use a different instance.")
				os.Exit(1)
			}
			address := util.UpdateAddress(false)
			auth.KeycloakAuth(address)
		}

		util.ValidatePath(getPath)
		util.ValidatePath(getMount)
		if getVersion == 0 { // get current version
			secret, err := auth.Client.KVv2(getMount).Get(context.Background(), getPath)
			if err != nil {
				log.Fatalf("unable to read secret: %v", err)
			}
			key := secret.Data //deleted secret check
			versio := secret.VersionMetadata.Destroyed
			secretTimeDeletion := secret.VersionMetadata.DeletionTime.Format("2006-01-02 15:04:05")
			if key == nil && versio {
				fmt.Println("Secret Destroyed")
				fmt.Println("Time Deleted: " + secretTimeDeletion)
				os.Exit(0)
			} else if key == nil {
				fmt.Println("Secret Deleted, not Destroyed")
				fmt.Println("Time Deleted: " + secretTimeDeletion)
				os.Exit(0)
			}
			printSecret(secret)
		} else {
			secret, err := auth.Client.KVv2(getMount).GetVersion(context.Background(), getPath, getVersion)
			if err != nil {
				log.Fatalf("unable to read secret: %v", err)
			}
			key := secret.Data //deleted secret check
			versio := secret.VersionMetadata.Destroyed
			secretTimeDeletion := secret.VersionMetadata.DeletionTime.Format("2006-01-02 15:04:05")
			if key == nil && versio {
				fmt.Println("Secret Destroyed")
				fmt.Println("Time Deleted: " + secretTimeDeletion)
				os.Exit(0)
			} else if key == nil {
				fmt.Println("Secret Deleted, not Destroyed")
				fmt.Println("Time Deleted: " + secretTimeDeletion)
				os.Exit(0)
			}
			printSecret(secret)
		}

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	//mount
	getCmd.Flags().StringVarP(&getMount, "mount", "m", "", "The mount path to retrive secrets from")

	if err := getCmd.MarkFlagRequired("mount"); err != nil {
		fmt.Println(err)
	}

	// path
	getCmd.Flags().StringVarP(&getPath, "path", "p", "", "path to the secret")

	if err := getCmd.MarkFlagRequired("path"); err != nil {
		fmt.Println(err)
	}

	// version
	getCmd.Flags().IntVarP(&getVersion, "version", "v", 0, "version of secret")

	// userpass
	getCmd.Flags().StringVarP(&u1, "user", "u", "", "Userpass username")
	getCmd.Flags().StringVarP(&p1, "pass", "a", "", "Userpass password")
	getCmd.MarkFlagsRequiredTogether("user", "pass")

	getCmd.Flags().BoolVarP(&instance, "instance", "i", false, "Use another Vault instance")

}

func printSecret(secret *vault.KVSecret) {
	secretVersion := secret.VersionMetadata.Version
	secretTimeCreation := secret.VersionMetadata.CreatedTime.Format("2006-01-02 15:04:05")

	fmt.Printf("Version: %d \n", secretVersion)
	fmt.Println("Time Created: ", secretTimeCreation)
	for key, value := range secret.Data {
		fmt.Printf("Key: %s Value: %v\n", key, value)
	}
}
