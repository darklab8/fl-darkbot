locals {
  datacenter  = "ash-dc1" # USA
  image       = "ubuntu-22.04"
  server_type = var.server_power
  task_name   = "cluster"
}

resource "hcloud_ssh_key" "darklab" {
  name       = "darklab_key"
  public_key = file("${path.module}/id_rsa.darklab.pub")
}

data "hcloud_image" "default" {
  name = "docker-ce"
}

resource "hcloud_server" "cluster" {
  name        = "${var.environment}-cluster"
  image       = data.hcloud_image.default.id
  datacenter  = local.datacenter
  server_type = local.server_type
  ssh_keys = [
    hcloud_ssh_key.darklab.id,
  ]
  public_net {
    ipv4_enabled = true
    ipv6_enabled = true
  }

  lifecycle {
    ignore_changes = [
      image,
    ]
  }

  backups = var.backups
}

output "cluster_ip" {
  value = hcloud_server.cluster.ipv4_address
}
