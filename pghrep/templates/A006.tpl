# {{ .checkId }} Postgres setting deviations #

## Observations ##
{{ if .diffData }}
### Settings (pg_settings) that differ ###
{{ if (index .diffData "pg_settings") }}
Setting | {{.hosts.master}} {{ range $skey, $host := .hosts.replicas }}| {{ $host }} {{ end }}
--------|-------{{ range $skey, $host := .hosts.replicas }}|-------- {{ end }}
{{ range $key, $value := (index .diffData "pg_settings") }}{{ $key }} {{ range $key, $value := $value }} | {{ if (index $value "unit") }}{{ UnitValue (index $value "value") (index $value "unit") }}{{else}}{{ index $value "value" }}{{ end  }}{{ end }}
{{ end }}
{{ else }}
No `pg_settings` differences
{{end}}
{{ if (index .diffData "pg_configs") }}
Configs(pg_config) that differ
### Configs(pg_config) that differ ###
{{ range $key, $value := (index .diffData "pg_configs") }}
Config {{ $key }}: {{ range $key, $value := $value }} On {{ $key }}: {{ if (index $value "unit") }}{{ UnitValue (index $value "value") (index $value "unit") }}{{else}}{{ index $value "value" }}{{ end  }}{{ end }}
{{ end }}
{{ else }}
No `pg_configs` differences
{{end}}
{{ else }}
No data
{{ end }}

## Conclusions ##


## Recommendations ##

