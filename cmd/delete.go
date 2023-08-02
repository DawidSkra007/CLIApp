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
	delMount   string
	delPath    string
	delVersion string
	u9         string
	p9         string
)

var versions []int

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes the data for the provided version",
	Long: `
	Deletes the data for the provided version and path in the key-value store. The
	versioned data will not be fully removed, but marked as deleted and will no
	longer be returned in normal get requests. If no version specified the latest 
	version of the secret will be deleted. If the version tag is
	provided then that secret version will be deleted. A mounting path and the path to
	the secret are needed.

	Examples of the delete command(Keycloak Authentication):
		$ ./cliapp delete --mount=secret --path=secret/my-secret

		$ ./cliapp delete --mount=secret --path=secret/my-secret --version=2

		$ ./cliapp delete --mount=secret --path=secret/my-secret --version=2,3,4

	To use Userpass Authentication:
		$ ./cliapp delete --mount=secret --path=secret/my-secret --user=username --pass=password

	To use a different instance:
		$ ./cliapp delete --mount=secret --path=secret/my-secret --user=username --pass=password --instance
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("user").Changed && cmd.Flag("pass").Changed {
			address := util.UpdateAddress(instance)
			auth.AuthenticateWithUserPass(u9, p9, address)
		} else {
			if cmd.Flag("instance").Changed {
				fmt.Println("Error: You must provide a username and password to use a different instance.")
				os.Exit(1)
			}
			address := util.UpdateAddress(false)
			auth.KeycloakAuth(address)
		}

		check()
		util.ValidatePath(delMount)
		util.ValidatePath(delPath)

		if strings.Contains(delVersion, ",") { // multiple versions
			vers := strings.Split(delVersion, ",")
			for _, kesy := range vers {
				i, err := strconv.Atoi(kesy)
				if err != nil || i < 1 {
					fmt.Println("Error: versions provided must be integers")
					os.Exit(1)
				}
				versions = append(versions, i)
			}

			err := auth.Client.KVv2(delMount).DeleteVersions(context.Background(), delPath, versions)
			if err != nil {
				log.Fatalf("Secret not Deleted error: %d ", err)
			}
		} else if delVersion != "" { // one version
			i, err := strconv.Atoi(delVersion)
			if err != nil || i < 1 {
				fmt.Println("Error: version provided must be integer")
				os.Exit(1)
			}
			versions = append(versions, i)

			err = auth.Client.KVv2(delMount).DeleteVersions(context.Background(), delPath, versions)
			if err != nil {
				log.Fatalf("Secret not Deleted error: %d ", err)
			}
		} else { //current version
			err := auth.Client.KVv2(delMount).Delete(context.Background(), delPath)
			if err != nil {
				log.Fatalf("Secret not Deleted error: %d ", err)
			}
		}

		fmt.Printf("Secret deleted.\n")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	//mount
	deleteCmd.Flags().StringVarP(&delMount, "mount", "m", "", "The mount path to retrive secrets from")

	if err := deleteCmd.MarkFlagRequired("mount"); err != nil {
		fmt.Println(err)
	}

	// path
	deleteCmd.Flags().StringVarP(&delPath, "path", "p", "", "path to the secret")

	if err := deleteCmd.MarkFlagRequired("path"); err != nil {
		fmt.Println(err)
	}

	// version
	deleteCmd.Flags().StringVarP(&delVersion, "version", "v", "", "version of secret")

	// userpass
	deleteCmd.Flags().StringVarP(&u9, "user", "u", "", "Userpass username")
	deleteCmd.Flags().StringVarP(&p9, "pass", "a", "", "Userpass password")
	deleteCmd.MarkFlagsRequiredTogether("user", "pass")

	deleteCmd.Flags().BoolVarP(&instance, "instance", "i", false, "Use another Vault instance")

}

func check() {
	_, error := auth.Client.KVv2(delMount).Get(context.Background(), delPath)
	if error != nil {
		log.Fatalf("Unable to read path: %s", path)
	}
}
