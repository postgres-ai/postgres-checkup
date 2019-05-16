# {{ .checkId }} Integer (int2, int4) Out-of-range Risks in PKs #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master)}}
{{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
| Table | PK | Type | Current max value | &#9660;&nbsp;Capacity used, % |
|------|----|------|-------------------|-------------------------------|
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
|`{{ index $value "Table"}}` | `{{ index $value "PK"}}` | {{ index $value "Type"}} | {{- RawIntFormat (index $value "Current max value")}} | {{ index $value "Capacity used, %"}}|
{{ end }}
{{- else -}}{{/*Master data*/}}
No data
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master results*/}}
No data
{{- end }}{{/*Master results*/}}
{{- else -}}{{/*Master*/}}
No data
{{ end }}{{/*Master*/}}

## Conclusions ##


## Recommendations ##
