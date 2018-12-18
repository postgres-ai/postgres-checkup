# {{ .checkId }} Postgres setting deviations #

## Observations ##

### Settings (pg_settings) that differ ###
{{ if .diffData }}
{{ if (index .diffData "pg_settings") }}
Setting | {{.hosts.master}} {{ range $skey, $host := .hosts.replicas }}| {{ $host }} {{ end }}
--------|-------{{ range $skey, $host := .hosts.replicas }}|-------- {{ end }}
{{ range $key, $value := (index .diffData "pg_settings") }}{{ $key }} {{ range $key, $value := $value }} |{{ index $value "value" }} {{ if (index $value "unit") }}({{ index $value "unit" }}){{ end  }}{{ end }}
{{ end }}{{end}}
{{ if (index .diffData "pg_configs") }}
Configs(pg_config) that differ
### Configs(pg_config) that differ ###
{{ range $key, $value := (index .diffData "pg_configs") }}
Config {{ $key }}: {{ range $key, $value := $value }} On {{ $key }}: `{{ index $value "value" }}` {{ if (index $value "unit") }}{{ index $value "unit" }}{{ end  }}{{ end }}
{{ end }}{{end}}
{{ else }}
No differences
{{ end }}

## Conclusions ##


## Recommendations ##

