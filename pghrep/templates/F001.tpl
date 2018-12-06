Current values
===

Heap bloat

Master DB server is {{.hosts.master}}
{{ range $key, $value := (index (index .results .hosts.master) "data") }}
  Table: {{ $key }}  Size: {{ ( index $value "Size") }}  Extra: {{ ( index $value "Extra") }}  Bloat: {{ ( index $value "Bloat") }}  Live: {{ ( index $value "Live") }} {{ if (index $value "Last Vaccuum") }} Last vaсcuum: {{ ( index $value "Last Vaccuum") }} {{ end }}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
Slave DB servers:
  {{ range $skey, $host := .hosts.replicas }}
  DB slave server: {{ $host }}
    {{ if (index $.results $host) }}
      {{ range $key, $value := (index (index $.results $host) "data") }}
          Table: {{ $key }}  Size: {{ ( index $value "Size") }}  Extra: {{ ( index $value "Extra") }}  Bloat: {{ ( index $value "Bloat") }}  Live: {{ ( index $value "Live") }} {{ if (index $value "Last Vaccuum") }} Last vaсcuum: {{ ( index $value "Last Vaccuum") }} {{ end }}
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
