storage "raft" {
  path    = "./vault-data2"
  node_id = "node2"
}

listener "tcp" {
  address     = "0.0.0.0:8400"
  tls_disable = "true" 
}

disable_mlock = true 
api_addr = "http://127.0.0.1:8400"
cluster_addr = "https://127.0.0.1:8401"
ui = true