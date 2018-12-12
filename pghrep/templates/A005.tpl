# {{ .checkId }} Extensions #

## Observations ##

### Master (`{{.hosts.master}}`) ###
Extension name | Installed version | Default version | Is old 
---------------|-------------------|-----------------|--------
{{ range $key, $value := (index (index .results .hosts.master) "data") }} {{ $key }} | {{ $value.installed_version }} | {{ $value.default_version }} | {{ $value.is_old }}
{{ end }}

{{ if gt (len .hosts.replicas) 9999 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}
Extension name | Installed version | Default version | Is old 
---------------|-------------------|-----------------|--------
{{ range $key, $value := (index (index $.results $host) "data") }} {{ $key }} | {{ $value.installed_version }} | {{ $value.default_version }} | {{ $value.is_old }}
{{ end }}
{{ else }}No data{{ end}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

