
locals {
  chart_path = "${path.module}/../charts/darkbot"
  # This hash forces Terraform to redeploy if a new template file is added or changed, or values are updated
  chart_hash  = sha1(join("", [for f in fileset(local.chart_path, "**/*ml") : filesha1("${local.chart_path}/${f}")]))
  environment = var.environment
}

resource "helm_release" "experiment" {
  name             = "darkbot"
  chart            = "../charts/darkbot"
  create_namespace = true
  namespace        = "darkbot-${var.environment}"
  force_update     = false
  reset_values     = true
  recreate_pods    = true

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
    value = var.DISCORDER_BOT_TOKEN
  }

  set {
    name  = "CONFIGURATOR_DBNAME"
    value = var.environment
  }

  set {
    name  = "ENVIRONMENT"
    value = var.environment
  }

  set {
    name  = "HOSTNAME"
    value = "${var.environment}-cluster"
  }

  set {
    name  = "CONSOLER_PREFIX"
    value = var.CONSOLER_PREFIX
  }
}
