Current values
===

Extensions

Master DB server is {{.hosts.master}}
{{ range $key, $value := (index (index .results .hosts.master) "data") }}
  Extension: {{ $key }}, Installed version: {{ $value.installed_version}}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
Slave DB servers:
  {{ range $skey, $host := .hosts.replicas }}
  DB slave server: {{ $host }}
    {{ if (index $.results $host) }}
      {{ range $key, $value := (index (index $.results $host) "data") }}
        Extension: {{ $key }}, Installed version: {{ $value.installed_version}}
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
