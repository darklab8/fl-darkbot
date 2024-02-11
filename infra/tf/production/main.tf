module "stack" {
  source       = "../modules/hetzner_server"
  environment  = "production"
  server_power = "cpx21"
  backups      = false
}

output "cluster_ip" {
  value = module.stack.cluster_ip
}

data "aws_ssm_parameter" "github.com/darklab8/fl-darkbot" {
  name = "/terraform/hetzner/darkbot/production"
}

locals {
  secrets = nonsensitive(jsondecode(data.aws_ssm_parameter.darkbot.value))
}

provider "docker" {
  host     = "ssh://root@${module.stack.cluster_ip}:22"
  ssh_opts = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-i", "~/.ssh/id_rsa.darklab"]
}

module "github.com/darklab8/fl-darkbot" {
  source              = "../modules/darkbot"
  configurator_dbname = "production"
  consoler_prefix     = "."
  secrets             = local.secrets
  tag_version         = "v1.5.1"
}