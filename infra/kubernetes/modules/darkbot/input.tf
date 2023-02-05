variable "environ" {
  type = object({
    SCRAPPY_PLAYER_URL  = string
    SCRAPPY_BASE_URL    = string
    DISCORDER_BOT_TOKEN = string
    CONSOLER_PREFIX     = string
    LOGGING             = bool
    LOOP_DELAY          = number
  })
}

variable "environment" {
  type = string
}

variable "image_version" {
  type        = string
  description = "darkbot image version"
}

variable "limit" {
  type = object({
    hard_memory = string
    hard_cpu    = string
  })
}
