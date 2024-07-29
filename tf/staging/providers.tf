terraform {
  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = ">=1.35.2"
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
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

data "external" "pass_storage" {
  program = ["bash", "${path.module}/pass_storage.sh"]
}

provider "hcloud" {
  token = data.external.pass_storage.result["hetzner_token"]
}

provider "cloudflare" {
  api_token = data.external.pass_storage.result["cloudflare_token"]
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "darklab"
}
