package auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/userpass"
	"golang.org/x/oauth2"
)

const (
	clientID    = "vault-client"
	keycloakURL = "http://127.0.0.1:8080/auth/realms/my_realm"
	redirectURL = "http://127.0.0.1:3000/callback"
)

var clientSecret string = ""

func KeycloakAuth(address string) {
	config := vault.DefaultConfig()
	config.Address = address
	localClient, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	secretBytes, err := ioutil.ReadFile("client_secret.txt")
	if err != nil {
		panic(err)
	}
	clientSecret = strings.TrimSpace(string(secretBytes))

	provider, err := oidc.NewProvider(context.Background(), keycloakURL)
	if err != nil {
		log.Fatalf("unable to initialize OIDC provider: %v", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	done := make(chan struct{})
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		token, err := oauth2Config.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}

		idTokenString := token.Extra("id_token").(string)

		vaultToken, err := authWithToken(idTokenString, localClient)
		if err != nil {
			log.Fatalf("unable to authenticate to Vault: %v", err)
		}

		Connect(vaultToken, address, localClient)

		close(done)
	})

	go func() {
		log.Fatal(http.ListenAndServe(":3000", nil))
	}()

	authURL := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Open the following URL in your browser to authenticate with Keycloak:\n\n%s\n\n", authURL)

	<-done
}

func authWithToken(idToken string, localClient *vault.Client) (string, error) {
	params := map[string]interface{}{
		"role": "user-policy",
		"jwt":  idToken,
	}

	secret, err := localClient.Logical().Write("auth/jwt/login", params) // JWT auth method (oidc -> ui)
	if err != nil {
		return "", fmt.Errorf("unable to authenticate with Vault using JWT: %w", err)
	}

	return secret.Auth.ClientToken, nil
}

func AuthenticateWithUserPass(username, password, address string) error {
	config := vault.DefaultConfig()
	config.Address = address
	localClient, err := vault.NewClient(config)
	if err != nil {
		return fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	userpassAuth, err := auth.NewUserpassAuth(username, &auth.Password{FromString: password})
	if err != nil {
		return fmt.Errorf("unable to initialize userpass auth method: %w", err)
	}
	authInfo, err := localClient.Auth().Login(context.TODO(), userpassAuth)
	if err != nil {
		fmt.Println("Error during login:", err)
		return fmt.Errorf("unable to login to userpass auth method: %w", err)
	}
	if authInfo == nil {
		return fmt.Errorf("no auth info was returned after login")
	}

	vaultToken := authInfo.Auth.ClientToken
	Connect(vaultToken, address, localClient)
	return nil
}

func GetAdminToken(keycloakUsername, keycloakPassword, keycloakClientID, keycloakSecret, keycloakAuthURL string) (*oauth2.Token, error) {
	cfg := &oauth2.Config{
		ClientID:     keycloakClientID,
		ClientSecret: keycloakSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: keycloakAuthURL,
		},
	}

	ctx := context.Background()
	token, err := cfg.PasswordCredentialsToken(ctx, keycloakUsername, keycloakPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to get Keycloak token: %w", err)
	}

	return token, nil
}
