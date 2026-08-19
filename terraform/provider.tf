terraform {
  required_version = ">= 1.6.0"

  required_providers {
    vscale = {
      source  = "vscale/vscale"
      version = "~> 1.0"
    }
  }
}

provider "vscale" {
  token = var.vscale_token
}
