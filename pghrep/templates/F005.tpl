# {{ .checkId }} Autovacuum: Index Bloat (Estimated) #
:warning: This report is based on estimations. The errors in bloat estimates may be significant (in some cases, up to 15% and even more). Use it only as an indicator of potential issues.
{{- $minRatioWarning:=40 }}

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master)}}
{{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
{{ if gt (len (index (index (index $.results $.hosts.master) "data") "index_bloat")) .ROWS_LIMIT }}The list is limited to {{.ROWS_LIMIT}} items.{{ end }}  

\# | Index (Table) | Table Size |Index Size | Extra | &#9660;&nbsp;Estimated bloat | Est. bloat, bytes | Est. bloat ratio, % | Live Data Size | Fill factor
---|---------------|------------|-----------|-------|------------------------------|-------------------|---------------------|----------------|-------------
&nbsp;|===== TOTAL ===== |
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "table_size_bytes_sum" ) 2 }} |
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "real_size_bytes_sum" ) 2 }} ||
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_size_bytes_sum" ) 2 }} |
{{- RawIntFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_size_bytes_sum" ) }}|
{{- if ge (Int (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_ratio_percent_avg" )) $minRatioWarning }}**{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_ratio_percent_avg" ) 2 }}**{{else}}{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "bloat_ratio_percent_avg" ) 2 }}{{end}}||
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "index_bloat") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "index_bloat") $key) -}}
{{- $tableIndex := Split $key "\n" -}}
{{ $value.num }} |
{{- $table := Trim (index $tableIndex 1) " ()"}}{{ (index $tableIndex 0) }} ({{ $table }}{{if $value.overrided_settings}}<sup>*</sup>{{ end }}) |
{{- ByteFormat ( index $value "table_size_bytes") 2 }} |
{{- ByteFormat ( index $value "real_size_bytes") 2 }} |
{{- if ( index $value "extra_size_bytes")}}{{- "~" }}{{ ByteFormat ( index $value "extra_size_bytes" ) 2 }} ({{- NumFormat ( index $value "extra_ratio_percent" ) 2 }}%){{end}} |
{{- if ( index $value "bloat_size_bytes")}}{{ ByteFormat ( index $value "bloat_size_bytes") 2 }}{{end}} |
{{- if ( index $value "bloat_size_bytes")}}{{ RawIntFormat ( index $value "bloat_size_bytes") }}{{end}} |
{{- if ge (Int (index $value "bloat_ratio_percent")) $minRatioWarning }} **{{- RawFloatFormat ( index $value "bloat_ratio_percent") 2 }}**{{else}}{{- RawFloatFormat ( index $value "bloat_ratio_percent") 2 }}{{end}} |
{{- "~" }}{{ ByteFormat ( index $value "live_bytes" ) 2 }} |
{{- ( index $value "fillfactor") }}
{{ end }}
{{- if gt (Int (index (index (index .results .hosts.master) "data") "overrided_settings_count")) 0 }}
<sup>*</sup> This table has specific autovacuum settings. See 'F001 Autovacuum: Current settings'
{{- end }}
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

