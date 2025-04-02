data "external" "secrets_darkbot" {
  program = ["pass", "personal/terraform/hetzner/darkbot/production"]
}

locals {
  secrets = nonsensitive(data.external.secrets_darkbot.result)
}

provider "docker" {
  host     = "ssh://root@${module.data_cluster.node_darklab.ipv4_address}:22"
  ssh_opts = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-i", "~/.ssh/id_rsa.darklab"]
}

module "darkbot" {
  source              = "../modules/darkbot"
  configurator_dbname = "production"
  consoler_prefix     = "."
  secrets             = local.secrets
  tag_version         = "production"
  mode                = "swarm"
  environment         = "production"
  debug               = false
}
