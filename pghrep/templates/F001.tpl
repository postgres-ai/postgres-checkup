# {{ .checkId }} Autovacuum: Current settings #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
&#9660;&nbsp;Setting name | Value | Unit | Pretty value
-------------|-------|------|--------------
{{ range $i, $key := (index (index (index (index (index .results .hosts.master) "data") "settings") "global_settings") "_keys") -}}
{{- if ne $key "hot_standby_feedback" -}}
{{- $value := (index (index (index (index (index $.results $.hosts.master) "data") "settings") "global_settings" ) $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting }}|{{ $value.unit }} | {{ UnitValue $value.setting $value.unit}}
{{ end -}}
{{ end }}{{/* range */}}

{{ if (index (index (index (index .results .hosts.master) "data") "settings") "table_settings") }}
#### Tables settings override ####
&#9660;&nbsp;Namespace | Relation | Options
----------|----------|------
{{ range $i, $key := (index (index (index (index (index .results .hosts.master) "data") "settings") "table_settings") "_keys") }}
{{- $value := (index (index (index (index (index $.results $.hosts.master) "data") "settings") "table_settings") $key) -}}
{{ $value.namespace }} | {{ $value.relname }}|{{ $value.reloptions }}
{{- end -}}{{/* range */}}
{{- end -}}{{/* if table_settings */}}
{{- else -}}
No data
{{- end -}}{{/* master */}}
{{- if gt (len .hosts.replicas) 0 -}}

### Replicas settings ###
Setting {{ range $skey, $host := .hosts.replicas }}| {{ $host }} {{ end }}
--------{{ range $skey, $host := .hosts.replicas }}|-------- {{ end }}
[hot_standby_feedback](https://postgresqlco.nf/en/doc/param/hot_standby_feedback)
{{- range $skey, $host := .hosts.replicas -}}| {{- $value := (index (index (index (index (index $.results $host) "data") "settings") "global_settings") "hot_standby_feedback") -}}{{- $value.setting -}}
{{- end -}}{{/* range replicas */}}
{{ end }}{{/* if replicas */}}

## Conclusions ##


## Recommendations ##

