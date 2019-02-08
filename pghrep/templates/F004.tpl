# {{ .checkId }} Autovacuum: Heap bloat #
:warning: This report is based on estimations. The errors in bloat estimates may be significant (in some cases, up to 15% and even more). Use it only as an indicator of potential issues.

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
{{ if (index (index .results .hosts.master) "data") }}
 Table | Size | Extra | &#9660;&nbsp;Estimated bloat | Est. bloat bytes | Est. bloat ratio,% | Live | Last vacuum | Fillfactor
-------|------|-------|------------------------------|------------------|--------------------|------|-------------|------------
**Total** | **{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "Real size bytes sum" ) 2 }}** ||**{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "Bloat size bytes sum" ) 2 }}** |**{{- NumFormat  (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "Bloat size bytes sum" ) -1 }}** | Avg: **{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "Avg bloat ratio" ) 2 }}** |||
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "heap_bloat") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "heap_bloat") $key ) -}}
{{ $key }} |
{{- ( index $value "Size") }} |
{{- ( index $value "Extra") }} |
{{- ByteFormat ( index $value "Bloat size bytes" ) 2 }} |
{{- NumFormat ( index $value "Bloat size bytes" ) -1 }} |
{{- RawFloatFormat ( index $value "Bloat ratio") 2 }} |
{{- ( index $value "Live") }} |
{{- if (index $value "Last Vaccuum") }} {{ ( index $value "Last Vaccuum") }} {{ end }} |
{{- ( index $value "Fillfactor") }}
{{ end }} {{/*range*/}}
{{- else -}}
No data
{{- end -}}
{{- else -}}
No data
{{ end }}

## Conclusions ##


## Recommendations ##

