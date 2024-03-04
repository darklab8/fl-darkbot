module "server" {
  source = "../../../../infra/tf/modules/hetzner_server/data"
  name   = "node-arm"
}

data "aws_ssm_parameter" "darkbot" {
  name = "/terraform/hetzner/darkbot/staging"
}

locals {
  secrets = nonsensitive(jsondecode(data.aws_ssm_parameter.darkbot.value))
}

provider "docker" {
  host     = "ssh://root@${module.server.ipv4_address}:22"
  ssh_opts = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-i", "~/.ssh/id_rsa.darklab"]
}

module "darkbot" {
  source              = "../modules/darkbot"
  configurator_dbname = "staging"
  consoler_prefix     = ","
  secrets             = local.secrets
  tag_version         = "v1.5.1-arm"
  debug               = false
  mode                = "kubernetes"
  environment         = "staging"
}
