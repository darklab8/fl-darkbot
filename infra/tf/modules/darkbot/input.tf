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
  type    = bool
  default = false
}

variable "secrets" {
  type = map(string)
}

variable "mode" {
  type = string

  validation {
    condition     = contains(["kubernetes", "docker"], var.mode)
    error_message = "Valid mode. should be docker or kubernetes"
  } 
}
