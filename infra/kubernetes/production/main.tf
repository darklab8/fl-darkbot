provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}


locals {
  chart_path = "${path.module}/../charts/darkbot"
  # This hash forces Terraform to redeploy if a new template file is added or changed, or values are updated
  chart_hash  = sha1(join("", [for f in fileset(local.chart_path, "**/*ml") : filesha1("${local.chart_path}/${f}")]))
  environment = "prod"
}

resource "helm_release" "experiment" {
  name             = "darkbot"
  chart            = "../charts/darkbot"
  create_namespace = true
  namespace        = "darkbot-${local.environment}"
  # force_update     = true
  # reset_values     = true
  # recreate_pods = true

  set {
    name  = "chartHash"
    value = local.chart_hash
  }

  set_sensitive {
    name  = "SCRAPPY_PLAYER_URL"
    value = var.SCRAPPY_PLAYER_URL
  }

  set_sensitive {
    name  = "SCRAPPY_BASE_URL"
    value = var.SCRAPPY_BASE_URL
  }

  set_sensitive {
    name  = "DISCORDER_BOT_TOKEN"
    value = var.PRODUCTION_DISCORDER_BOT_TOKEN
  }

  set {
    name  = "PRODUCTION_CONFIGURATOR_DBNAME"
    value = "prod"
  }

  set {
    name  = "ENVIRONMENT"
    value = "prod"
  }

  set {
    name  = "HOSTNAME"
    value = "production-cluster"
  }
}
