# {{ .checkId }} Autovacuum: Heap bloat #
:warning: This report is based on estimations. The errors in bloat estimates may be significant (in some cases, up to 15% and even more). Use it only as an indicator of potential issues.

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
{{ if (index (index .results .hosts.master) "data") }}
 Table | &#9660;&nbsp;Size | Extra | Bloat | Live | Last vacuum | Fillfactor
-------|------|-------|-------|------|-------------|-------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
{{ $key }} | {{ ( index $value "Size") }} | {{ ( index $value "Extra") }} | {{ ( index $value "Bloat") }} | {{ ( index $value "Live") }} | {{ if (index $value "Last Vaccuum") }} {{ ( index $value "Last Vaccuum") }} {{ end }} | {{ ( index $value "Fillfactor") }}
{{ end }}
{{- else -}}
No data
{{- end -}}
{{- else -}}
No data
{{ end }}

## Conclusions ##


## Recommendations ##

