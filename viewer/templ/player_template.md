**{{.Header}}** (last updated: {{.LastUpdated}})
{{.TableName}}
{{- if and (.Players) (contains .TableName "Neutral") }}```json
{{range $val := .Players -}}
name: {{$val.Name | printf "%q"}}, #system: {{$val.System | printf "%q"}}, #region: {{$val.Region | printf "%q"}}, #time: {{$val.Time | printf "%q"}}
{{ end -}}
```{{ end -}}
{{- if and (.Players) (contains .TableName "Friend") }}```diff
{{range $val := .Players -}}
+name: {{$val.Name | printf "%q"}}, #system: {{$val.System | printf "%q"}}, #region: {{$val.Region | printf "%q"}}, #time: {{$val.Time | printf "%q"}}
{{ end -}}
```{{ end -}}
{{- if and (.Players) (contains .TableName "Enemy") }}```diff
{{range $val := .Players -}}
-name: {{$val.Name | printf "%q"}}, #system: {{$val.System | printf "%q"}}, #region: {{$val.Region | printf "%q"}}, #time: {{$val.Time | printf "%q"}}
{{ end -}}
```{{ end -}}