# Create a docker image resource
# -> docker pull nginx:latest
resource "docker_image" "darkbot" {
  name         = "darkwind8/darkbot:${var.tag_version}"
  keep_locally = true
}

variable "tag_version" {
  type = string
}

variable "configurator_dbname" {
  type = string
}

variable "consoler_prefix" {
  type = string
}

variable "debug" {
  type = bool
}

variable "secrets" {
  type = map(string)
}

# # Create a docker container resource
# # -> same as 'docker run --name nginx -p8080:80 -d nginx:latest'
resource "docker_container" "darkbot" {
  name  = "darkbot"
  image = docker_image.darkbot.image_id

  env = [
    "SCRAPPY_PLAYER_URL=${var.secrets["SCRAPPY_PLAYER_URL"]}",
    "SCRAPPY_BASE_URL=${var.secrets["SCRAPPY_BASE_URL"]}",
    "DISCORDER_BOT_TOKEN=${var.secrets["DISCORDER_BOT_TOKEN"]}",
    "CONFIGURATOR_DBNAME=${var.configurator_dbname}",
    "CONSOLER_PREFIX=${var.consoler_prefix}",
    "LOOP_DELAY=60",
    "DEVENV_MOCK_API=false",
    "DARKBOT_LOG_LEVEL=${var.debug ? "DEBUG" : "WARN"}"
  ]

  restart = "always"
  volumes {
    container_path = "/code/data"
    read_only      = false
    host_path      = "/var/lib/darklab/darkbot"
  }
}
