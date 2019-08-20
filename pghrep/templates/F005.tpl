# {{ .checkId }} Autovacuum: Index Bloat (Estimated) #
:warning: This report is based on estimations. The errors in bloat estimates may be significant (in some cases, up to 15% and even more). Use it only as an indicator of potential issues.
{{- $minRatioWarning:=40 }}

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{ if gt (Int (index (index (index .results .reorderedHosts.master) "data") "min_table_size_bytes")) 0 }}NOTICE: only tables larger than {{ ByteFormat (index (index (index .results .reorderedHosts.master) "data") "min_table_size_bytes") 0 }} are analyzed.  {{end}}
{{ if .hosts.master }}
{{ if (index .results .hosts.master)}}
{{ if (index (index .results .hosts.master) "data") }}
{{ if (index (index (index .results .hosts.master) "data") "index_bloat") }}
### Master (`{{.hosts.master}}`) ###
{{ if ge (len (index (index (index $.results $.hosts.master) "data") "index_bloat")) .LISTLIMIT }}The list is limited to {{.LISTLIMIT}} items. Total: {{ Sub (len (index (index (index $.results $.hosts.master) "data") "index_bloat")) 1 }}. {{ end }}  

| \# | Index (Table) | Index Size | Table Size | &#9660;&nbsp;Estimated bloat | Est. bloat, bytes | Est. bloat factor |Est. bloat level, % | Live Data Size | Fillfactor  |
|----|---------------|------------|------------|------------------------------|-------------------|-------------------|--------------------|----------------|-------------|
|&nbsp;|===== TOTAL ===== |
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "real_size_bytes_sum" ) 2 }} |
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "table_size_bytes_sum" ) 2 }} |
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_size_bytes_sum" ) 2 }} |
{{- RawIntFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_size_bytes_sum" ) }}|
{{- if (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_ratio_factor_avg") }}{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_ratio_factor_avg" ) 2 }}{{ end }} |
{{- if ge (Int (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_ratio_percent_avg" )) $minRatioWarning }}**{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_ratio_percent_avg" ) 2 }}**{{else}}{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_ratio_percent_avg" ) 2 }}{{end}}|
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "live_data_size_bytes_sum" ) 2 }} |||
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "index_bloat") "_keys") }}
{{- if lt $i $.LISTLIMIT -}}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "index_bloat") $key) -}}
|{{ $value.num }} | `{{ $value.index_name }}` (`{{ $value.table_name }}`{{if $value.overrided_settings}}\*{{ end }}) |
{{- ByteFormat ( index $value "real_size_bytes") 2 }} |
{{- ByteFormat ( index $value "table_size_bytes") 2 }} |
{{- if ( index $value "bloat_size_bytes")}}{{ ByteFormat ( index $value "bloat_size_bytes") 2 }}{{end}} |
{{- if ( index $value "bloat_size_bytes")}}{{ RawIntFormat ( index $value "bloat_size_bytes") }}{{end}} |
{{- if ( index $value "bloat_ratio_factor")}}{{ RawFloatFormat ( index $value "bloat_ratio_factor") 2 }}{{end}} |
{{- if ge (Int (index $value "bloat_ratio_percent")) $minRatioWarning }} **{{- RawFloatFormat ( index $value "bloat_ratio_percent") 2 }}**{{else}}{{- RawFloatFormat ( index $value "bloat_ratio_percent") 2 }}{{end}} |
{{- "~" }}{{ ByteFormat ( index $value "live_data_size_bytes" ) 2 }} |
{{- ( index $value "fillfactor") }} |
{{/* if limit list */}}{{ end -}}
{{ end }}
{{- if gt (Int (index (index (index .results .hosts.master) "data") "overrided_settings_count")) 0 }}
\* This table has specific autovacuum settings. See 'F001 Autovacuum: Current settings'
{{- end }}
{{- else -}}{{/*Index bloat*/}}
Nothing found
{{- end }}{{/*Index bloat*/}}
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
