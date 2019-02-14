# {{ .checkId }} Postgres settings #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
{{ if .hosts.master }}
{{ if and (index .results .hosts.master) (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###  
&#9660;&nbsp;Category | Setting | Value | Unit | Pretty value
---------|---------|-------|------|--------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
{{ $value.category }}|[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }}) | {{ Br $value.setting }} | {{ if $value.unit }}{{ $value.unit }} {{ end }} | {{ UnitValue $value.setting $value.unit }}
{{ end }}
{{- else -}}{{/*Master data*/}}
No data  
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master*/}}
No data  
{{ end }}{{/*Master*/}}
  