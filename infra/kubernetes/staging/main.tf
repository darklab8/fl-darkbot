provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}


locals {
  chart_path = "${path.module}/../charts/darkbot"
  # This hash forces Terraform to redeploy if a new template file is added or changed, or values are updated
  chart_hash  = sha1(join("", [for f in fileset(local.chart_path, "**/*ml") : filesha1("${local.chart_path}/${f}")]))
  environment = "staging"
}

resource "helm_release" "experiment" {
  name             = "darkbot"
  chart            = "../charts/darkbot"
  create_namespace = true
  namespace        = "darkbot-${local.environment}"

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
    value = var.STAGING_DISCORDER_BOT_TOKEN
  }

  set {
    name  = "CONFIGURATOR_DBNAME"
    value = "dev"
  }

  set {
    name  = "ENVIRONMENT"
    value = "staging"
  }

  set {
    name  = "HOSTNAME"
    value = "staging-cluster"
  }

  set {
    name  = "CONSOLER_PREFIX"
    value = ","
  }
}
