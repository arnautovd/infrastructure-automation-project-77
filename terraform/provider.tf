terraform {
  required_version = ">= 1.6.0"

  required_providers {
    vscale = {
      source  = "arnautovd/vscale"
      version = "0.1.1"
    }
  }
}

provider "vscale" {
  token = var.vscale_token
}
