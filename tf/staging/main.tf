module "server" {
  source = "../../../infra/tf/modules/hetzner_server/data"
  name   = "node-darkbot"
}

data "external" "secrets_darkbot" {
  program = ["pass", "personal/terraform/hetzner/darkbot/staging"]
}

locals {
  secrets = nonsensitive(data.external.secrets_darkbot.result)
}

provider "docker" {
  host     = "ssh://root@${module.server.ipv4_address}:22"
  ssh_opts = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-i", "~/.ssh/id_rsa.darklab"]
}

# limitation of `redock` usage with local terraform state. Repair container if necessary.
# Have remote accessable terraform state if u wish it being resolvable from CI automatically / or utilize docker swarm hmm
# cd tf/staging
# tofu state rm module.darkbot.docker_container.darkbot[0]
# export container_id=$(DOCKER_HOST=ssh://root@darkbot docker inspect --format="{{.Id}}" darkbot-staging)
# tofu import module.darkbot.docker_container.darkbot[0] $container_id

module "darkbot" {
  source              = "../modules/darkbot"
  configurator_dbname = "staging"
  consoler_prefix     = ","
  secrets             = local.secrets
  tag_version         = "staging"
  debug               = false
  mode                = "docker"
  environment         = "staging"
}
