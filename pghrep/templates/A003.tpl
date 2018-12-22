# {{ .checkId }} Postgres settings #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
Setting | Value | Unit | Pretty value
--------|-------|------|--------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }}) | {{ Br $value.setting }} | {{ if $value.unit }}{{ $value.unit }} {{ end }} | {{ UnitValue $value.setting $value.unit }}
{{ end }}
{{ else }}
No data
{{ end }}

## Conclusions ##


## Recommendations ##

