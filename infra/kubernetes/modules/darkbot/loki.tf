resource "helm_release" "loki" {
  name             = "loki"
  chart            = "../charts/loki"
  create_namespace = true
  namespace        = "loki"

  values = [
    "${file("${path.module}/../../charts/loki/values.yaml")}"
  ]
}
