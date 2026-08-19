variable "vscale_token" {
  description = "Vscale API token. Pass it through TF_VAR_vscale_token or VSCALE_API_TOKEN."
  type        = string
  sensitive   = true
}

variable "web_server_count" {
  description = "Number of web servers."
  type        = number
  default     = 2

  validation {
    condition     = var.web_server_count == 2
    error_message = "This project requires exactly two web servers."
  }
}

variable "web_server_image" {
  description = "Vscale image identifier used to create each web server."
  type        = string
  default     = "debian_12_64"
}

variable "web_server_plan" {
  description = "Vscale resource plan identifier used by each web server."
  type        = string
  default     = "small"
}

variable "vscale_location" {
  description = "Vscale location identifier."
  type        = string
  default     = "msk0"
}

variable "vscale_ssh_key_ids" {
  description = "IDs of SSH keys already registered in Vscale."
  type        = list(string)
  default     = []
}
