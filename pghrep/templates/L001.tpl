# {{ .checkId }} Table Sizes #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
{{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
{{ if gt (len (index (index .results .hosts.master) "data")) .LISTLIMIT }}The list is limited to {{.LISTLIMIT}} items.{{ end }}  

| \# | Table | Rows | &#9660;&nbsp;Total size | Table size | Index(es) Size | TOAST Size |
|---|---|------|------------|------------|----------------|------------|
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- if le $i $.LISTLIMIT -}}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
|{{ $value.num}} | {{if eq (index $value "table") "=====TOTAL=====" }}{{ index $value "table" }}{{else}}`{{ index $value "table" }}`{{end}} | {{ index $value "rows"}} | {{ index $value "total_size"}} | {{ index $value "table_size"}} | {{ index $value "indexes_size"}} | {{ index $value "toast_size"}}|
{{/* if limit list */}}{{ end -}}
{{ end }}
{{- else -}}{{/*Master data*/}}
Nothing found
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master results*/}}
Nothing found
{{- end }}{{/*Master results*/}}
{{- else -}}{{/*Master*/}}
Nothing found
{{ end }}{{/*Master*/}}

## Conclusions ##


## Recommendations ##

