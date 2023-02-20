{{.Header}} (last updated: {{.LastUpdated}})
**Bases:**
{{range $val := .Bases -}}
'''scss
BaseName({{$val.Name | printf "%q"}});
Health({{$val.Health | printf "\"%.1f\""}});
HealthChange({{$val.HealthChange | printf "\"%.6f\""}});
Affiliation({{$val.Affiliation | printf "%q"}});
{{ $val.IsHealthDecreasing }}
{{ $val.isUnderAttack }}
'''
{{ end }}
