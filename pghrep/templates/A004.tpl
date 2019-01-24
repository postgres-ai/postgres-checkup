# {{ .checkId }} Cluster information #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
{{ if (index (index (index .results .hosts.master) "data") "general_info") }}
 Indicator | Value
-----------|-------
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "general_info") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "general_info") $key) -}}
{{ $key }} | {{ Nobr (index $value "value") }}
{{ end }}
{{- end -}}
{{ if (index (index (index .results .hosts.master) "data") "database_sizes") }}
#### Databases sizes ####
Database | &#9660;&nbsp;Size
---------|------
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "database_sizes") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "database_sizes") $key) -}}
{{ $key }} | {{ ByteFormat $value 2 }}
{{ end }}
{{- end -}}
{{- end -}}
{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}
 Indicator | Value
-----------|-------
{{ range $i, $key := (index (index (index (index $.results $host) "data") "general_info") "_keys") }}
{{- $value := (index (index (index (index $.results $host) "data") "general_info") $key) -}}
{{ $key }} | {{ Nobr (index $value "value") }}
{{ end }}
{{- else -}}
No data
{{- end -}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

