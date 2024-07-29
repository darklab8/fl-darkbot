terraform {
  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = ">=1.45.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">=2.7.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = ">=3.7.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = ">=3.0.2"
    }
  }
}

data "external" "secrets_cloudflare" {
  program = ["pass", "api/personal/terraform/cloudflare/dd84ai"]
}

data "external" "secrets_hetzner" {
  program = ["pass", "api/personal/terraform/hetzner/production"]
}

provider "hcloud" {
  token = data.external.secrets_hetzner.result["token"]
}

provider "cloudflare" {
  api_token = data.external.secrets_cloudflare.result["token"]
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "darklab"
}
