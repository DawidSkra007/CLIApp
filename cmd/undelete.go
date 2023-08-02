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
	undelMount   string
	undelPath    string
	undelVersion string
	u5           string
	p5           string
)

var delVersions []int

// undeleteCmd represents the undelete command
var undeleteCmd = &cobra.Command{
	Use:   "undelete",
	Short: "Undelets the data for the provided version",
	Long: `
	Undeletes the data for the provided version and path in the key-value store.
  	This restores the data, allowing it to be returned on get requests. Version must
	be provided or multiple versions.

	Examples of the undelete command(Keycloak Authentication):
		$ ./cliapp undelete --mount=secret --path=secret/my-secret --version=2

		$ ./cliapp undelete --mount=secret --path=secret/my-secret --version=2,3,4

	To use Userpass Authentication:
		$ ./cliapp undelete --mount=secret --path=secret/my-secret --version=2 --user=username --pass=password

	To use a different instance:
		$ ./cliapp undelete --mount=secret --path=secret/my-secret --version=2 --user=username --pass=password --instance
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("user").Changed && cmd.Flag("pass").Changed {
			address := util.UpdateAddress(instance)
			auth.AuthenticateWithUserPass(u5, p5, address)
		} else {
			if cmd.Flag("instance").Changed {
				fmt.Println("Error: You must provide a username and password to use a different instance.")
				os.Exit(1)
			}
			address := util.UpdateAddress(false)
			auth.KeycloakAuth(address)
		}

		util.ValidatePath(undelMount)
		util.ValidatePath(undelPath)

		if strings.Contains(undelVersion, ",") { // multiple versions
			vers := strings.Split(undelVersion, ",")
			for _, kesy := range vers {
				i, err := strconv.Atoi(kesy)
				if err != nil || i < 1 {
					fmt.Println("Error: versions provided must be integers")
					os.Exit(1)
				}
				delVersions = append(delVersions, i)
			}

			err := auth.Client.KVv2(undelMount).Undelete(context.Background(), undelPath, delVersions)
			if err != nil {
				log.Fatalf("Secret not unDeleted error: %d ", err)
			}
		} else if undelVersion != "" { // one version
			i, err := strconv.Atoi(undelVersion)
			if err != nil || i < 1 {
				fmt.Println("Error: version provided must be integer")
				os.Exit(1)
			}
			delVersions = append(delVersions, i)
			err = auth.Client.KVv2(undelMount).Undelete(context.Background(), undelPath, delVersions)
			if err != nil {
				log.Fatalf("Secret not unDeleted error: %d ", err)
			}
		}

		fmt.Printf("Secret undeleted.\n")

	},
}

func init() {
	rootCmd.AddCommand(undeleteCmd)

	//mount
	undeleteCmd.Flags().StringVarP(&undelMount, "mount", "m", "", "The mount path to undelete secrets")

	if err := undeleteCmd.MarkFlagRequired("mount"); err != nil {
		fmt.Println(err)
	}

	// path
	undeleteCmd.Flags().StringVarP(&undelPath, "path", "p", "", "path to the secret")

	if err := undeleteCmd.MarkFlagRequired("path"); err != nil {
		fmt.Println(err)
	}

	// version
	undeleteCmd.Flags().StringVarP(&undelVersion, "version", "v", "", "version of secret")

	if err := undeleteCmd.MarkFlagRequired("version"); err != nil {
		fmt.Println(err)
	}

	// userpass
	undeleteCmd.Flags().StringVarP(&u5, "user", "u", "", "Userpass username")
	undeleteCmd.Flags().StringVarP(&p5, "pass", "a", "", "Userpass password")
	undeleteCmd.MarkFlagsRequiredTogether("user", "pass")

	undeleteCmd.Flags().BoolVarP(&instance, "instance", "i", false, "Use another Vault instance")

}
