{{.Header}} (last updated: {{.LastUpdated}})
**Bases:** `Health` - `Name` - `Affiliation`
```{{ if .Bases }}arm{{ end }}
{{range $val := .Bases -}}
health: {{$val.Health | printf "\"%.1f\""}}, #name: {{$val.Name | printf "%q"}}, #affiliation: {{$val.Affiliation | printf "%q"}}
{{ end -}}
```

{{- /* How it should be looking. Based on Darkbot1.X   */ -}}
{{- /* **Bases:** `Health` - `Name` - `Affiliation`   */ -}}
{{- /* ```json   */ -}}
{{- /* {% for key, value in data.items() %}   */ -}}
{{- /* {"{{ value["diff"] }}"}["{{value["health"]}}"]{"{{ key }}"}["{{value["affiliation"]}}"]{% endfor %}   */ -}}
{{- /* ``` -->   */ -}}