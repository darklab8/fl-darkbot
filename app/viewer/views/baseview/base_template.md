```scss
BaseName({{.Name | printf "%q"}});
Health({{derefF64 .Health | printf "\"%.1f\""}});
HealthChangeInLast15m({{.HealthChange | printf "%q"}});
Affiliation({{derefS .FactionName | printf "%q"}});
{{- if (ne (derefI .Money) -1) }}
Money({{derefI .Money | printf "\"%d\""}});
{{- end -}}
{{- if (ne (derefI .CargoSpaceLeft) -1) }}
CargoSpaceLeft({{derefI .CargoSpaceLeft | printf "\"%d\""}});
{{- end -}}
{{- if .IsHealthDecreasing -}}{{ .HealthDecreasePhrase }}{{- end }}
{{- if .IsUnderAttack -}}{{ .UnderAttackPhrase }}{{- end }}
```
