module "ssh_key" {
   source       = "../../../../infra/tf/modules/hetzner_ssh_key/data"
}

module "stack" {
  source       = "../../../../infra/tf/modules/hetzner_server"
  environment  = "production"
  name         = "node-darkbot"
  hardware     = "cpx21"
  ssh_key_id   = module.ssh_key.id
  datacenter   = "ash-dc1"
}


data "aws_ssm_parameter" "darkbot" {
  name = "/terraform/hetzner/darkbot/production"
}

locals {
  secrets = nonsensitive(jsondecode(data.aws_ssm_parameter.darkbot.value))
}

provider "docker" {
  host     = "ssh://root@${module.stack.ipv4_address}:22"
  ssh_opts = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-i", "~/.ssh/id_rsa.darklab"]
}

module "darkbot" {
  source              = "../modules/darkbot"
  configurator_dbname = "production"
  consoler_prefix     = "."
  secrets             = local.secrets
  tag_version         = "v1.5.1"
  mode                = "docker"
}