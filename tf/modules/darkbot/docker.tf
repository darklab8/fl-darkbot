resource "docker_image" "darkbot" {
  count        = var.mode == "docker" ? 1 : 0
  name         = local.image_name
  keep_locally = true
}

resource "docker_container" "darkbot" {
  count = var.mode == "docker" ? 1 : 0

  name  = "darkbot-${var.environment}"
  image = docker_image.darkbot[0].image_id

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
