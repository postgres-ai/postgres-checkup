# {{ .checkId }} Autovacuum: Heap Bloat (Estimated) #
:warning: This report is based on estimations. The errors in bloat estimates may be significant (in some cases, up to 15% and even more). Use it only as an indicator of potential issues.
{{- $minRatioWarning:=40 }}

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
### Master (`{{.hosts.master}}`) ###
{{ if (index (index .results .hosts.master) "data") }}
{{ if (index (index (index .results .hosts.master) "data") "heap_bloat") }}
{{ if gt (len (index (index (index $.results $.hosts.master) "data") "heap_bloat")) .LISTLIMIT }}The list is limited to {{.LISTLIMIT}} items.{{ end }}  

| \# | Table | Real Size | &#9660;&nbsp;Estimated bloat | Est. bloat, bytes | Est. bloat factor | Est. bloat level, % | Live Data Size | Last vacuum | Fillfactor |
|----|-------|------|------------------------------|-------------------|-------------------|---------------------|----------------|-------------|------------|
{{ if (index (index (index .results .hosts.master) "data") "heap_bloat_total") }}|&nbsp;|===== TOTAL ===== |
{{- if (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "real_size_bytes_sum") }}{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "real_size_bytes_sum") 2 }}{{ end }} |
{{- if (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "bloat_size_bytes_sum") }}{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "bloat_size_bytes_sum") 2 }}{{ end }} |
{{- if (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "bloat_size_bytes_sum") }}{{- RawIntFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "bloat_size_bytes_sum" ) }}{{ end }} |
{{- if (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "bloat_ratio_avg") }}{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "bloat_ratio_avg" ) 2 }}{{ end }} |
{{- if (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "bloat_ratio_percent_avg") }}{{- if ge (Int (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "bloat_ratio_percent_avg" )) $minRatioWarning }}**{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "bloat_ratio_percent_avg" ) 2 }}**{{else}}{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "bloat_ratio_percent_avg") 2 }}{{ end }} |
{{- if (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "live_data_size_bytes_sum") }} ~{{ ByteFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "live_data_size_bytes_sum") 2 }}{{ end }} |||{{ end }}|{{ end }}
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "heap_bloat") "_keys") }}
{{- if lt $i $.LISTLIMIT -}}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "heap_bloat") $key ) -}}
|{{ $value.num }} |`{{- $key }}`{{if $value.overrided_settings}}\*{{ end }} |
{{- ByteFormat ( index $value "real_size_bytes" ) 2 }} |
{{- if ( index $value "bloat_size_bytes")}}{{ ByteFormat ( index $value "bloat_size_bytes") 2 }}{{end}} |
{{- if ( index $value "bloat_size_bytes")}}{{ RawIntFormat ( index $value "bloat_size_bytes") }}{{end}} |
{{- if ( index $value "bloat_ratio")}}{{ RawFloatFormat ( index $value "bloat_ratio") 2 }}{{end}} |
{{- if ge (Int (index $value "bloat_ratio_percent")) $minRatioWarning }} **{{- RawFloatFormat ( index $value "bloat_ratio_percent") 2 }}**{{else}}{{- RawFloatFormat ( index $value "bloat_ratio_percent") 2 }}{{end}} |
{{- "~" }}{{ ByteFormat ( index $value "live_data_size_bytes" ) 2 }} |
{{- if (index $value "last_vaccuum") }} {{ ( index $value "last_vaccuum") }} {{ end }} |
{{- ( index $value "fillfactor") }} |
{{/* if limit list */}}{{ end -}}
{{ end }}{{/*range*/}}

{{- if gt (Int (index (index (index .results .hosts.master) "data") "overrided_settings_count")) 0 }}
\* This table has specific autovacuum settings. See 'F001 Autovacuum: Current settings'
{{- end }}
{{- else }}{{/* if heap_bloat */}}
Nothing found
{{- end -}}{{/* if heap_bloat */}}
{{- else -}}
Nothing found
{{- end -}}
{{- else -}}{{/*Master data*/}}
Nothing found
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master*/}}
Nothing found
{{ end }}{{/*Master*/}}

## Conclusions ##

{{- if .conclusions }}
{{ range $conclusion := .conclusions -}}
{{ $conclusion }}  
{{ end }}
{{ end }}

## Recommendations ##

{{- if .recommendations }}
{{ range $recommendation := .recommendations -}}
{{ $recommendation }}  
{{ end }}
{{ end }}

