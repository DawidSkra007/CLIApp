# List existing secrets engines.
path "sys/mounts"
{
  capabilities = ["read"]
}

# List existing policies
path "sys/policies/acl"
{
  capabilities = ["list"]
}

# List existing policies
path "sys/policies/acl"
{
  capabilities = ["list"]
}

path "sys/capabilities-self"
{
  capabilities = ["update"]
}

# Create and manage ACL policies
path "sys/policies/acl/*"
{
  capabilities = ["create", "read", "update", "list", "sudo"]
}

# List auth methods
path "sys/auth"
{
  capabilities = ["read"]
}

# Manage secrets engines
path "sys/mounts/*"
{
  capabilities = ["create", "read", "update", "list"]
}

# Perform action on paths
# test path
path "test/data/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "patch"] 
}

path "test/undelete/*" {
  capabilities = ["update"]
}

# kv path
path "kv/data/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "patch"] 
}

path "kv/undelete/*" {
  capabilities = ["update"]
}
