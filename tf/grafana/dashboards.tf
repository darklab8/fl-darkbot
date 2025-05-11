
locals {
  grafana-dashboards = {
    darkbot_dashboard = {
      json = file("${path.module}/dashboards/darkbot.json")
    }
  }
}

resource "grafana_dashboard" "dashboard" {
  for_each    = local.grafana-dashboards
  config_json = each.value.json
}
