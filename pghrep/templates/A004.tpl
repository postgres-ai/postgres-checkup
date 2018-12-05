Current values
===

Cluster info

Master DB server is {{.hosts.master}}
{{ range $key, $value := (index (index .results .hosts.master) "data") }}
  {{ $key }}: {{ (index $value "value") }}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
Slave DB servers:
  {{ range $skey, $host := .hosts.replicas }}
  DB slave server: {{ $host }}
    {{ if (index $.results $host) }}
      {{ range $key, $value := (index (index $.results $host) "data") }}
        {{ $key }}: {{ (index $value "value") }}
      {{ end }}
    {{ else }}
      No data
    {{ end}}
  {{ end }}
{{ end }}

Conclusions
===
{{.Conclusion}}

Recommendations
===
{{.Recommended}}
