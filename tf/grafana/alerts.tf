locals {
    alert_disable_provenance = true
}

resource "grafana_rule_group" "rule_group_15cad967beae7328" {
  org_id             = 1
  name               = "darkbot_uptime_health"
  folder_uid         = "fei55y3yi1340d"
  interval_seconds   = 60
  disable_provenance = local.alert_disable_provenance


  rule {
    name      = "darkbot uptime health"
    condition = "C"

    data {
      ref_id = "A"

      relative_time_range {
        from = 600
        to   = 0
      }

      datasource_uid = "prometheus-datasource"
      model          = "{\"disableTextWrap\":false,\"editorMode\":\"builder\",\"expr\":\"sum by(environment) (darkbot_uptime_seconds)\",\"fullMetaSearch\":false,\"includeNullMetadata\":true,\"instant\":true,\"intervalMs\":1000,\"legendFormat\":\"__auto\",\"maxDataPoints\":43200,\"range\":false,\"refId\":\"A\",\"useBackend\":false}"
    }
    data {
      ref_id = "C"

      relative_time_range {
        from = 0
        to   = 0
      }

      datasource_uid = "__expr__"
      model          = "{\"conditions\":[{\"evaluator\":{\"params\":[60],\"type\":\"lt\"},\"operator\":{\"type\":\"and\"},\"query\":{\"params\":[\"C\"]},\"reducer\":{\"params\":[],\"type\":\"last\"},\"type\":\"query\"}],\"datasource\":{\"type\":\"__expr__\",\"uid\":\"__expr__\"},\"expression\":\"A\",\"intervalMs\":1000,\"maxDataPoints\":43200,\"refId\":\"C\",\"type\":\"threshold\"}"
    }

    no_data_state  = "NoData"
    exec_err_state = "Error"
    for            = "2m"
    labels = {
      ping_channel = "true"
    }
    is_paused = false
  }
}
