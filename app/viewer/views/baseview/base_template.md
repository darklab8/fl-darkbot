```scss
BaseName({{.Name | printf "%q"}});
Health({{.Health | printf "\"%.1f\""}});
HealthChangeInLast15m({{.HealthChange | printf "%q"}});
Affiliation({{.FactionName | printf "%q"}});
{{- if .IsHealthDecreasing -}}{{ .HealthDecreasePhrase }}{{- end }}
{{- if .IsUnderAttack -}}{{ .UnderAttackPhrase }}{{- end }}
```
