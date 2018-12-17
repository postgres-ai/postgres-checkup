# {{ .checkId }} Cluster information #

## Observations ##

### Master (`{{.hosts.master}}`) ###
 Indicator | Value
-----------|-------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
{{ $key }} | {{ Nobr (index $value "value") }}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}
 Indicator | Value
-----------|-------
{{ range $i, $key := (index (index (index $.results $host) "data") "_keys") }}
{{- $value := (index (index (index $.results $host) "data") $key) -}}
{{ $key }} | {{ Nobr (index $value "value") }}
{{ end }}
{{- else -}}
No data
{{- end -}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

