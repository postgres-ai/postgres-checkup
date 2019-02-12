# {{ .checkId }} Autovacuum: Heap bloat #
:warning: This report is based on estimations. The errors in bloat estimates may be significant (in some cases, up to 15% and even more). Use it only as an indicator of potential issues.

## Observations ##
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
### Master (`{{.hosts.master}}`) ###
{{ if (index (index .results .hosts.master) "data") }}
 Table | Size | Extra | &#9660;&nbsp;Estimated bloat | Est. bloat, bytes | Est. bloat ratio,% | Live | Last vacuum | Fillfactor
-------|------|-------|------------------------------|------------------|--------------------|------|-------------|------------
**Total** | **{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "Real size bytes sum" ) 2 }}** ||**{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "Bloat size bytes sum" ) 2 }}** |**{{- RawIntFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "Bloat size bytes sum" ) }}** |**{{- RawFloatFormat (index (index (index (index $.results $.hosts.master) "data") "heap_bloat_total") "Bloat ratio" ) 2 }}**|||
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "heap_bloat") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "heap_bloat") $key ) -}}
{{ $key }} |
{{- ByteFormat ( index $value "Real size bytes" ) 2 }} |
{{- "~" }}{{ ByteFormat ( index $value "Extra size bytes" ) 2 }} ({{- NumFormat ( index $value "Extra_ratio" ) 2 }}%)|
{{- if ( index $value "Bloat size bytes")}}{{- if ge (Int (index $value "Bloat size bytes")) 1073741824 }}**{{ ByteFormat ( index $value "Bloat size bytes") 2 }}**{{else}}{{ ByteFormat ( index $value "Bloat size bytes") 2 }}{{end}}{{end}} |
{{- if ( index $value "Bloat size bytes")}}{{- if ge (Int (index $value "Bloat size bytes")) 1073741824 }}**{{ RawIntFormat ( index $value "Bloat size bytes") }}**{{else}}{{ RawIntFormat ( index $value "Bloat size bytes") }}{{end}}{{end}} |
{{- if ge (Int (index $value "Bloat ratio")) 10 }} **{{- RawFloatFormat ( index $value "Bloat ratio") 2 }}**{{else}}{{- RawFloatFormat ( index $value "Bloat ratio") 2 }}{{end}} |
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

