#!/bin/bash

VAULT_PID=$(pgrep -f "vault server")
if [ ! -z "$VAULT_PID" ]; then
  echo "Shutting down existing Vault server (PID: $VAULT_PID)..."
  kill "$VAULT_PID"
  sleep 2
fi

echo "Starting Vault server..."
vault server -config=config.hcl > vault.log 2>&1 &

echo "Waiting for Vault server to initialize..."
sleep 10

export VAULT_ADDR="http://127.0.0.1:8200"

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
echo "Waiting for raft storage to finish initialization..."
sleep 10

# write the admin and user policies

vault policy write admin-policy policies/admin-policy.hcl
vault policy write user-policy policies/user-policy.hcl

# enable the userpass auth method

vault auth enable userpass
vault write auth/userpass/users/admin password=admin policies=admin-policy
vault write auth/userpass/users/user password=pass policies=user-policy

# enable the kv secrets engine at the 'kv' path

vault secrets enable -version=2 kv

# setup oidc method
CLIENT_SECRET=$(sed 's/^[[:space:]]*//;s/[[:space:]]*$//' client_secret.txt)

vault auth enable oidc # in ui, token type is not set to service and no TTLs

vault write auth/oidc/role/default \
    allowed_redirect_uris="http://localhost:8250/oidc/callback,http://localhost:8200/ui/vault/auth/oidc/oidc/callback" \
    user_claim="sub" \
    groups_claim="groups" \
    policies="default"

vault write auth/oidc/config oidc_discovery_url=http://localhost:8080/auth/realms/my_realm oidc_client_id=vault-client oidc_client_secret="$CLIENT_SECRET" default_role=default 

GROUP_ID=$(vault write -format=json identity/group @group_payload.json | jq -r '.data.id')

MOUNT_ACCESSOR=$(vault auth list -format=json | jq -r '.["oidc/"].accessor')

cat << EOF > group_alias_payload.json
{
    "name": "vault-client",
    "mount_accessor": "$MOUNT_ACCESSOR",
    "canonical_id": "$GROUP_ID"
}
EOF

vault write -format=json identity/group-alias @group_alias_payload.json

# jwt auth method

vault auth enable jwt

vault write auth/jwt/role/user-policy role_type=jwt \
  bound_audiences="vault-client" \
  user_claim="sub" \
  policies="user-policy" \
  ttl=1h \
  max_ttl=24h

vault write auth/jwt/config \
    oidc_discovery_url="http://localhost:8080/auth/realms/my_realm" \
    default_role="user-policy"

