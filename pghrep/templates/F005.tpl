# {{ .checkId }} Autovacuum: Index bloat #
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

\# | Index (Table) | &#9660;&nbsp;Size | Extra | Estimated bloat | Est. bloat, bytes | Est. bloat ratio, % | Live | Fill factor
---|------------|-------------------|-------|-------|-------------|-------------|------|-------------
&nbsp;|===== TOTAL ===== |
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "Real size bytes sum" ) 2 }} ||
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "Bloat size bytes sum" ) 2 }} |
{{- RawIntFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "Bloat size bytes sum" ) }}|
{{- if ge (Int (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "Bloat ratio" )) $minRatioWarning }}**{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "Bloat ratio" ) 2 }}**{{else}}{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "Bloat ratio" ) 2 }}{{end}}||
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "index_bloat") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "index_bloat") $key) -}}
{{- $tableIndex := Split $key "\n" -}}
{{ $value.num }} |
{{- $table := Trim (index $tableIndex 1) " ()"}}{{ (index $tableIndex 0) }} ({{ $table }}{{if $value.overrided_settings}}<sup>*</sup>{{ end }}) |
{{- ByteFormat ( index $value "Real size bytes") 2 }} |
{{- if ( index $value "Extra size bytes")}}{{- "~" }}{{ ByteFormat ( index $value "Extra size bytes" ) 2 }} ({{- NumFormat ( index $value "Extra_ratio" ) 2 }}%){{end}} |
{{- if ( index $value "Bloat size bytes")}}{{ ByteFormat ( index $value "Bloat size bytes") 2 }}{{end}} |
{{- if ( index $value "Bloat size bytes")}}{{ RawIntFormat ( index $value "Bloat size bytes") }}{{end}} |
{{- if ge (Int (index $value "Bloat ratio")) $minRatioWarning }} **{{- RawFloatFormat ( index $value "Bloat ratio") 2 }}**{{else}}{{- RawFloatFormat ( index $value "Bloat ratio") 2 }}{{end}} |
{{- "~" }}{{ ByteFormat ( index $value "Live bytes" ) 2 }} |
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

