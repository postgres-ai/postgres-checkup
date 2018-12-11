# {{ .checkId }} Cluster information #

## Observations ##

### Master (`{{.hosts.master}}`) ###
 Indicator | Value
-----------|-------
{{ range $key, $value := (index (index .results .hosts.master) "data") }}{{ $key }} | {{ Nobr (index $value "value") }}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}
 Indicator | Value
-----------|-------
{{ range $key, $value := (index (index $.results $host) "data") }}{{ $key }} | {{ Nobr (index $value "value") }}
{{ end }}
    {{ else }}
      No data
{{ end}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

