{{.Header}} (last updated: {{.LastUpdated}})
**Bases:**
{{range $val := .Bases -}}
```scss
BaseName({{$val.Name | printf "%q"}});
Health({{$val.Health | printf "\"%.1f\""}});
Affiliation({{$val.Affiliation | printf "%q"}});
```
{{ end -}}

{{- /* How it should be looking. Based on Darkbot1.X   */ -}}
{{- /* **Bases:** `Health` - `Name` - `Affiliation`   */ -}}
{{- /* ```json   */ -}}
{{- /* {% for key, value in data.items() %}   */ -}}
{{- /* {"{{ value["diff"] }}"}["{{value["health"]}}"]{"{{ key }}"}["{{value["affiliation"]}}"]{% endfor %}   */ -}}
{{- /* ``` -->   */ -}}