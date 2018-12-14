# {{ .checkId }} Postgres settings #

## Observations ##

### Master (`{{.hosts.master}}`) ###
Setting | Value | Unit
--------|-------|------
{{ range $i, $key := (index (index .results .hosts.master) "data_keys") }}[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }}){{ $value := (index (index (index $.results $.hosts.master) "data") $key) }} | {{ $value.setting}} | {{ if $value.unit }}{{ $value.unit }} {{ end }}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}
Setting | Value | Unit
--------|-------|------
{{ range $i, $key := (index (index $.results $host) "data_keys") }}[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }}){{ $value := (index (index (index $.results $host) "data") $key) }} | {{ $value.setting}} | {{ if $value.unit }}{{ $value.unit }} {{ end }}
{{ end }}
    {{ else }}
No data
{{ end}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

