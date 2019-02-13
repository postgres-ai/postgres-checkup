# {{ .checkId }} Autovacuum: Heap bloat #
:warning: This report is based on estimations. The errors in bloat estimates may be significant (in some cases, up to 15% and even more). Use it only as an indicator of potential issues.

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
### Master (`{{.hosts.master}}`) ###
{{ if (index (index .results .hosts.master) "data") }}
 Table | Size | Extra | &#9660;&nbsp;Estimated bloat | Est. bloat, bytes | Est. bloat ratio,% | Live | Last vacuum | Fillfactor
-------|------|-------|------------------------------|------------------|--------------------|------|-------------|------------
**Total** | **{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "Real size bytes sum" ) 2 }}** ||**{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "Bloat size bytes sum" ) 2 }}** |**{{- RawIntFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "Bloat size bytes sum" ) }}** ||||
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "heap_bloat") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "heap_bloat") $key ) -}}
{{ $key }} |
{{- ByteFormat ( index $value "Real size bytes" ) 2 }} |
{{- "~" }}{{ ByteFormat ( index $value "Extra size bytes" ) 2 }} ({{- NumFormat ( index $value "Extra_ratio" ) 2 }}%)|
{{- ByteFormat ( index $value "Bloat size bytes" ) 2 }} |
{{- RawIntFormat ( index $value "Bloat size bytes" ) }} |
{{- RawFloatFormat ( index $value "Bloat ratio") 2 }} |
{{- "~" }}{{ ByteFormat ( index $value "Live bytes" ) 2 }} |
{{- if (index $value "Last Vaccuum") }} {{ ( index $value "Last Vaccuum") }} {{ end }} |
{{- ( index $value "Fillfactor") }}
{{ end }} {{/*range*/}}
{{- else -}}
No data
{{- end -}}
{{- else -}}{{/*Master data*/}}
No data
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master*/}}
No data
{{ end }}{{/*Master*/}}

## Conclusions ##


## Recommendations ##

