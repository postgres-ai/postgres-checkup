# {{ .checkId }} Useful Linux tools

## Observations ##
{{ if .hosts.master }}
{{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
### {{ $key }}
Utility | Availability
--------|--------------
{{- range $k, $util_name := (index (index $value) "_keys") }}
{{ $util_name }} | {{ (index (index $value) $util_name) }}
{{- end }}
{{ end }}{{ end }}{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers:  
    {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`)  
{{ $.results }}


{{ end }} {{ end }}

