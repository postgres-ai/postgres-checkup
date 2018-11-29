Current values
===

Master DB server is {{.hosts.master}}
  Version: {{(index .results .hosts.master).version}}
  Server version number: {{(index .results .hosts.master).server_version_num}}
  Server major version: {{(index .results .hosts.master).server_major_ver}}
  Server minor version: {{(index .results .hosts.master).server_minor_ver}}

Slave DB servers:
{{ range $key, $value := .hosts.replicas }}
  DB slave server: {{ $value }}
  {{ if (index $.results $value) }}
    Version: {{(index $.results $value).version}}
    Server version number: {{(index $.results $value).server_version_num}}
    Server major version: {{(index $.results $value).server_major_ver}}
    Server minor version: {{(index $.results $value).server_minor_ver}}
  {{ else }}
    No data
  {{ end}}
{{ end }}

Conclusions
===
{{.Conclusion}}

Recommendations
===
{{.Recommended}}
