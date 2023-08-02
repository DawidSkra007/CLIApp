/*
Copyright © 2023 Dawid Skraba <dawid.skraba@ucdconnect.ie>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cliapp",
	Short: `A CLI app for secrets managment`,
	Long: `
░█████╗░██╗░░░░░██╗░█████╗░██████╗░██████╗░
██╔══██╗██║░░░░░██║██╔══██╗██╔══██╗██╔══██╗
██║░░╚═╝██║░░░░░██║███████║██████╔╝██████╔╝
██║░░██╗██║░░░░░██║██╔══██║██╔═══╝░██╔═══╝░
╚█████╔╝███████╗██║██║░░██║██║░░░░░██║░░░░░
░╚════╝░╚══════╝╚═╝╚═╝░░╚═╝╚═╝░░░░░╚═╝░░░░░

Welcome to the CLI app for secrets managment in Vault. 
You can use this app to create and manage policies, users, secrets and more.
There are two ways to authenticate to Vault. You can use the keycloak server or userpass.
If using userpass, use the commands with the -u and -a flags(see help for more info).
The default is keycloak, so just enter your username and password when prompted in browser 
This app is a part of my final year project in University College Dublin.
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cliapp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
