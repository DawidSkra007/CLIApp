#!/bin/bash

export VAULT_ADDR="http://localhost:8400"

INIT_OUTPUT=$(curl --request PUT --data '{"secret_shares": 5, "secret_threshold": 3}' ${VAULT_ADDR}/v1/sys/init)
INIT_EXIT_CODE=$?

if [ $INIT_EXIT_CODE -ne 0 ]; then
  echo "Error during Vault initialization. Please check your Vault server configuration or connection."
  exit 1
fi

UNSEAL_KEYS=$(echo "$INIT_OUTPUT" | jq -r '.keys[]')
ROOT_TOKEN=$(echo "$INIT_OUTPUT" | jq -r '.root_token')

echo "Initial Root Token: $ROOT_TOKEN"

KEY_COUNT=0
for KEY in $UNSEAL_KEYS; do
  vault operator unseal "$KEY"
  KEY_COUNT=$((KEY_COUNT + 1))
  if [ $KEY_COUNT -ge 3 ]; then
    break
  fi
done

export VAULT_TOKEN="$ROOT_TOKEN"

echo "Vault server is now unsealed and ready to use."

# write the admin and user policies

vault policy write admin-policy policies/admin-policy.hcl
vault policy write user-policy policies/user-policy.hcl

# enable the userpass auth method

vault auth enable userpass
vault write auth/userpass/users/admin password=admin policies=admin-policy
vault write auth/userpass/users/user password=pass policies=user-policy

# enable the kv secrets engine at the 'kv' path

vault secrets enable -version=2 kv



