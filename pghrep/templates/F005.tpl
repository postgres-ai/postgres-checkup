# {{ .checkId }} Autovacuum: Index bloat #
:warning: This report is based on estimations. The errors in bloat estimates may be significant (in some cases, up to 15% and even more). Use it only as an indicator of potential issues.

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
 Index (Table) | &#9660;&nbsp;Size | Extra | Bloat | Bloat, bytes | Bloat ratio,% | Live | Fill factor
---------------|-------------------|-------|-------|-------------|-------------|------|-------------
**Total** |**{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "Real size bytes sum" ) 2 }}** ||**{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "Bloat size bytes sum" ) 2 }}** |**{{- RawIntFormat (index (index (index (index $.results $.hosts.master) "data") "index_bloat_total") "Bloat size bytes sum" ) }}**|||
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "index_bloat") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "index_bloat") $key) -}}
{{- $tableIndex := Split $key "\n" -}}
{{ $table := Trim (index $tableIndex 1) " ()"}}{{ (index $tableIndex 0) }} ({{ $table }}) |
{{- ByteFormat ( index $value "Real size bytes") 2 }} |
{{- ( index $value "Extra") }} |
{{- if ( index $value "Bloat size bytes")}}{{ ByteFormat ( index $value "Bloat size bytes") 2 }}{{end}} |
{{- if ( index $value "Bloat size bytes")}}{{- RawIntFormat ( index $value "Bloat size bytes") }}{{end}} |
{{- if ( index $value "Bloat ratio")}}{{- RawFloatFormat ( index $value "Bloat ratio") 2 }}{{end}} |
{{- ( index $value "Live") }} |
{{- ( index $value "fillfactor") }}
{{ end }}
{{- else -}}
No data
{{ end }}

## Conclusions ##


## Recommendations ##

