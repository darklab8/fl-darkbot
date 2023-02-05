
locals {
  chart_path = "${path.module}/../../charts/darkbot"
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

  values = [
    <<-EOT
    hard_memory_limit: "${var.limit.hard_memory}"
    hard_cpu_limit: "${var.limit.hard_cpu}"
    hostname: "${var.environment}-cluster"
    image_version: "${var.image_version}"

    ENVIRONMENT: "${var.environment}"
    CONFIGURATOR_DBNAME: "${var.environment}"
    CONSOLER_PREFIX: "${var.environ.CONSOLER_PREFIX}"
    LOGGING: "${var.environ.LOGGING}"
    LOOP_DELAY: "${var.environ.LOOP_DELAY}"
    EOT
  ]
  set {
    name  = "chartHash"
    value = local.chart_hash
  }

  set_sensitive {
    name  = "SCRAPPY_PLAYER_URL"
    value = var.environ.SCRAPPY_PLAYER_URL
  }

  set_sensitive {
    name  = "SCRAPPY_BASE_URL"
    value = var.environ.SCRAPPY_BASE_URL
  }

  set_sensitive {
    name  = "DISCORDER_BOT_TOKEN"
    value = var.environ.DISCORDER_BOT_TOKEN
  }
}
