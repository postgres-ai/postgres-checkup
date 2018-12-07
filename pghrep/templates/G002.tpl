Current values
===

Connections and current activity

Master DB server '{{.hosts.master}}':

```
{{ (index (index .results .hosts.master) "data").raw }}
```
{{/* newline */}}
{{/* newline */}}

{{- if gt (len .hosts.replicas) 0 -}}
Slave DB servers
{{/* newline */}}
{{/* newline */}}
  {{- range $skey, $host := .hosts.replicas -}}
  DB slave server: '{{ $host }}':
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
Conclusions
===
{{.Conclusion}}

Recommendations
===
{{.Recommended}}
