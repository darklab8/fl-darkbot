locals {
  datacenter  = "ash-dc1" # USA
  image       = "ubuntu-22.04"
  server_type = "cpx21"
  task_name   = "cluster"
}

resource "hcloud_ssh_key" "darklab" {
  name       = "darklab_key"
  public_key = file("../../../id_rsa.darklab.pub")
}

resource "hcloud_server" "cluster" {
  name        = "${var.envornment}-cluster"
  image       = local.image
  datacenter  = local.datacenter
  server_type = local.server_type
  ssh_keys = [
    hcloud_ssh_key.darklab.id,
  ]
  public_net {
    ipv4_enabled = true
    ipv6_enabled = true
  }
}
