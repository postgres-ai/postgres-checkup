# {{ .checkId }} Autovacuum info #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
Setting name | Value | Unit | Pretty value
-------------|-------|------|--------------
{{ range $i, $key := (index (index (index (index (index .results .hosts.master) "data") "settings") "global_settings") "_keys") }}
{{- $value := (index (index (index (index (index $.results $.hosts.master) "data") "settings") "global_settings" ) $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting }}|{{ $value.unit }} | {{ UnitValue $value.setting $value.unit}}
{{ end }}
{{ if (index (index (index (index .results .hosts.master) "data") "settings") "table_settings") }}
#### Tables settings override ####
Namespace | Relation | Options
----------|----------|------
{{ range $i, $key := (index (index (index (index (index .results .hosts.master) "data") "settings") "table_settings") "_keys") }}
{{- $value := (index (index (index (index (index $.results $.hosts.master) "data") "settings") "table_settings") $key) -}}
{{ $value.namespace }} | {{ $value.relname }}|{{ $value.reloptions }}
{{ end }}
{{- end -}}

{{- if (index (index .results .hosts.master) "data").iotop -}}
#### iotop information ####
Command: `{{ (index (index (index .results .hosts.master) "data") "iotop").cmd }}`
Result:

```
{{- (index (index (index .results .hosts.master) "data") "iotop").data -}}
```

{{- end -}}
{{- else -}}
No data
{{ end }}

## Conclusions ##


## Recommendations ##

