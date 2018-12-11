# {{ .checkId }} Extensions #

## Observations ##

### Master (`{{.hosts.master}}`) ###
Extension name | Installed version
---------------|-------------------
{{ range $key, $value := (index (index .results .hosts.master) "data") }}{{ $key }} | {{ $value.installed_version}}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}
Extension name | Installed version
---------------|-------------------
{{ range $key, $value := (index (index $.results $host) "data") }}{{ $key }} | {{ $value.installed_version}}
{{ end }}
{{ else }}No data{{ end}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

