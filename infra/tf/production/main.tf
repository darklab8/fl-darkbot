terraform {
  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "~> 1.35.2"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">= 2.7.0"
    }
  }
}

variable "production_hcloud_token" {
  type      = string
  sensitive = true
}

provider "hcloud" {
  token = var.production_hcloud_token
}

module "stack" {
  source       = "../modules/darkbot"
  environment  = "prod"
  server_power = "cpx21"
}

output "cluster_ip" {
  value = module.stack.cluster_ip
}
