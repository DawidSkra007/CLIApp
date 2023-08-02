/*
Copyright © 2023 Dawid Skraba <dawid.skraba@ucdconnect.ie>
*/
package cmd

import (
	"bufio"
	"cliapp/auth"
	"cliapp/util"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	mountPath string
	path      string
	key       string
	value     string
	u2        string
	p2        string
)

var keys []string
var values []string
var Stdininput string

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Writes data to the vault at the path",
	Long: `	
	This data should be specified as a key=value pair. The input could also be specified
	from a file, when using the @ symbol at the start of the value variable, this file being made up of multiple key value pairs. 
	Multiple key value pairs are also permitted. You can write from stdin also, when specified the value with the '-' sign. The mount 
	should be specified with a already running engine at that path(check documentation
	for full details) and this is where the secrets will be mounted.

	Examples of the write command(Keycloak Authentication):
		$ ./cliapp write --mount=secret --path=secret/my-secret --key=customer_name --value=Apple_Inc.
		
		$ ./cliapp write --mount=secret --path=secret/my-secret --key=key1,key2 --value=val1,val2

		$ ./cliapp write --mount=secret --path=secret/my-secret --key=key --value=-

		$ ./cliapp write --mount=secret --path=secret/my-secret --key=@file.json --value=file
	
	To use with Userpass Authentication:
		$ ./cliapp write --mount=secret --path=secret/my-secret --key=customer_name --value=Apple_Inc. --user=username --pass=password

	To use a different instance:
		$ ./cliapp write --mount=secret --path=secret/my-secret --key=customer_name --value=Apple_Inc. --user=username --pass=password --instance
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("user").Changed && cmd.Flag("pass").Changed {
			address := util.UpdateAddress(instance)
			auth.AuthenticateWithUserPass(u2, p2, address)
		} else {
			if cmd.Flag("instance").Changed {
				fmt.Println("Error: You must provide a username and password to use a different instance.")
				os.Exit(1)
			}
			address := util.UpdateAddress(false)
			auth.KeycloakAuth(address)
		}

		if strings.Contains(key, ",") || strings.Contains(value, ",") { // multiple key value pairs
			key := strings.Split(key, ",")
			value := strings.Split(value, ",")
			for _, kesy := range key {
				keys = append(keys, kesy)
			}
			for _, vasl := range value {
				values = append(values, vasl)
			}
		} else {
			keys = append(keys, key)
			values = append(values, value)
		}

		if len(keys) != len(values) && keys[0] != "@" { // check if all keys have corresponding values
			fmt.Println("Error: All keys must have corresponding values")
			os.Exit(1)
		}

		if strings.Contains(key, "@") { // JSON file as input
			file := strings.Split(key, "@")
			var nfile = file[1]
			jsonBytes, err := ioutil.ReadFile(nfile)
			if err != nil {
				fmt.Println("Error reading JSON file:", err)
				return
			}

			var data map[string]string
			err = json.Unmarshal(jsonBytes, &data)
			if err != nil {
				fmt.Println("Error unmarshalling JSON data:", err)
				return
			}
			keys = []string{}
			values = []string{}
			for k, v := range data {
				keys = append(keys, k)
				values = append(values, v)
			}
		}

		for i, s := range values { // STDIN
			if s == "-" {
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					Stdininput = scanner.Text()
				}
				values[i] = Stdininput
			}
		}

		for i := 0; i < len(keys) && i < len(values); i++ {
			util.ValidateKVsecret(keys[i], values[i])
		}
		util.ValidatePath(path)
		util.ValidatePath(mountPath)
		secretData := make(map[string]interface{})
		for i := 0; i < len(keys) && i < len(values); i++ {
			key := keys[i]
			value := values[i]
			secretData[key] = value
		}

		_, err := auth.Client.KVv2(mountPath).Put(context.Background(), path, secretData)
		if err != nil {
			log.Fatalf("%v", err)
		}

		fmt.Println("Secret written successfully.")
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)

	//mount
	writeCmd.Flags().StringVarP(&mountPath, "mount", "m", "", "The mount path to write secrets to")

	if err := writeCmd.MarkFlagRequired("mount"); err != nil {
		fmt.Println(err)
	}

	//path
	writeCmd.Flags().StringVarP(&path, "path", "p", "", "path to the secret")

	if err := writeCmd.MarkFlagRequired("path"); err != nil {
		fmt.Println(err)
	}

	//key
	writeCmd.Flags().StringVarP(&key, "key", "k", "", "key of the secret")

	if err := writeCmd.MarkFlagRequired("key"); err != nil {
		fmt.Println(err)
	}

	//value
	writeCmd.Flags().StringVarP(&value, "value", "v", "", "value of the secret")

	if err := writeCmd.MarkFlagRequired("value"); err != nil {
		fmt.Println(err)
	}

	// userpass
	writeCmd.Flags().StringVarP(&u2, "user", "u", "", "Userpass username")
	writeCmd.Flags().StringVarP(&p2, "pass", "a", "", "Userpass password")
	writeCmd.MarkFlagsRequiredTogether("user", "pass")

	writeCmd.Flags().BoolVarP(&instance, "instance", "i", false, "Use another Vault instance")

}
