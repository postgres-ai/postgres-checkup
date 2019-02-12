# {{ .checkId }} Autovacuum: Index bloat #
:warning: This report is based on estimations. The errors in bloat estimates may be significant (in some cases, up to 15% and even more). Use it only as an indicator of potential issues.

## Observations ##
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
### Master (`{{.hosts.master}}`) ###
 Index (Table) | &#9660;&nbsp;Size | Extra | Bloat | Live | Fill factor
---------------|------|-------|-------|------|-------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
{{- $tableIndex := Split $key "\n" -}}
{{ $table := Trim (index $tableIndex 1) " ()"}}{{ (index $tableIndex 0) }} ({{ $table }}) | {{ ( index $value "Size") }} | {{ ( index $value "Extra") }} | {{ ( index $value "Bloat") }} | {{ ( index $value "Live") }} | {{ ( index $value "fillfactor") }}
{{ end }}
{{- else -}}{{/*Master data*/}}
No data
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master*/}}
No data
{{ end }}{{/*Master*/}}

## Conclusions ##


## Recommendations ##

