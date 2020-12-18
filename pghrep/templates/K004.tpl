# {{ .checkId }} Top-{{.LISTLIMIT}} Query Groups by `calls`

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
{{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
Start: {{ (index (index (index .results .hosts.master) "data") "start_timestamptz") }}  
End: {{ (index (index (index .results .hosts.master) "data") "end_timestamptz") }}  
Period seconds: {{ (index (index (index .results .hosts.master) "data") "period_seconds") }}  
Period age: {{ (index (index (index .results .hosts.master) "data") "period_age") }}  

Error (calls): {{ NumFormat (index (index (index .results .hosts.master) "data") "absolute_error_calls") 2 }} ({{ NumFormat (index (index (index .results .hosts.master) "data") "relative_error_calls") 2 }}%)  
Error (total time): {{ NumFormat (index (index (index .results .hosts.master) "data") "absolute_error_total_time") 2 }} ({{ NumFormat (index (index (index .results .hosts.master) "data") "relative_error_total_time") 2 }}%)

{{ if (index (index (index .results .hosts.master) "data") "top_frequent") }}
{{ if gt (len (index (index (index .results .hosts.master) "data") "top_frequent")) .LISTLIMIT }}The list is limited to {{.LISTLIMIT}} items.{{ end }}  

| \# | Query | &#9660;&nbsp;per_sec_calls | ratio_calls | per_call_total_time | ratio_total_time | per_call_rows | ratio_rows |
|----|-------|----------------------------|-------------|---------------------|------------------|---------------|------------|
{{ range $i, $value := (index (index (index .results .hosts.master) "data") "top_frequent") }}
{{- if lt $i $.LISTLIMIT -}}
{{- $num:= Add $i 1 -}}
| {{- $num }} | 
{{- EscapeQuery (WordWrap (LimitStr $value.query 1000) 30) }}<br/>[Full query]({{ $value.link }}) | 
{{- NumFormat $value.per_sec_calls 2 }}/sec | 
{{- NumFormat $value.ratio_calls 2 }}% |
{{- MsFormat $value.per_call_total_time }}/call |
{{- NumFormat $value.ratio_total_time 2 }}% |
{{- NumFormat $value.per_call_rows 2 }}/call | 
{{- NumFormat $value.ratio_rows 2 }}% |
{{/* if limit list */}}{{ end -}}
{{ end }}{{/* range */}}

{{ else }}
Nothing found
{{ end }}{{/* top_frequernt exists*/}}

{{- end }}{{/*Master data*/}}
{{- end }}{{/*Master data*/}}
{{ end }}{{/*Master*/}}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
{{ range $key, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
{{ if (index $.results $host) }}
Start: {{ (index (index (index $.results $host) "data") "start_timestamptz") }}  
End: {{ (index (index (index $.results $host) "data") "end_timestamptz") }}  
Period seconds: {{ (index (index (index $.results $host) "data") "period_seconds") }}  
Period age: {{ (index (index (index $.results $host) "data") "period_age") }}  

{{ if (index (index (index $.results $host) "data") "top_frequent") }}
{{ if gt (len (index (index (index $.results $host) "data") "top_frequent")) 50 }}Top 50 rows{{ end }}  

| \# | Query | &#9660;&nbsp;per_sec_calls | ratio_calls | per_call_total_time | ratio_total_time | per_call_rows | ratio_rows |
|----|-------|----------------------------|-------------|---------------------|------------------|---------------|------------|
{{ range $i, $value := (index (index (index $.results $host) "data") "top_frequent") }}
{{- if lt $i $.LISTLIMIT -}}
{{- $num:= Add $i 1 -}}
| {{- $num }} | 
{{- EscapeQuery (WordWrap (LimitStr $value.query 1000) 30) }}<br/>[Full query]({{ $value.link }}) | 
{{- NumFormat $value.per_sec_calls 2 }}/sec | 
{{- NumFormat $value.ratio_calls 2 }}% |
{{- MsFormat $value.per_call_total_time }}/call |
{{- NumFormat $value.ratio_total_time 2 }}% |
{{- NumFormat $value.per_call_rows 2 }}/call | 
{{- NumFormat $value.ratio_rows 2 }}% |
{{/* if limit list */}}{{ end -}}
{{ end }}{{/* range */}}
  
{{ else }}
Nothing found
{{ end }}{{/* top_frequernt exists*/}}

{{- else -}}{{/* if host data */}}
Nothing found
{{- end -}}{{/* if host data */}}
{{- end -}}{{/* hosts range */}}
{{- end -}}{{/* if replicas */}}

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
 {{end}}
{{ end }}
