storage "raft" {
  path    = "./vault"
  node_id = "node5"
}

listener "tcp" {
  address     = "127.0.0.1:8200"
  tls_disable = "true" 
}

disable_mlock = true 
api_addr = "http://127.0.0.1:8200"
cluster_addr = "https://127.0.0.1:8201"
ui = true
