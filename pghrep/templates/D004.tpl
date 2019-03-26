# {{ .checkId }} pg_stat_statements and pg_stat_kcache Settings #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
{{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
{{ if (index (index (index .results .hosts.master) "data") "pg_stat_statements") }}
#### `pg_stat_statements` extension settings ####
Setting | Value | Unit | Type | Min value | Max value
--------|-------|------|------|-----------|-----------
{{ range $i, $setting_name := (index (index (index (index .results .hosts.master) "data") "pg_stat_statements") "_keys") }}
{{- $setting_data := (index (index (index (index $.results $.hosts.master) "data") "pg_stat_statements") $setting_name) -}}
[{{ $setting_name }}](https://postgresqlco.nf/en/doc/param/{{ $setting_name }})|{{ $setting_data.setting }}|{{ if $setting_data.unit }}{{ $setting_data.unit }} {{ end }}|{{ $setting_data.vartype }}|{{ if $setting_data.min_val }}{{ $setting_data.min_val }} {{ end }}|{{ if $setting_data.max_val }}{{ $setting_data.max_val }} {{ end }}
{{ end }}
{{- end -}}

{{ if (index (index (index .results .hosts.master) "data") "kcache") }}
#### `kcache` extension settings ####
Setting | Value | Unit | Type | Min value | Max value
--------|-------|------|------|-----------|-----------
{{ range $i, $setting_name := (index (index (index (index .results .hosts.master) "data") "kcache") "_keys") }}
{{- $setting_data := (index (index (index (index $.results $.hosts.master) "data") "kcache") $setting_name) -}}
[{{ $setting_name }}](https://postgresqlco.nf/en/doc/param/{{ $setting_name }})|{{ $setting_data.setting }}|{{ if $setting_data.unit }}{{ $setting_data.unit }} {{ end }}|{{ $setting_data.vartype }}|{{ if $setting_data.min_val }}{{ $setting_data.min_val }} {{ end }}|{{ if $setting_data.max_val }}{{ $setting_data.max_val }} {{ end }}
{{ end }}
{{- end -}}
{{- else -}}{{/*Master data*/}}
No data
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master results*/}}
No data
{{- end }}{{/*Master results*/}}
{{- else -}}{{/*Master*/}}
No data
{{ end }}{{/*Master*/}}

## Conclusions ##


## Recommendations ##

