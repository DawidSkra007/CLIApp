/*
Copyright © 2023 Dawid Skraba <dawid.skraba@ucdconnect.ie>
*/
package auth

import (
	vault "github.com/hashicorp/vault/api"
)

var Client *vault.Client

func Connect(token, address string, client *vault.Client) {
	client.SetToken(token)

	Client = client
}
