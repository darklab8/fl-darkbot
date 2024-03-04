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

data "aws_ssm_parameter" "hetzner" {
  name = "/terraform/hetzner/production"
}
provider "hcloud" {
  token = data.aws_ssm_parameter.hetzner.value
}

data "aws_ssm_parameter" "cloudflare_key" {
  name = "/terraform/cloudflare/dd84ai"
}

provider "cloudflare" {
  api_token = data.aws_ssm_parameter.cloudflare_key.value
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "darklab"
}
