locals {
  image_name = "darkwind8/darkbot:${var.tag_version}"
}

locals {
  envs = {
    SCRAPPY_PLAYER_URL    = "${var.secrets["SCRAPPY_PLAYER_URL"]}"
    SCRAPPY_BASE_URL      = "${var.secrets["SCRAPPY_BASE_URL"]}"
    DISCORDER_BOT_TOKEN   = "${var.secrets["DISCORDER_BOT_TOKEN"]}"
    CONFIGURATOR_DBNAME   = "${var.configurator_dbname}"
    CONSOLER_PREFIX       = "${var.consoler_prefix}"
    SCRAPPY_LOOP_DELAY    = "60"
    VIEWER_LOOP_DELAY     = "10"
    DEVENV_MOCK_API       = "false" # legacy. Delete when will deploy new one
    DEV_ENV_MOCK_API      = "false"
    DARKBOT_LOG_LEVEL     = "${var.debug ? "DEBUG" : "WARN"}"
    TYPELOG_LOG_JSON      = "true"
    UTILS_ENVIRONMENT     = var.environment
    UTILS_USERAGENT       = "darkwind/1.0"
    DARKBOT_DARKSTAT_HOST = "darkstat"
    DARKBOT_DARKSTAT_PORT = "8100"
  }
}
