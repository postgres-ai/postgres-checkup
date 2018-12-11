# Connections and current activity #

## Current values ##

### Master DB server is `{{.hosts.master}}` ###

{{ Code (index (index .results .hosts.master) "data").raw false }}

{{/* newline */}}
{{/* newline */}}

{{- if gt (len .hosts.replicas) 0 -}}
### Slave DB servers: ###
{{/* newline */}}
{{/* newline */}}
  {{- range $skey, $host := .hosts.replicas -}}
#### DB slave server: `{{ $host }}` ####
    {{- if (index $.results $host) -}}
{{/* newline */}}
{{/* newline */}}
```
{{ (index (index $.results $host) "data").raw }}
```
    {{- else -}}
```
No data
```
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/* newline */}}
{{/* newline */}}
## Conclusions ##

{{.Conclusion}}

## Recommendations ##

{{.Recommended}}
