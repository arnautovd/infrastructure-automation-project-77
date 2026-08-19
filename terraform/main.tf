resource "vscale_scalet" "web" {
  count = var.web_server_count

  name      = format("web-%02d", count.index + 1)
  make_from = var.web_server_image
  rplan     = var.web_server_plan
  location  = var.vscale_location
  ssh_keys  = var.vscale_ssh_key_ids
}

output "web_server_ids" {
  description = "IDs of the two web servers."
  value       = vscale_scalet.web[*].id
}

output "web_server_addresses" {
  description = "Public IPv4 addresses of the two web servers."
  value       = vscale_scalet.web[*].public_address
}
