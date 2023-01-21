terraform {
  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "~> 1.35.2"
    }
  }
}

variable "staging_hcloud_token" {
  type      = string
  sensitive = true
}

provider "hcloud" {
  token = var.staging_hcloud_token
}

module "stack" {
  source      = "../modules/darkbot"
  environment = "production"
}
