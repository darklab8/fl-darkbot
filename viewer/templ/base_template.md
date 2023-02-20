{{.Header}} (last updated: {{.LastUpdated}})
**Bases:**
{{range $val := .Bases -}}
```scss
BaseName({{$val.Name | printf "%q"}});
Health({{$val.Health | printf "\"%.1f\""}});
HealthChange({{$val.HealthChange | printf "\"%.6f\""}});
Affiliation({{$val.Affiliation | printf "%q"}});
{{- if $val.IsHealthDecreasing -}}{{ $val.HealthDecreasePhrase }}{{- end }}
{{- if $val.IsUnderAttack -}}{{ $val.UnderAttackPhrase }}{{- end }}
```
{{ end }}
