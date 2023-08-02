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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	desMount   string
	desPath    string
	desVersion string
	u8         string
	p8         string
)

var desVersions []int

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Permanently removes the specified versions' data from the key-value store",
	Long: `
	Permanently removes the specified versions' data from the key-value store. This
	also destroys all the metadata.If no key exists at the path, no action is taken. 
	A Version must be provided or multiple versions to destroy a secret at the path.

	Example of the destroy command(Keycloak Authentication):
		$ ./cliapp destroy --mount=test --path=path --version=1

		$ ./cliapp destroy --mount=test --path=path --version=1,2,3

	To use Userpass Authentication:
		$ ./cliapp destroy --mount=test --path=path --version=1 --user=username --pass=password

	To use a different instance:
		$ ./cliapp destroy --mount=test --path=path --version=1 --user=username --pass=password --instance
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("user").Changed && cmd.Flag("pass").Changed {
			address := util.UpdateAddress(instance)
			auth.AuthenticateWithUserPass(u8, p8, address)
		} else {
			if cmd.Flag("instance").Changed {
				fmt.Println("Error: You must provide a username and password to use a different instance.")
				os.Exit(1)
			}
			address := util.UpdateAddress(false)
			auth.KeycloakAuth(address)
		}

		util.ValidatePath(desMount)
		util.ValidatePath(desPath)

		if strings.Contains(desVersion, ",") { // multiple versions
			vers := strings.Split(desVersion, ",")
			for _, kesy := range vers {
				i, err := strconv.Atoi(kesy)
				if err != nil || i < 1 {
					fmt.Println("Error: versions provided must be integers")
					os.Exit(1)
				}
				desVersions = append(desVersions, i)
			}

			err := auth.Client.KVv2(desMount).Destroy(context.Background(), desPath, desVersions)
			if err != nil {
				log.Fatalf("Secret not Destroyed error: %d ", err)
			}
		} else if desVersion != "" { // one version
			i, err := strconv.Atoi(desVersion)
			if err != nil || i < 1 {
				fmt.Println("Error: version provided must be integer")
				os.Exit(1)
			}
			desVersions = append(delVersions, i)
			err = auth.Client.KVv2(desMount).Destroy(context.Background(), desPath, desVersions)
			if err != nil {
				log.Fatalf("Secret not Destroyed error: %d ", err)
			}
		}

		fmt.Printf("Secret Destroyed.\n")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	//mount
	destroyCmd.Flags().StringVarP(&desMount, "mount", "m", "", "The mount path to destroy secrets")

	if err := destroyCmd.MarkFlagRequired("mount"); err != nil {
		fmt.Println(err)
	}

	// path
	destroyCmd.Flags().StringVarP(&desPath, "path", "p", "", "path to the secret")

	if err := destroyCmd.MarkFlagRequired("path"); err != nil {
		fmt.Println(err)
	}

	// version
	destroyCmd.Flags().StringVarP(&desVersion, "version", "v", "", "version of secret")

	if err := destroyCmd.MarkFlagRequired("version"); err != nil {
		fmt.Println(err)
	}

	// userpass
	destroyCmd.Flags().StringVarP(&u8, "user", "u", "", "Userpass username")
	destroyCmd.Flags().StringVarP(&p8, "pass", "a", "", "Userpass password")
	destroyCmd.MarkFlagsRequiredTogether("user", "pass")

	destroyCmd.Flags().BoolVarP(&instance, "instance", "i", false, "Use another Vault instance")
}
