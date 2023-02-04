provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

module "darkbot" {
  source = "../modules/darkbot"

  environment         = "staging"
  CONSOLER_PREFIX     = "!"
  SCRAPPY_PLAYER_URL  = var.SCRAPPY_PLAYER_URL
  SCRAPPY_BASE_URL    = var.SCRAPPY_BASE_URL
  DISCORDER_BOT_TOKEN = var.STAGING_DISCORDER_BOT_TOKEN
}
