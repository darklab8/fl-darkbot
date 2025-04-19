resource "docker_image" "darkbot" {
  count        = var.mode == "docker" || var.mode == "swarm" ? 1 : 0
  name         = local.image_name
  keep_locally = true
}

locals {
  service_name = "${var.environment}-darkbot-app"
}

resource "docker_network" "darkbot" {
  name       = "darkbot-${var.environment}"
  driver     = "overlay"
  attachable = true
}

data "docker_network" "grafana" {
  name = "grafana"
}

resource "docker_container" "darkbot" {
  count = var.mode == "docker" ? 1 : 0

  name  = local.service_name
  image = docker_image.darkbot[0].name

  networks_advanced {
    name    = docker_network.darkbot.id
    aliases = ["darkbot"]
  }

  networks_advanced {
    name    = data.docker_network.grafana.id
    aliases = ["${var.environment}-darkbot"]
  }

  env = [for k, v in local.envs : "${k}=${v}"]

  restart = "always"
  volumes {
    container_path = "/code/data"
    read_only      = false
    host_path      = "/var/lib/darklab/darkbot-${var.environment}"
  }

  volumes {
    container_path = "/tmp/darkstat"
    read_only      = false
    host_path      = "/tmp/darkstat-${var.environment}"
  }


  labels {
    label = "prometheus"
    value = "true"
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

  name = local.service_name

  task_spec {
    networks_advanced {
      name    = docker_network.darkbot.id
      aliases = ["darkbot"]
    }

    networks_advanced {
      name    = data.docker_network.grafana.id
      aliases = ["${var.environment}-darkbot"]
    }

    container_spec {
      image = docker_image.darkbot[0].name
      env   = local.envs

      labels {
        label = "prometheus"
        value = "true"
      }

      mounts {
        target    = "/code/data"
        source    = "/var/lib/darklab/darkbot-${var.environment}"
        type      = "bind"
        read_only = false

        bind_options {
          propagation = "rprivate"
        }
      }
      mounts { // darkstat socks
        target    = "/tmp/darkstat"
        source    = "/tmp/darkstat-${var.environment}"
        type      = "bind"
        read_only = false
        bind_options {
          propagation = "rprivate"
        }
      }
    }
    restart_policy {
      condition = "any"
      delay     = "5m"
      window    = "10s"
    }
    resources {
      limits {
        memory_bytes = 1000 * 1000 * 1000 # 1 gb
      }
    }
  }

  lifecycle {
    ignore_changes = [
      task_spec[0].restart_policy[0].window,
      task_spec[0].container_spec[0].image,
    ]
  }
}

