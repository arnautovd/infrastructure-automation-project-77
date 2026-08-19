output "web_server_names" {
  description = "Names of the web servers."
  value       = vscale_scalet.web[*].name
}
