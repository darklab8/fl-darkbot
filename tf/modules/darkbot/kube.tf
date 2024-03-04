locals {
    namespace = "darkbot-${var.environment}"
}

resource "kubernetes_namespace_v1" "example" {
  count = var.mode == "kubernetes" ? 1 : 0
  metadata {
    name = local.namespace
  }
}

resource "kubernetes_pod_v1" "test" {
  count = var.mode == "kubernetes" ? 1 : 0
  metadata {
    name      = "darkbot"
    namespace = local.namespace
  }

  spec {
    affinity {
      node_affinity {
        required_during_scheduling_ignored_during_execution {
          node_selector_term {
            match_expressions {
              key      = "node"
              operator = "In"
              values   = ["arm"]
            }
          }
        }
      }
    }

    restart_policy = "Always"
    container {
      image = local.image_name
      name  = "app"

      volume_mount {
        name       = "volv"
        mount_path = "/code/data"
        read_only  = false
      }

      dynamic "env" {
        for_each = local.envs
        content {
          name  = env.key
          value = env.value
        }
      }
    }
    volume {
      name = "volv"
      host_path {
        path = "/var/lib/darklab/darkbot-${var.environment}"
        type = "DirectoryOrCreate"
      }
    }
  }

  lifecycle {
    ignore_changes = [
      metadata,
    ]
  }
}
