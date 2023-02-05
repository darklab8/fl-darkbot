provider "helm" {
  kubernetes {
    config_path = "~/.kube/prod_darklab_config"
  }
}


module "darkbot" {
  source = "../modules/darkbot"

  environment         = "prod"
  CONSOLER_PREFIX     = "."
  SCRAPPY_PLAYER_URL  = var.SCRAPPY_PLAYER_URL
  SCRAPPY_BASE_URL    = var.SCRAPPY_BASE_URL
  DISCORDER_BOT_TOKEN = var.PRODUCTION_DISCORDER_BOT_TOKEN
  DARKBOT_VERSION     = "v0.1.1"
}
