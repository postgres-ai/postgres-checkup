# {{ .checkId }} Postgres settings #

## Observations ##

### Master (`{{.hosts.master}}`) ###
Setting | Value | Unit
--------|-------|------
{{ range $key, $value := (index (index .results .hosts.master) "data") }}{{ $key }} | {{ $value.setting}} | {{ if $value.unit }}{{ $value.unit }} {{ end }}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}
Setting | Value | Unit
--------|-------|------
{{ range $key, $value := (index (index $.results $host) "data") }}{{ $key }} | {{ $value.setting}} | {{ if $value.unit }}{{ $value.unit }} {{ end }}
{{ end }}
    {{ else }}
No data
{{ end}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

