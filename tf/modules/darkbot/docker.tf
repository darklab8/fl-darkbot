resource "docker_image" "darkbot" {
  count        = var.mode == "docker" || var.mode == "swarm" ? 1 : 0
  name         = local.image_name
  keep_locally = true
}

resource "docker_container" "darkbot" {
  count = var.mode == "docker" ? 1 : 0

  name  = "darkbot-${var.environment}"
  image = docker_image.darkbot[0].name

  env = [for k, v in local.envs : "${k}=${v}"]

  restart = "always"
  volumes {
    container_path = "/code/data"
    read_only      = false
    host_path      = "/var/lib/darklab/darkbot-${var.environment}"
  }

  memory = 1000 # MBs

  lifecycle {
    ignore_changes = [
      memory_swap,
      network_mode,
    ]
  }
}

resource "docker_service" "darkbot" {
  count = var.mode == "swarm" ? 1 : 0

  name = "darkbot-${var.environment}"

  task_spec {
    container_spec {
      image = docker_image.darkbot[0].name
      env   = local.envs

      mounts {
        target    = "/code/data"
        source    = "/var/lib/darklab/darkbot-${var.environment}"
        type      = "bind"
        read_only = false

        bind_options {
          propagation = "rprivate"
        }
      }
    }
    restart_policy {
      condition = "any"
      delay     = "20s"
    }
    resources {
      limits {
        memory_bytes = 1000 * 1000 * 1000 # 1 gb
      }
    }
  }
}

