# {{ .checkId }} Integer (int2, int4) Out-of-range Risks in PKs #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{- if .hosts.master }}
{{- if (index .results .hosts.master)}}
{{- if (index (index .results .hosts.master) "data") }}
{{ if gt (Int (index (index (index .results .reorderedHosts.master) "data") "min_table_size_bytes")) 0 }}NOTICE: only tables larger than {{ ByteFormat (index (index (index .results .reorderedHosts.master) "data") "min_table_size_bytes") 0 }} are analyzed.  
  {{end}}
{{- if (index (index (index .results .hosts.master) "data") "tables") }}
### Master (`{{.hosts.master}}`) ###
{{ if ge (len (index (index .results .hosts.master) "data")) .LISTLIMIT }}The list is limited to {{.LISTLIMIT}} items. Total: {{ Sub (len (index (index .results .hosts.master) "data")) 1 }}.{{ end }}  

| Table | PK | Type | Current max value | &#9660;&nbsp;Capacity used, % |
|------|----|------|-------------------|-------------------------------|
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "tables") "_keys") }}
{{- if lt $i $.LISTLIMIT -}}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "tables") $key) -}}
|`{{ index $value "table"}}` | `{{ index $value "pk"}}` | {{ index $value "type"}} | {{- RawIntFormat (index $value "current_max_value")}} | {{ index $value "capacity_used_percent"}}|
{{ end }}
{{- else -}}{{/*Tables data*/}}
Nothing found
{{- end }}{{/*Tables data*/}}
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

{{- if .processed }}
 {{- if .conclusions }}
  {{ range $conclusion := .conclusions -}}
   - {{ $conclusion.Message }}
  {{ end }}
 {{else}}
 {{end}}
{{ end }}

## Recommendations ##

{{- if .processed }}
 {{- if .recommendations }}
  {{ range $recommendation := .recommendations -}}
   - {{ $recommendation.Message }}
  {{ end }}
 {{else}}
  All good, no recommendations here.
 {{end}}
{{ end }}
