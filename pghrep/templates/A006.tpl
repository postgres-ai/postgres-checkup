# {{ .checkId }} Postgres Setting Deviations #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
{{ if .diffData }}
### Settings (pg_settings) that Differ ###
{{ if (index .diffData "pg_settings") }}
&#9660;&nbsp;Setting | {{.hosts.master}} {{ range $skey, $host := .hosts.replicas }}| {{ $host }} {{ end }}
--------|-------{{ range $skey, $host := .hosts.replicas }}|-------- {{ end }}
{{ range $key, $value := (index .diffData "pg_settings") }}{{ $key }} {{ range $key, $value := $value }} |{{ if and ( ne (index $value "value") "-1") ( ne (index $value "value") "0") (index $value "unit") }}{{ if (UnitValue (index $value "value") (index $value "unit")) }}{{ UnitValue (index $value "value") (index $value "unit") }}{{else}}{{(index $value "value")}} {{(index $value "unit") }}{{end}}{{else}}{{ index $value "value" }}{{ end }}{{ end }}
{{ end }}
{{ else }}
No differences in `pg_settings` are found.
{{end}}
### Configs(pg_config) that differ ###
{{ if (index .diffData "pg_configs") }}
&#9660;&nbsp;Config | {{.hosts.master}} {{ range $skey, $host := .hosts.replicas }}| {{ $host }} {{ end }}
--------|-------{{ range $skey, $host := .hosts.replicas }}|-------- {{ end }}
{{ range $key, $value := (index .diffData "pg_configs") }}{{ $key }} {{ range $key, $value := $value }} | {{ if and ( ne (index $value "value") "-1") ( ne (index $value "value") "0") (index $value "unit") }}{{ if (UnitValue (index $value "value") (index $value "unit")) }}{{ UnitValue (index $value "value") (index $value "unit") }}{{else}}{{(index $value "value")}} {{(index $value "unit") }}{{end}}{{else}}{{ index $value "value" }}{{ end }}{{ end }}
{{ end }}
{{ else }}
No differences in `pg_config` are found.
{{end}}
{{ else }}
    No data
{{ end }}

## Conclusions ##


## Recommendations ##

