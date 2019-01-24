# {{ .checkId }} Table sizes #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
Table | Rows | &#9660;&nbsp;Total size | Table size | Index(es) Size | TOAST Size
------|------|------------|------------|----------------|------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
{{ index $value "Table"}} | {{ index $value "Rows"}} | {{ index $value "Total Size"}} | {{ index $value "Table Size"}} | {{ index $value "Index(es) Size"}} | {{ index $value "TOAST Size"}}
{{ end }}
{{- else }}
No data
{{- end }}

## Conclusions ##


## Recommendations ##

