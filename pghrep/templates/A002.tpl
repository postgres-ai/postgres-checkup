Current values
===

Versions

Master DB server is {{.hosts.master}}
  Version: {{ (index (index .results .hosts.master) "data").version }}
  Server version number: {{ (index (index .results .hosts.master) "data").server_version_num }}
  Server major version: {{ (index (index .results .hosts.master) "data").server_major_ver }}
  Server minor version: {{ (index (index .results .hosts.master) "data").server_minor_ver }}

{{ if gt (len .hosts.replicas) 0 }}
Slave DB servers:
    {{ range $key, $value := .hosts.replicas }}
  DB slave server: {{ $value }}
      {{ if (index $.results $value) }}
    Version: {{ (index (index $.results $value) "data").version }}
    Server version number: {{ (index (index $.results $value) "data").server_version_num }}
    Server major version: {{ (index (index $.results $value) "data").server_major_ver }}
    Server minor version: {{ (index (index $.results $value) "data").server_minor_ver }}
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
