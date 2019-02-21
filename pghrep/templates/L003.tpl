# {{ .checkId }} Integer (int2, int4) out-of-range risks in PKs #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{ if .hosts.master }}
{{ if and (index .results .hosts.master) (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
Table | Column | Type | &#9660;&nbsp;Reached value
------|--------|------|---------------------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
{{ index $value "Table"}} | {{ index $value "Column"}} | {{ index $value "Type"}} | {{ index $value "Reached value"}}
{{ end }}
{{- else -}}{{/*Master data*/}}
No data
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master*/}}
No data
{{ end }}{{/*Master*/}}

## Conclusions ##


## Recommendations ##

