terraform {
  required_providers {
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

data "external" "secrets_cloudflare" {
  program = ["pass", "personal/terraform/cloudflare/dd84ai"]
}

module "data_cluster" {
  source = "../../../infra/tf/production/output/deserializer"
}

provider "cloudflare" {
  api_token = data.external.secrets_cloudflare.result["token"]
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "darkbot-context"
}
