# {{ .checkId }} Postgres setting deviations #

## Observations ##

### Settings (pg_settings) that differ ###
{{ if (index .diffData "pg_settings") }}
Setting | {{.hosts.master}}
{{- range $skey, $host := .hosts.replicas }}| {{ $host }} {{ end }}
--------|-------{{ range $skey, $host := .hosts.replicas }}|-------- {{ end }}
    {{- range $key, $value := (index .diffData "pg_settings") }}
{{ $key }} | {{ $hostValue := (index $value "master") }}{{ index $hostValue "value" }}{{ if (index $hostValue "unit") }}({{ index $hostValue "unit" }}){{ end }}
{{- range $h, $host := $.hosts.replicas }}| {{ $hostValue := (index $value $host)}}{{ index $hostValue "value" }}{{ if (index $hostValue "unit") }}({{ index $hostValue "unit" }}){{ end }}{{ end -}}
    {{ end }}
{{ end -}}

{{ if (index .diffData "pg_configs") }}
### Configs(pg_config) that differ ###
Setting | {{.hosts.master}}
{{- range $skey, $host := .hosts.replicas }}| {{ $host }} {{ end }}
--------|-------{{ range $skey, $host := .hosts.replicas }}|-------- {{ end }}
    {{- range $key, $value := (index .diffData "pg_configs") }}
{{ $key }} | {{ $hostValue := (index $value "master") }}{{ index $hostValue "value" }}{{ if (index $hostValue "unit") }}({{ index $hostValue "unit" }}){{ end }}
{{- range $h, $host := $.hosts.replicas }}| {{ $hostValue := (index $value $host)}}{{ index $hostValue "value" }}{{ if (index $hostValue "unit") }}({{ index $hostValue "unit" }}){{ end }}{{ end -}}
    {{ end }}
{{ end }}

## Conclusions ##


## Recommendations ##

