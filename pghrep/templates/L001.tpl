# {{ .checkId }} Table sizes #

## Observations ##
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
### Master (`{{.hosts.master}}`) ###
Table | Rows | &#9660;&nbsp;Total size | Table size | Index(es) Size | TOAST Size
------|------|------------|------------|----------------|------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
{{ index $value "Table"}} | {{ index $value "Rows"}} | {{ index $value "Total Size"}} | {{ index $value "Table Size"}} | {{ index $value "Index(es) Size"}} | {{ index $value "TOAST Size"}}
{{ end }}
{{- else -}}{{/*Master data*/}}
No data
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master*/}}
No data
{{ end }}{{/*Master*/}}

## Conclusions ##


## Recommendations ##

