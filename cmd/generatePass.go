/*
Copyright © 2023 Dawid Skraba <dawid.skraba@ucdconnect.ie>
*/
package cmd

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	length           int
	includeNumbers   bool
	includeSymbols   bool
	includeLowercase bool
	includeUppercase bool
	outputFile       string = "password.txt"
)

// generatePassCmd represents the generatePass command
var generatePassCmd = &cobra.Command{
	Use:   "generatePass",
	Short: "Generate a secure password",
	Long: `Generate a secure password with multiple given inputs by the user.

	Examples of the generatePass command:
		$ ./cliapp generatePass --length 16 --numbers --symbols --lowercase
	`,
	Run: func(cmd *cobra.Command, args []string) {
		charsets := []string{
			"0123456789",
			"!@#$%^&*()-_=+<>,.?/:;{}[]|",
			"abcdefghijklmnopqrstuvwxyz",
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		}

		var charset string = ""
		if includeNumbers {
			charset += charsets[0]
		}
		if includeSymbols {
			charset += charsets[1]
		}
		if includeLowercase {
			charset += charsets[2]
		}
		if includeUppercase {
			charset += charsets[3]
		}

		if charset == "" {
			fmt.Println("Error: At least one character set must be included.")
			os.Exit(1)
		}

		password, err := generateRandomPassword(length, charset)
		if err != nil {
			log.Fatalf("Error generating password: %v", err)
		}

		err = ioutil.WriteFile(outputFile, []byte(password), 0644)
		if err != nil {
			log.Fatalf("Error writing password to file: %v", err)
		}

		fmt.Printf("Generated password in file: %s\n", outputFile)

	},
}

func init() {
	rootCmd.AddCommand(generatePassCmd)

	generatePassCmd.Flags().IntVarP(&length, "length", "l", 15, "Length of the generated password")

	generatePassCmd.Flags().BoolVarP(&includeNumbers, "numbers", "n", false, "Include numbers")

	generatePassCmd.Flags().BoolVarP(&includeSymbols, "symbols", "s", false, "Include symbols")

	generatePassCmd.Flags().BoolVarP(&includeLowercase, "lowercase", "w", false, "Include lowercase letters")

	generatePassCmd.Flags().BoolVarP(&includeUppercase, "uppercase", "u", false, "Include uppercase letters")

}

func generateRandomPassword(length int, charset string) (string, error) {
	var password strings.Builder
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		password.WriteByte(charset[n.Int64()])
	}

	return password.String(), nil
}
