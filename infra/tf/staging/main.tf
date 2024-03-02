# module "stack" {
#   source       = "../modules/hetzner_server"
#   environment  = "staging"
#   server_power = "cpx21"
#   backups      = false
# }

# output "cluster_ip" {
#   value = module.stack.cluster_ip
# }

# provider "docker" {
#   host     = "ssh://root@${module.stack.cluster_ip}:22"
#   ssh_opts = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-i", "~/.ssh/id_rsa.darklab"]
# }

data "aws_ssm_parameter" "darkbot" {
  name = "/terraform/hetzner/darkbot/staging"
}

locals {
  secrets = nonsensitive(jsondecode(data.aws_ssm_parameter.darkbot.value))
}

# module "darkbot" {
#   source              = "../modules/darkbot"
#   configurator_dbname = "staging"
#   consoler_prefix     = ","
#   secrets             = local.secrets
#   tag_version         = "v1.5.0"
#   debug               = false
# }
