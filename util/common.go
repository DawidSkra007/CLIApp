/*
Copyright © 2023 Dawid Skraba <dawid.skraba@ucdconnect.ie>
*/
package util

import (
	"fmt"
	"os"
)

func ValidateKVsecret(key string, value string) {
	if key == "" || value == "" {
		fmt.Println("Error: A key and value are required to write a KV secret")
		os.Exit(1)
	} else {
		return
	}
}

func ValidatePath(path string) {
	if path == "" {
		fmt.Println("Error: A path is required")
		os.Exit(1)
	} else {
		return
	}
}

func GetAuth(user, pass string) bool {
	if user != "" && pass != "" {
		return true
	} else {
		return false
	}
}

func UpdateAddress(instance bool) string {
	address := ""
	if instance {
		address = "http://127.0.0.1:8400"
	} else {
		address = "http://127.0.0.1:8200"
	}
	return address
}
