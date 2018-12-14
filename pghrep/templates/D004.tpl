# {{ .checkId }} pg_stat_statements and kcache settings #

## Observations ##

### Master (`{{.hosts.master}}`) ###

#### `pg_stat_statements` extension settings ####
Setting | Value | Unit | Type | Min value | Max value
--------|-------|------|------|-----------|-----------
{{ range $i, $setting_name := (index (index (index .results .hosts.master) "data") "pg_stat_statements_keys") }}
{{- $setting_data := (index (index (index (index $.results $.hosts.master) "data") "pg_stat_statements") $setting_name) -}}
[{{ $setting_name }}](https://postgresqlco.nf/en/doc/param/{{ $setting_name }})|{{ $setting_data.setting }}|{{ if $setting_data.unit }}{{ $setting_data.unit }} {{ end }}|{{ $setting_data.vartype }}|{{ if $setting_data.min_val }}{{ $setting_data.min_val }} {{ end }}|{{ if $setting_data.max_val }}{{ $setting_data.max_val }} {{ end }}
{{ end }}

#### `kcache` extension settings ####
Setting | Value | Unit | Type | Min value | Max value
--------|-------|------|------|-----------|-----------
{{ range $i, $setting_name := (index (index (index .results .hosts.master) "data") "kcache_keys") }}
{{- $setting_data := (index (index (index (index $.results $.hosts.master) "data") "kcache") $setting_name) -}}
[{{ $setting_name }}](https://postgresqlco.nf/en/doc/param/{{ $setting_name }})|{{ $setting_data.setting }}|{{ if $setting_data.unit }}{{ $setting_data.unit }} {{ end }}|{{ $setting_data.vartype }}|{{ if $setting_data.min_val }}{{ $setting_data.min_val }} {{ end }}|{{ if $setting_data.max_val }}{{ $setting_data.max_val }} {{ end }}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
    {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
        {{ if (index $.results $host) }}  
#### `pg_stat_statements` settings ####
Setting | Value | Unit | Type | Min value | Max value
--------|-------|------|------|-----------|-----------
{{ range $i, $setting_name := (index (index (index $.results $host) "data") "pg_stat_statements_keys") }}
{{- $setting_data := (index (index (index (index $.results $host) "data") "pg_stat_statements") $setting_name) -}}
[{{ $setting_name }}](https://postgresqlco.nf/en/doc/param/{{ $setting_name }})|{{ $setting_data.setting }}|{{ if $setting_data.unit }}{{ $setting_data.unit }} {{ end }}|{{ $setting_data.vartype }}|{{ if $setting_data.min_val }}{{ $setting_data.min_val }} {{ end }}|{{ if $setting_data.max_val }}{{ $setting_data.max_val }} {{ end }}
{{ end }}

#### `kcache` settings ####
Setting | Value | Unit | Type | Min value | Max value
--------|-------|------|------|-----------|-----------
{{ range $i, $setting_name := (index (index (index $.results $host) "data") "kcache_keys") }}
{{- $setting_data := (index (index (index (index $.results $host) "data") "kcache") $setting_name) -}}
[{{ $setting_name }}](https://postgresqlco.nf/en/doc/param/{{ $setting_name }})|{{ $setting_data.setting }}|{{ if $setting_data.unit }}{{ $setting_data.unit }} {{ end }}|{{ $setting_data.vartype }}|{{ if $setting_data.min_val }}{{ $setting_data.min_val }} {{ end }}|{{ if $setting_data.max_val }}{{ $setting_data.max_val }} {{ end }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}

## Conclusions ##


## Recommendations ##

