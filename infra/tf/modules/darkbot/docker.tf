# Create a docker image resource
# -> docker pull nginx:latest
resource "docker_image" "darkbot" {
  count = var.mode == "docker" ? 1 : 0
  name         = "darkwind8/darkbot:${var.tag_version}"
  keep_locally = true
}

# # Create a docker container resource
# # -> same as 'docker run --name nginx -p8080:80 -d nginx:latest'
resource "docker_container" "darkbot" {
  count = var.mode == "docker" ? 1 : 0

  name  = "darkbot"
  image = docker_image.darkbot[0].image_id

  env = [
    "SCRAPPY_PLAYER_URL=${var.secrets["SCRAPPY_PLAYER_URL"]}",
    "SCRAPPY_BASE_URL=${var.secrets["SCRAPPY_BASE_URL"]}",
    "DISCORDER_BOT_TOKEN=${var.secrets["DISCORDER_BOT_TOKEN"]}",
    "CONFIGURATOR_DBNAME=${var.configurator_dbname}",
    "CONSOLER_PREFIX=${var.consoler_prefix}",
    "SCRAPPY_LOOP_DELAY=60",
    "VIEWER_LOOP_DELAY=10",
    "DEVENV_MOCK_API=false",
    "DARKBOT_LOG_LEVEL=${var.debug ? "DEBUG" : "WARN"}"
  ]

  restart = "always"
  volumes {
    container_path = "/code/data"
    read_only      = false
    host_path      = "/var/lib/darklab/darkbot"
  }

  memory = 1000 # MBs

  lifecycle {
    ignore_changes = [
      memory_swap,
    ]
  }
}
