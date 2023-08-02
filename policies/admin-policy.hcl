# Read system health check
path "sys/health"
{
  capabilities = ["read", "sudo"]
}

# Create and manage ACL policies broadly across Vault

# List existing policies
path "sys/policies/acl"
{
  capabilities = ["list"]
}

# Create and manage ACL policies
path "sys/policies/acl/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "sys/capabilities-self"
{
  capabilities = ["update"]
}

# Enable and manage authentication methods broadly across Vault

# Manage auth methods broadly across Vault
path "auth/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
# Create, update, and delete auth methods
path "sys/auth/*"
{
  capabilities = ["create", "update", "delete", "sudo"]
}

# List auth methods
path "sys/auth"
{
  capabilities = ["read"]
}

# Enable and manage the key/value secrets engine at `secret/` path

# List, create, update, and delete key/value secrets
path "secret/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo", "patch"]
}

# Manage secrets engines
path "sys/mounts/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List existing secrets engines.
path "sys/mounts"
{
  capabilities = ["read"]
}

# Perform action on paths ** ALL PATHS NEED TO BE SPECIFIED **
path "test/data/*"
{
  capabilities = ["create", "read", "update", "delete", "patch"]
}

path "test/destroy/*" {
  capabilities = ["update"]
}

path "test/metadata/*" {
  capabilities = ["list", "read", "delete"]
}

path "test/undelete/*" {
  capabilities = ["update"]
}

path "kv/data/*"
{
  capabilities = ["create", "read", "update", "delete", "patch"]
}

path "kv/destroy/*" {
  capabilities = ["update"]
}

path "kv/metadata/*" {
  capabilities = ["list", "read", "delete"]
}

path "kv/undelete/*" {
  capabilities = ["update"]
}
