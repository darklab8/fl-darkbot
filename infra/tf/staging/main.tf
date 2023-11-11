module "stack" {
  source       = "../modules/hetzner_server"
  environment  = "staging"
  server_power = "cpx21"
  backups      = false
}

output "cluster_ip" {
  value = module.stack.cluster_ip
}

provider "docker" {
  host     = "ssh://root@${module.stack.cluster_ip}:22"
  ssh_opts = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-i", "~/.ssh/id_rsa.darklab"]
}

data "aws_ssm_parameter" "darkbot" {
  name = "/terraform/hetzner/darkbot/staging"
}

locals {
  secrets = nonsensitive(jsondecode(data.aws_ssm_parameter.darkbot.value))
}

module "darkbot" {
  source              = "../modules/darkbot"
  configurator_dbname = "staging"
  consoler_prefix     = ","
  secrets             = local.secrets
  tag_version         = "v1.0.2-rc.1"
  debug               = true
}

# # Create a docker image resource
# # -> docker pull nginx:latest
# resource "docker_image" "nginx" {
#   name         = "nginx:latest"
#   keep_locally = true
# }

# # Create a docker container resource
# # -> same as 'docker run --name nginx -p8080:80 -d nginx:latest'
# resource "docker_container" "nginx" {
#   name    = "nginx"
#   image   = docker_image.nginx.image_id

#   ports {
#     external = 8080
#     internal = 80
#   }
# }
