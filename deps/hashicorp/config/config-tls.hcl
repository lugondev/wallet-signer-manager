backend "file" {
  path = "/vault/file"
}

listener "tcp" {
  address = "hashicorp:8200"
  tls_disable = false
  tls_client_ca_file = "/vault/tls/ca.crt"
  tls_cert_file = "/vault/tls/tls.crt"
  tls_key_file = "/vault/tls/tls.key"
}

default_lease_ttl = "15m"
max_lease_ttl = 99999999
api_addr = "https://hashicorp:8200"
plugin_directory = "/vault/plugins"
log_level = "Debug"

ui = false
