variable "vscale_token" {
  description = "Vscale API token. Pass it through TF_VAR_vscale_token or VSCALE_API_TOKEN."
  type        = string
  sensitive   = true
}
