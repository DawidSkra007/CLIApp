#!/bin/bash

KEYCLOAK_BASE_URL="http://localhost:8080/auth"  
KEYCLOAK_USER="admin"
KEYCLOAK_PASSWORD="password"
REALM_NAME="my_realm"
CLIENT_NAME="vault-client"
# sample user
KC_USERNAME="user"
FIRST_NAME="user"
LAST_NAME="user"
EMAIL="user@gmail.com"
PASSWORD="foo"
MAPPER_NAME="groups"


echo "Fetching access token..."
ACCESS_TOKEN=$(curl -s -X POST "${KEYCLOAK_BASE_URL}/realms/master/protocol/openid-connect/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=${KEYCLOAK_USER}" \
    -d "password=${KEYCLOAK_PASSWORD}" \
    -d "grant_type=password" \
    -d "client_id=admin-cli" | jq -r '.access_token')

# Create the realm
echo "Creating realm..."
curl -X POST "${KEYCLOAK_BASE_URL}/admin/realms" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{\"id\": \"${REALM_NAME}\", \"realm\": \"${REALM_NAME}\", \"enabled\": true, \"registrationAllowed\": true}"

# Create the client
echo "Creating client..."
curl -X POST "${KEYCLOAK_BASE_URL}/admin/realms/${REALM_NAME}/clients" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
          \"clientId\": \"${CLIENT_NAME}\",
          \"name\": \"${CLIENT_NAME}\",
          \"enabled\": true,
          \"standardFlowEnabled\": true,
          \"directAccessGrantsEnabled\": true,
          \"rootUrl\": \"http://localhost:8200\",
          \"redirectUris\": [
            \"http://127.0.0.1:3000/callback\",
            \"http://localhost:8200/ui/vault/auth/oidc/oidc/callback\",
            \"http://localhost:8250/oidc/callback\"
          ],
          \"webOrigins\": [
            \"http://127.0.0.1:3000\",
            \"http://localhost:8200\"
          ],
          \"attributes\": { 
          \"backchannel_logout_session_required\": true
           }
        }"

echo "Client '${CLIENT_NAME}' created in the '${REALM_NAME}' realm with the desired settings."


CLIENT_ID=$(curl -X GET "${KEYCLOAK_BASE_URL}/admin/realms/${REALM_NAME}/clients" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" | jq -r --arg CLIENT_NAME "$CLIENT_NAME" '.[] | select(.clientId==$CLIENT_NAME) | .id')

#client secret
CLIENT_SECRET=$(curl -X POST "${KEYCLOAK_BASE_URL}/admin/realms/${REALM_NAME}/clients/${CLIENT_ID}/client-secret" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" | jq -r '.value')

echo "Client Secret: ${CLIENT_SECRET}"
echo "${CLIENT_SECRET}" > client_secret.txt


# Create the user
echo "Creating user..."
curl -X POST "${KEYCLOAK_BASE_URL}/admin/realms/${REALM_NAME}/users" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
          \"username\": \"${KC_USERNAME}\",
          \"firstName\": \"${FIRST_NAME}\",
          \"lastName\": \"${LAST_NAME}\",
          \"email\": \"${EMAIL}\",
          \"enabled\": true,
          \"credentials\": [
            {
              \"type\": \"password\",
              \"value\": \"${PASSWORD}\",
              \"temporary\": false
            }
          ]
        }"

# create group and add user to group
echo "Creating group and adding user..."
curl -X POST "${KEYCLOAK_BASE_URL}/admin/realms/${REALM_NAME}/groups" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
          \"name\": \"vault-client\",
          \"path\": \"vault-client\",
          \"attributes\": {}
        }"

USER_ID=$(curl -X GET "${KEYCLOAK_BASE_URL}/admin/realms/${REALM_NAME}/users" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -G \
    --data-urlencode "username=${KC_USERNAME}" | jq -r '.[0].id')

GROUP_ID=$(curl -X GET "${KEYCLOAK_BASE_URL}/admin/realms/${REALM_NAME}/groups" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" | jq -r --arg GROUP_NAME "vault-client" '.[] | select(.name==$GROUP_NAME) | .id')

curl -X PUT "${KEYCLOAK_BASE_URL}/admin/realms/${REALM_NAME}/users/${USER_ID}/groups/${GROUP_ID}" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json"

# create mapper
echo "Creating mapper..."
curl -X POST "${KEYCLOAK_BASE_URL}/admin/realms/${REALM_NAME}/clients/${CLIENT_ID}/protocol-mappers/models" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
          \"name\": \"${MAPPER_NAME}\",
          \"protocol\": \"openid-connect\",
          \"protocolMapper\": \"oidc-group-membership-mapper\",
          \"config\": {
            \"full.path\": \"false\",
            \"id.token.claim\": \"true\",
            \"access.token.claim\": \"true\",
            \"claim.name\": \"${MAPPER_NAME}\",
            \"userinfo.token.claim\": \"true\"
          }
        }"


