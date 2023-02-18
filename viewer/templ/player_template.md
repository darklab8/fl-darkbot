**{{.Header}}** (last updated: {{.LastUpdated}})
{{.TableName}}
{{- if contains .TableName "Neutral" }}```json
{{range $val := .Players -}}
name: {{$val.Name | printf "%q"}}, #system: {{$val.System | printf "%q"}}, #region: {{$val.Region | printf "%q"}}, #time: {{$val.Time | printf "%q"}}
{{ end -}}
```{{ end -}}
{{- if contains .TableName "Friend" }}```diff
{{range $val := .Players -}}
+name: {{$val.Name | printf "%q"}}, #system: {{$val.System | printf "%q"}}, #region: {{$val.Region | printf "%q"}}, #time: {{$val.Time | printf "%q"}}
{{ end -}}
```{{ end -}}
{{- if contains .TableName "Enemy" }}```diff
{{range $val := .Players -}}
-name: {{$val.Name | printf "%q"}}, #system: {{$val.System | printf "%q"}}, #region: {{$val.Region | printf "%q"}}, #time: {{$val.Time | printf "%q"}}
{{ end -}}
```{{ end -}}